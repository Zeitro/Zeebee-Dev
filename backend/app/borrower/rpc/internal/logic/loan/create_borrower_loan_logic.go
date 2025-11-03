package loanlogic

import (
	"context"
	"strings"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"zeebee-backend/app/borrower/rpc/internal/repo/model"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/constants"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/stringx"
	"zeebee-backend/common/timex"
	"zeebee-backend/common/uuid"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type CreateBorrowerLoanLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateBorrowerLoanLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateBorrowerLoanLogic {
	return &CreateBorrowerLoanLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateBorrowerLoanLogic) CreateBorrowerLoan(in *borrower.LoanCreateRequest) (resp *borrower.LoanCreateResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	loan, loanProperty, applicants, loanBorrower, err := l.buildCreateModels(in.Data)
	if err != nil {
		return nil, err
	}

	// create loan
	err = l.svcCtx.LoanRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
		var innerErr error

		// loan.PropertyId = property.Id
		uid, innerErr := uuid.GetSonyflakeIDString()
		if innerErr != nil {
			return innerErr
		}
		loan.Code = uid
		loan, innerErr = l.svcCtx.LoanRepo.CreateWithTx(l.ctx, tx, loan, false)
		if innerErr != nil {
			return innerErr
		}

		// create loan property
		loanProperty.LoanId = loan.Id
		// loanProperty.PropertyBaseId = propertyBase.Id
		loanProperty, innerErr = l.svcCtx.LoanPropertyRepo.CreateWithTx(l.ctx, tx, loanProperty, false)
		if innerErr != nil {
			return innerErr
		}

		// loan borrower
		loanBorrower.LoanId = loan.Id
		loanBorrower, innerErr = l.svcCtx.LoanBorrowerRepo.CreateWithTx(l.ctx, tx, loanBorrower, false)
		if innerErr != nil {
			return innerErr
		}

		// create applicants
		applicants = slice.Map(applicants, func(index int, item *model.LoanApplicant) *model.LoanApplicant {
			item.LoanId = loan.Id
			item.LoanPropertyId = loanProperty.Id
			item.LoanBorrowerId = loanBorrower.Id
			return item
		})
		_, innerErr = l.svcCtx.LoanApplicantRepo.BatchCreateWithTx(l.ctx, tx, applicants, false)
		if innerErr != nil {
			return innerErr
		}

		// sync total amount
		innerErr = l.svcCtx.BorrowerPropertyRepo.SyncTotalAmountAndLoanCountWithTx(l.ctx, tx, loanProperty.BorrowerPropertyId)
		return innerErr
	})
	if err != nil {
		return nil, errors.Wrap(err, "Failed To Create Loan")
	}

	return &borrower.LoanCreateResponse{
		Data: &borrower.LoanCreateResponseData{
			Id: loan.Id,
		},
	}, nil
}

// buildCreateModels
func (l *CreateBorrowerLoanLogic) buildCreateModels(in *borrower.LoanCreateRequestData) (loan *model.Loan, loanProperty *model.LoanProperty, applicants []*model.LoanApplicant, loanBorrower *model.LoanBorrower, err error) {
	// build applicants
	for _, applicant := range in.Applicants {
		if len(applicant.FirstName) == 0 {
			return nil, nil, nil, nil, errors.New("Applicant Name Is Required")
		}
		applicantName := strings.TrimSpace(applicant.FirstName)
		if len(applicantName) == 0 {
			return nil, nil, nil, nil, errors.New("Applicant Name Is Empty")
		}
		applicantItem := &model.LoanApplicant{}
		err = copier.CopyWithOption(applicantItem, applicant, copier.Option{
			IgnoreEmpty:   false,
			DeepCopy:      true,
			CaseSensitive: true,
			Converters: []copier.TypeConverter{
				copierx.StrToSqlTimePtr(timex.LayoutUsStyle),
			},
		})
		if err != nil {
			return nil, nil, nil, nil, errorx.GetFailedTransErr(err)
		}
		// if !applicantItem.Type.IsValid() {
		// 	return nil, nil, nil, nil, errors.New("Invalid Applicant Type")
		// }
		applicantItem.FullName = stringx.FormatFullName(applicant.FirstName, applicant.MidName, applicant.LastName)
		applicantItem.BorrowerId = in.BorrowerId
		applicantItem.InvitationStatus = constants.ApplicantInvitationStatusAccepted
		applicants = append(applicants, applicantItem)
	}
	slice.SortBy(applicants, func(a, b *model.LoanApplicant) bool {
		return b.IsPrimary
	})
	applicants = slice.UniqueByComparator(applicants, func(item *model.LoanApplicant, other *model.LoanApplicant) bool {
		return item.FullName == other.FullName
	})
	if len(applicants) == 0 {
		return nil, nil, nil, nil, errors.New("No Applicant")
	}
	primaryCount := slice.CountBy(applicants, func(index int, item *model.LoanApplicant) bool {
		return item.IsPrimary
	})
	if primaryCount == 0 {
		return nil, nil, nil, nil, errors.New("No Primary Applicant")
	}
	if primaryCount > 1 {
		return nil, nil, nil, nil, errors.New("Multiple Primary Applicant")
	}
	// loan borrower
	loanBorrower = &model.LoanBorrower{
		Name: applicants[0].FullName,
		// Type:       applicants[0].Type,
		OwnPercent: 100,
	}

	// update loan
	loan = &model.Loan{}
	err = copier.CopyWithOption(loan, in, copier.Option{
		IgnoreEmpty:   false,
		DeepCopy:      true,
		CaseSensitive: true,
		Converters: []copier.TypeConverter{
			copierx.StrToSqlTimePtr(timex.LayoutUsStyle),
		},
	})
	if err != nil {
		return nil, nil, nil, nil, errorx.GetFailedTransErr(err)
	}
	loan.BorrowerId = in.BorrowerId
	loan.PendingStatus = constants.LoanPendingStatusNone
	loan.Source = constants.LoanSourceBorrower

	borrowerProperty, err := l.svcCtx.BorrowerPropertyRepo.FindOneByColumns(l.ctx, map[string]interface{}{
		"id":          in.LoanProperty.BorrowerPropertyId,
		"borrower_id": in.BorrowerId,
	})
	if err != nil {
		return nil, nil, nil, nil, err
	}

	loanProperty = &model.LoanProperty{}

	err = copier.CopyWithOption(loanProperty, borrowerProperty, copier.Option{
		IgnoreEmpty:   false,
		DeepCopy:      true,
		CaseSensitive: true,
		Converters: []copier.TypeConverter{
			copierx.StrToSqlTimePtr(timex.LayoutUsStyle),
		},
	})
	if err != nil {
		return nil, nil, nil, nil, errorx.GetFailedTransErr(err)
	}
	// Clear ID to allow auto-generation by database
	loanProperty.Id = 0
	loanProperty.BorrowerPropertyId = in.LoanProperty.BorrowerPropertyId
	loanProperty.Owner = in.LoanProperty.Owner
	loanProperty.PurchaseAmount = in.LoanProperty.PurchaseAmount

	if loan.TotalAmount < loan.InitialAmount {
		return nil, nil, nil, nil, errors.New("Total Amount Cannot Be Less Than Initial Amount")
	}
	switch constants.LoanStage(in.Stage) {
	case constants.LoanStagePending:

		if loan.ClosingDate == nil {
			return nil, nil, nil, nil, errors.New("Closing Date Is Required")
		}
		loan.HoldbackAmount = loan.TotalAmount - loan.InitialAmount
		loan.PrincipalAmount = loan.TotalAmount
		// loan.PendingStatus = constants.LoanPendingStatus(in.PendingStatus)

	case constants.LoanStageActive:
		// loan.PrincipalAmount = in.PrincipalAmount
		if loan.MaturityDate == nil {
			return nil, nil, nil, nil, errors.New("Maturity Date Is Required")
		}

	case constants.LoanStageArchived:

		// loan.PrincipalAmount = in.PrincipalAmount
		if loan.PaidOffDate == nil {
			return nil, nil, nil, nil, errors.New("Paid Off Date Is Required")
		}

	default:
		return nil, nil, nil, nil, errors.WithStack(errors.New("Invalid Loan Stage"))
	}
	return
}
