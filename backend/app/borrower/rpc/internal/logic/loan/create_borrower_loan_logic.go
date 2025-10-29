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
	"zeebee-backend/common/types"
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
	loan, propertyBase, loanProperty, applicants, loanBorrower, err := l.buildCreateModels(in.Data)
	if err != nil {
		return nil, err
	}

	// create loan
	err = l.svcCtx.LoanRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
		var innerErr error
		propertyBase, innerErr = l.svcCtx.PropertyBaseRepo.GetOrCreateWithTx(l.ctx, tx, propertyBase)
		if innerErr != nil {
			return innerErr
		}

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
		loanProperty.PropertyBaseId = propertyBase.Id
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
		innerErr = l.svcCtx.PropertyBaseRepo.SyncTotalAmountWithTx(l.ctx, tx, propertyBase.Id)
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
func (l *CreateBorrowerLoanLogic) buildCreateModels(in *borrower.LoanCreateRequestData) (loan *model.Loan, propertyBase *model.PropertyBase, loanProperty *model.LoanProperty, applicants []*model.Applicant, loanBorrower *model.LoanBorrower, err error) {
	// build applicants
	for _, applicant := range in.Applicants {
		if len(applicant.FirstName) == 0 {
			return nil, nil, nil, nil, nil, errors.New("applicant name is required")
		}
		applicantName := strings.TrimSpace(applicant.FirstName)
		if len(applicantName) == 0 {
			return nil, nil, nil, nil, nil, errors.New("applicant name is empty")
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
			return nil, nil, nil, nil, nil, errors.New("invalid applicant type")
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
		return nil, nil, nil, nil, nil, errors.New("no applicant")
	}
	primaryCount := slice.CountBy(applicants, func(index int, item *model.Applicant) bool {
		return item.IsPrimary
	})
	if primaryCount == 0 {
		return nil, nil, nil, nil, nil, errors.New("no primary applicant")
	}
	if primaryCount > 1 {
		return nil, nil, nil, nil, nil, errors.New("multiple primary applicant")
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
		return nil, nil, nil, nil, nil, errorx.GetFailedTransErr(err)
	}
	loan.BorrowerId = in.BorrowerId
	loan.PendingStatus = constants.LoanPendingStatusNone
	loan.Source = constants.LoanSourceBorrower

	// loan = &model.Loan{
	// 	BorrowerId:       in.BorrowerId,
	// 	LenderCompany:    in.LenderCompany,
	// 	ServiceCompany:   in.ServiceCompany,
	// 	Stage:            constants.LoanStage(in.Stage),
	// 	PendingStatus:    constants.LoanPendingStatusNone,
	// 	LoanPurpose:      constants.LoanPurpose(in.LoanPurpose),
	// 	LoanType:         constants.LoanType(in.LoanType),
	// 	TotalAmount:      in.TotalAmount,
	// 	PrincipalAmount:  in.PrincipalAmount,
	// 	InitialAmount:    in.InitialAmount,
	// 	Rate:             in.Rate,
	// 	BusinessStrategy: pq.StringArray(in.BusinessStrategy),
	// 	Source:           constants.LoanSourceBorrower,
	// 	ClosingFund:      in.ClosingFund,
	// 	// float purchase_amount = 24;
	// 	//  string purchase_date = 25;
	// 	//  LoanExtendValueEntry extend = 26;
	// 	//  string existing_lender_company = 27;
	// 	//  string original_contract_date = 28;
	// 	//  string extension_date = 29;
	// 	//  string loan_processor = 30;
	// 	//  string owner = 31;
	// }
	propertyBase = &model.PropertyBase{
		BorrowerId: in.BorrowerId,
		// PropertyType: constants.PropertyType(in.Property.PropertyType),
		Unit: in.LoanProperty.PropertyBase.Unit,
		Address: &types.Address{
			State:    in.LoanProperty.PropertyBase.Address.State,
			City:     in.LoanProperty.PropertyBase.Address.City,
			Address1: in.LoanProperty.PropertyBase.Address.Address1,
			Address2: in.LoanProperty.PropertyBase.Address.Address2,
			ZipCode:  in.LoanProperty.PropertyBase.Address.ZipCode,
		},
	}

	loanProperty = &model.LoanProperty{
		PropertyType: constants.PropertyType(in.LoanProperty.PropertyType),
	}

	if loan.TotalAmount < loan.InitialAmount {
		return nil, nil, nil, nil, nil, errors.New("Total amount cannot be less than initial amount")
	}
	switch constants.LoanStage(in.Stage) {
	case constants.LoanStagePending:
		closingDate, err := timex.FormateTimeToStrWithUsStyle(in.ClosingDate)

		if err != nil {
			return nil, nil, nil, nil, nil, errors.Wrap(err, "invalid closing date")
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
			return nil, nil, nil, nil, nil, errors.Wrap(err, "invalid maturity date")
		}
		loan.MaturityDate = &sql.NullTime{
			Time:  maturityDate,
			Valid: true,
		}
		loan.PrincipalAmount = in.PrincipalAmount
		// loan.MonthlyPayment = in.MonthlyPayment
		// loan.LienPosition = in.LienPosition
		// loan.ServiceCompany = in.ServiceCompany
		// loan.PrepaymentPenalty = in.PrepaymentPenalty
	case constants.LoanStageArchived:
		paidOffDate, err := timex.FormateTimeToStrWithUsStyle(in.PaidOffDate)
		if err != nil {
			return nil, nil, nil, nil, nil, errors.Wrap(err, "invalid paid off date")
		}
		loan.PaidOffDate = &sql.NullTime{
			Time:  paidOffDate,
			Valid: true,
		}
		loan.PrincipalAmount = in.PrincipalAmount
	default:
		return nil, nil, nil, nil, nil, errors.WithStack(errors.New("invalid loan stage"))
	}
	return
}
