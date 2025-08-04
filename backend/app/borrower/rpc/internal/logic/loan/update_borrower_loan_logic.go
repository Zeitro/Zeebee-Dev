package loanlogic

import (
	"context"
	"database/sql"
	"strings"

	"github.com/lib/pq"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"zeebee-backend/app/borrower/rpc/internal/repo/model"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/constants"
	"zeebee-backend/common/stringx"
	"zeebee-backend/common/timex"
	"zeebee-backend/common/types"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type UpdateBorrowerLoanLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateBorrowerLoanLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateBorrowerLoanLogic {
	return &UpdateBorrowerLoanLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateBorrowerLoanLogic) UpdateBorrowerLoan(in *borrower.LoanUpdateRequest) (resp *borrower.EmptyResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	loan, property, applicant, err := l.buildCreateModels(in)
	if err != nil {
		return nil, err
	}
	loanModel, err := l.svcCtx.LoanRepo.GetBorrowerLoan(l.ctx, in.Data.Id, in.Data.BorrowerId)
	if err != nil {
		return nil, err
	}
	if len(loanModel.Applicants) == 0 {
		return nil, errors.WithStack(errors.New("applicant is not found"))
	}
	if loanModel.Source != constants.LoanSourceBorrower {
		return nil, errors.WithStack(errors.New("loan source is not 3rd party"))
	}

	// create loan
	err = l.svcCtx.LoanRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
		var innerErr error
		// handle property
		property, innerErr = l.handleProperty(tx, loanModel.PropertyId, property)
		if innerErr != nil {
			return innerErr
		}
		previousPropertyId := loan.PropertyId

		// update loan
		loan.Id = in.Data.Id
		loan.PropertyId = property.Id
		_, innerErr = l.svcCtx.LoanRepo.UpdateOneWithTx(l.ctx, tx, in.Data.Id, loan)
		if innerErr != nil {
			return innerErr
		}

		// update applicant
		applicant.Id = loanModel.Applicants[0].Id
		applicant.LoanId = loan.Id
		applicant.PropertyId = property.Id
		_, innerErr = l.svcCtx.ApplicantRepo.UpdateOneWithTx(l.ctx, tx, in.Data.Id, applicant)
		if innerErr != nil {
			return innerErr
		}
		_ = l.svcCtx.PropertyRepo.SyncTotalAmountWithTx(l.ctx, tx, previousPropertyId)
		if property.Id != previousPropertyId {
			_ = l.svcCtx.PropertyRepo.SyncTotalAmountWithTx(l.ctx, tx, property.Id)
		}
		return innerErr
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to update loan")
	}

	return &borrower.EmptyResponse{}, nil
}

// handle property
func (l *UpdateBorrowerLoanLogic) handleProperty(tx *gorm.DB, propertyId int64, property *model.Property) (newProperty *model.Property, err error) {
	// get existed property
	newProperty, err = l.svcCtx.PropertyRepo.GetOneWithTxByProperty(l.ctx, tx, property)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get property")
	}
	if newProperty != nil {
		return newProperty, nil
	}
	// get previous property
	prevProperty, err := l.svcCtx.PropertyRepo.GetOneWithTx(l.ctx, tx, propertyId, []string{"Loans"})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get property")
	}

	// create new property if relate zero or more than one loan
	if len(prevProperty.Loans) == 0 || len(prevProperty.Loans) > 1 {
		return l.svcCtx.PropertyRepo.GetOrCreateWithTx(l.ctx, tx, property)
	}

	// update property
	property.Id = prevProperty.Id
	newProperty, err = l.svcCtx.PropertyRepo.UpdateOneWithTx(l.ctx, tx, prevProperty.Id, property)
	return
}

func (l *UpdateBorrowerLoanLogic) buildCreateModels(in *borrower.LoanUpdateRequest) (loan *model.Loan, property *model.Property, applicant *model.Applicant, err error) {
	in.Data.ApplicantName = strings.TrimSpace(in.Data.ApplicantName)
	if len(in.Data.ApplicantName) == 0 {
		return nil, nil, nil, errors.New("invalid applicant name")
	}
	firstName, lastName := stringx.ParseUserName(in.Data.ApplicantName)

	loan = &model.Loan{
		BorrowerId:       in.Data.BorrowerId,
		LenderCompany:    in.Data.LenderCompany,
		ServiceCompany:   in.Data.ServiceCompany,
		Stage:            constants.LoanStage(in.Data.Stage),
		PendingStatus:    constants.LoanPendingStatusClosing,
		LoanPurpose:      constants.LoanPurpose(in.Data.LoanPurpose),
		LoanType:         constants.LoanType(in.Data.LoanType),
		TotalAmount:      in.Data.TotalAmount,
		PrincipalAmount:  in.Data.PrincipalAmount,
		InitialAmount:    in.Data.InitialAmount,
		Rate:             in.Data.Rate,
		BusinessStrategy: pq.StringArray(in.Data.BusinessStrategy),
		Source:           constants.LoanSourceBorrower,
	}
	property = &model.Property{
		BorrowerId:   in.Data.BorrowerId,
		PropertyType: constants.PropertyType(in.Data.Property.PropertyType),
		Unit:         in.Data.Property.Unit,
		Address: types.Address{
			State:    in.Data.Property.Address.State,
			City:     in.Data.Property.Address.City,
			Address1: in.Data.Property.Address.Address1,
			Address2: in.Data.Property.Address.Address2,
			ZipCode:  in.Data.Property.Address.ZipCode,
		},
	}
	applicant = &model.Applicant{
		BorrowerId: in.Data.BorrowerId,
		LoanId:     0,
		FirstName:  firstName,
		LastName:   lastName,
		FullName:   in.Data.ApplicantName,
		IsPrimary:  true,
	}
	if loan.TotalAmount < loan.InitialAmount {
		return nil, nil, nil, errors.New("Total amount cannot be less than initial amount")
	}
	switch constants.LoanStage(in.Data.Stage) {
	case constants.LoanStagePending:
		closingDate, err := timex.FormateTimeToStrWithUsStyle(in.Data.ClosingDate)

		if err != nil {
			return nil, nil, nil, errors.Wrap(err, "invalid closing date")
		}
		loan.ClosingDate = &sql.NullTime{
			Time:  closingDate,
			Valid: true,
		}
		loan.HoldbackAmount = loan.TotalAmount - loan.InitialAmount
		loan.PrincipalAmount = loan.TotalAmount
		loan.PendingStatus = constants.LoanPendingStatus(in.Data.PendingStatus)

	case constants.LoanStageActive:

		maturityDate, err := timex.FormateTimeToStrWithUsStyle(in.Data.MaturityDate)
		if err != nil {
			return nil, nil, nil, errors.Wrap(err, "invalid maturity date")
		}
		loan.MaturityDate = &sql.NullTime{
			Time:  maturityDate,
			Valid: true,
		}
		loan.PrincipalAmount = in.Data.PrincipalAmount
		loan.MonthlyPayment = in.Data.MonthlyPayment
		loan.LienPosition = in.Data.LienPosition
		loan.ServiceCompany = in.Data.ServiceCompany
		loan.PrepaymentPenalty = in.Data.PrepaymentPenalty
	case constants.LoanStageArchived:
		paidOffDate, err := timex.FormateTimeToStrWithUsStyle(in.Data.PaidOffDate)
		if err != nil {
			return nil, nil, nil, errors.Wrap(err, "invalid paid off date")
		}
		loan.PaidOffDate = &sql.NullTime{
			Time:  paidOffDate,
			Valid: true,
		}
		loan.PrincipalAmount = in.Data.PrincipalAmount
	default:
		return nil, nil, nil, errors.WithStack(errors.New("invalid loan stage"))
	}
	return
}
