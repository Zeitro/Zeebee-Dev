package loanlogic

import (
	"context"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/constants"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/timex"
	"zeebee-backend/common/utils"
	"zeebee-backend/common/uuid"
	"zeebee-backend/domains/zeebee/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateBorrowerLoanLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewCreateBorrowerLoanLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateBorrowerLoanLogic {
	return &CreateBorrowerLoanLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.NewBaseLogic(ctx, svcCtx),
	}
}

func (l *CreateBorrowerLoanLogic) CreateBorrowerLoan(in *borrower.LoanCreateRequest) (resp *borrower.LoanCreateResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	loan, loanProperty, applicants, loanBorrower, err := l.buildCreateModels(in.Data, claim)
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
		return nil, errorx.Wrap(err, "Failed To Create Loan")
	}

	return &borrower.LoanCreateResponse{
		Data: &borrower.LoanCreateResponseData{
			Id: loan.Id,
		},
	}, nil
}

// buildCreateModels
func (l *CreateBorrowerLoanLogic) buildCreateModels(in *borrower.LoanCreateRequestData, claim utils.JwtUser) (loan *model.Loan, loanProperty *model.LoanProperty, applicants []*model.LoanApplicant, loanBorrower *model.LoanBorrower, err error) {
	// build applicants
	var totalOwnPercent int32
	for _, applicant := range in.Applicants {
		totalOwnPercent += applicant.OwnPercent
		borrowerApplicant, err := l.svcCtx.BorrowerApplicantRepo.GetBorrowerApplicantById(l.ctx, applicant.BorrowerApplicantId, claim.OrganizationId)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		applicantItem := &model.LoanApplicant{}
		err = copier.CopyWithOption(applicantItem, borrowerApplicant, copier.Option{
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
		applicantItem.Id = 0
		applicantItem.IsPrimary = applicant.IsPrimary
		applicantItem.BorrowerApplicantId = applicant.BorrowerApplicantId
		applicantItem.InvitationStatus = constants.ApplicantInvitationStatusAccepted
		applicants = append(applicants, applicantItem)
	}
	if totalOwnPercent > 100 {
		return nil, nil, nil, nil, errorx.New("Total Own Percent Cannot Be Greater Than 100")
	}

	slice.SortBy(applicants, func(a, b *model.LoanApplicant) bool {
		return b.IsPrimary
	})

	primaryCount := slice.CountBy(applicants, func(index int, item *model.LoanApplicant) bool {
		return item.IsPrimary
	})
	if primaryCount == 0 {
		return nil, nil, nil, nil, errorx.New("No Primary Applicant")
	}
	if primaryCount > 1 {
		return nil, nil, nil, nil, errorx.New("Multiple Primary Applicant")
	}

	// loan borrower
	loanBorrower = &model.LoanBorrower{}

	// LoanProperty
	borrowerProperty, err := l.svcCtx.BorrowerPropertyRepo.GetOneByIdForBorrower(l.ctx, in.LoanProperty.BorrowerPropertyId, claim.OrganizationId)
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
	err = copier.CopyWithOption(loanProperty, in.LoanProperty, copier.Option{
		CaseSensitive: true,
		Converters: []copier.TypeConverter{
			copierx.StrToSqlTimePtr(timex.LayoutUsStyle),
		},
	})
	if err != nil {
		return nil, nil, nil, nil, errorx.GetFailedTransErr(err)
	}
	loanProperty.Id = 0

	// loan
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
	// creator, err := l.svcCtx.UserRepo.GetUser(l.ctx, claim.OrganizationUserId, constants.UserTypeBorrower)
	// if err != nil {
	// 	return nil, nil, nil, nil, err
	// }

	loan.CreatorId = claim.OrganizationUserId
	loan.CreatorName = claim.Username
	loan.BorrowerOrganizationId = claim.OrganizationId
	loan.PendingStatus = constants.LoanPendingStatusNone
	loan.Source = constants.LoanSourceBorrower
	if loan.TotalAmount < loan.InitialAmount {
		return nil, nil, nil, nil, errorx.New("Total Amount Cannot Be Less Than Initial Amount")
	}
	switch constants.LoanStage(in.Stage) {
	case constants.LoanStagePending:
		if loan.ClosingDate == nil {
			return nil, nil, nil, nil, errorx.New("Closing Date Is Required")
		}
		loan.HoldbackAmount = loan.TotalAmount - loan.InitialAmount
		loan.PrincipalAmount = loan.TotalAmount
		loan.PendingStatus = constants.LoanPendingStatus(in.PendingStatus)

	case constants.LoanStageActive:
		// loan.PrincipalAmount = in.PrincipalAmount
		if loan.MaturityDate == nil {
			return nil, nil, nil, nil, errorx.New("Maturity Date Is Required")
		}

	case constants.LoanStageArchived:

		// loan.PrincipalAmount = in.PrincipalAmount
		if loan.PaidOffDate == nil {
			return nil, nil, nil, nil, errorx.New("Paid Off Date Is Required")
		}

	default:
		return nil, nil, nil, nil, errorx.WithStack(errorx.New("Invalid Loan Stage"))
	}
	return
}
