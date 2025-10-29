package loanlogic

import (
	"context"
	"database/sql"
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
		// propertyBase, innerErr = l.svcCtx.PropertyBaseRepo.GetOrCreateWithTx(l.ctx, tx, propertyBase)
		// if innerErr != nil {
		// 	return innerErr
		// }

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
		applicants = slice.Map(applicants, func(index int, item *model.Applicant) *model.Applicant {
			item.LoanId = loan.Id
			item.LoanPropertyId = loanProperty.Id
			item.LoanBorrowerId = loanBorrower.Id
			return item
		})
		_, innerErr = l.svcCtx.ApplicantRepo.BatchCreateWithTx(l.ctx, tx, applicants, false)
		if innerErr != nil {
			return innerErr
		}

		// sync total amount
		innerErr = l.svcCtx.PropertyBaseRepo.SyncTotalAmountAndLoanCountWithTx(l.ctx, tx, loanProperty.PropertyBaseId)
		return innerErr
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create loan")
	}

	return &borrower.LoanCreateResponse{
		Data: &borrower.LoanCreateResponseData{
			Id: loan.Id,
		},
	}, nil
}

// buildCreateModels
func (l *CreateBorrowerLoanLogic) buildCreateModels(in *borrower.LoanCreateRequestData) (loan *model.Loan, loanProperty *model.LoanProperty, applicants []*model.Applicant, loanBorrower *model.LoanBorrower, err error) {
	// build applicants
	for _, applicant := range in.Applicants {
		if len(applicant.FirstName) == 0 {
			return nil, nil, nil, nil, errors.New("applicant name is required")
		}
		applicantName := strings.TrimSpace(applicant.FirstName)
		if len(applicantName) == 0 {
			return nil, nil, nil, nil, errors.New("applicant name is empty")
		}
		var applicantItem *model.Applicant
		applicantItem = &model.Applicant{
			// Name:      applicantName,
			FirstName: applicant.FirstName,
			LastName:  applicant.LastName,
			MidName:   applicant.MidName,
			FullName:  stringx.FormatFullName(applicant.FirstName, applicant.MidName, applicant.LastName),
			IsPrimary: applicant.IsPrimary,
			Type:      constants.ApplicantType(applicant.Type),
		}

		if !applicantItem.Type.IsValid() {
			return nil, nil, nil, nil, errors.New("invalid applicant type")
		}
		applicantItem.BorrowerId = in.BorrowerId
		applicantItem.InvitationStatus = constants.ApplicantInvitationStatusAccepted
		applicants = append(applicants, applicantItem)
	}
	slice.SortBy(applicants, func(a, b *model.Applicant) bool {
		return b.IsPrimary
	})
	applicants = slice.UniqueByComparator(applicants, func(item *model.Applicant, other *model.Applicant) bool {
		return item.FullName == other.FullName
	})
	if len(applicants) == 0 {
		return nil, nil, nil, nil, errors.New("no applicant")
	}
	primaryCount := slice.CountBy(applicants, func(index int, item *model.Applicant) bool {
		return item.IsPrimary
	})
	if primaryCount == 0 {
		return nil, nil, nil, nil, errors.New("no primary applicant")
	}
	if primaryCount > 1 {
		return nil, nil, nil, nil, errors.New("multiple primary applicant")
	}
	// loan borrower
	loanBorrower = &model.LoanBorrower{
		Name:       applicants[0].FullName,
		Type:       applicants[0].Type,
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

	propertyBase, err := l.svcCtx.PropertyBaseRepo.FindOneByColumns(l.ctx, map[string]interface{}{
		"id":          in.LoanProperty.PropertyBaseId,
		"borrower_id": in.BorrowerId,
	})
	if err != nil {
		return nil, nil, nil, nil, err
	}

	loanProperty = &model.LoanProperty{
		PropertyBaseId:       in.LoanProperty.PropertyBaseId,
		Unit:                 propertyBase.Unit,
		PropertyType:         propertyBase.PropertyType,
		LivingAreaSqft:       propertyBase.LivingAreaSqft,
		Bedrooms:             propertyBase.Bedrooms,
		Bathrooms:            propertyBase.Bathrooms,
		FuturePropertyType:   propertyBase.PropertyType,
		FutureUnit:           propertyBase.Unit,
		FutureLivingAreaSqft: propertyBase.LivingAreaSqft,
		FutureBedrooms:       propertyBase.Bedrooms,
		FutureBathrooms:      propertyBase.Bathrooms,
		Address:              propertyBase.Address,
		Owner:                in.LoanProperty.Owner,
		PurchaseAmount:       in.LoanProperty.PurchaseAmount,
	}
	if len(in.LoanProperty.PurchaseDate) > 0 {
		purchaseDate, err := timex.FormateTimeToStrWithUsStyle(in.LoanProperty.PurchaseDate)
		if err != nil {
			return nil, nil, nil, nil, errors.Wrap(err, "Invalid Purchase Date")
		}
		loanProperty.PurchaseDate = &sql.NullTime{
			Time:  purchaseDate,
			Valid: true,
		}

	}

	if loan.TotalAmount < loan.InitialAmount {
		return nil, nil, nil, nil, errors.New("Total Amount Cannot Be Less Than Initial Amount")
	}
	switch constants.LoanStage(in.Stage) {
	case constants.LoanStagePending:
		closingDate, err := timex.FormateTimeToStrWithUsStyle(in.ClosingDate)

		if err != nil {
			return nil, nil, nil, nil, errors.Wrap(err, "Invalid Closing Date")
		}
		loan.ClosingDate = &sql.NullTime{
			Time:  closingDate,
			Valid: true,
		}
		loan.HoldbackAmount = loan.TotalAmount - loan.InitialAmount
		loan.PrincipalAmount = loan.TotalAmount
		loan.PendingStatus = constants.LoanPendingStatus(in.PendingStatus)

	case constants.LoanStageActive:

		maturityDate, err := timex.FormateTimeToStrWithUsStyle(in.MaturityDate)
		if err != nil {
			return nil, nil, nil, nil, errors.Wrap(err, "Invalid Maturity Date")
		}
		loan.MaturityDate = &sql.NullTime{
			Time:  maturityDate,
			Valid: true,
		}
		loan.PrincipalAmount = in.PrincipalAmount
	case constants.LoanStageArchived:
		paidOffDate, err := timex.FormateTimeToStrWithUsStyle(in.PaidOffDate)
		if err != nil {
			return nil, nil, nil, nil, errors.Wrap(err, "Invalid Paid Off Date")
		}
		loan.PaidOffDate = &sql.NullTime{
			Time:  paidOffDate,
			Valid: true,
		}
		loan.PrincipalAmount = in.PrincipalAmount
	default:
		return nil, nil, nil, nil, errors.WithStack(errors.New("Invalid Loan Stage"))
	}
	return
}
