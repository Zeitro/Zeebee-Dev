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

func (l *CreateBorrowerLoanLogic) CreateBorrowerLoan(in *borrower.CreateLoanRequest) (resp *borrower.CreateLoanResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	loan, property, applicant, err := l.buildCreateModels(in)
	if err != nil {
		return nil, err
	}

	// create loan
	err = l.svcCtx.LoanRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
		var innerErr error
		property, innerErr = l.svcCtx.PropertyRepo.GetOrCreateWithTx(l.ctx, tx, property)
		if innerErr != nil {
			return innerErr
		}

		loan.PropertyID = property.ID
		loan, innerErr = l.svcCtx.LoanRepo.CreateWithTx(l.ctx, tx, loan)
		if innerErr != nil {
			return innerErr
		}
		applicant.LoanID = loan.ID
		applicant, innerErr = l.svcCtx.ApplicantRepo.CreateWithTx(l.ctx, tx, applicant)

		return innerErr
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create loan")
	}

	return &borrower.CreateLoanResponse{
		Id: loan.ID,
	}, nil
}

// buildCreateModels
func (l *CreateBorrowerLoanLogic) buildCreateModels(in *borrower.CreateLoanRequest) (loan *model.Loan, property *model.Property, applicant *model.Applicant, err error) {
	in.ApplicantName = strings.TrimSpace(in.ApplicantName)
	if len(in.ApplicantName) == 0 {
		return nil, nil, nil, errors.New("invalid applicant name")
	}
	firstName, lastName := stringx.ParseUserName(in.ApplicantName)

	loan = &model.Loan{
		BorrowerID:      in.BorrowerId,
		LenderCompany:   in.LenderCompany,
		ServiceCompany:  in.ServiceCompany,
		Stage:           constants.LoanStage(in.Stage),
		PendingStatus:   constants.LoanPendingStatusClosing,
		Purpose:         constants.LoanPurpose(in.LoanPurpose),
		LoanType:        constants.LoanType(in.LoanType),
		TotalAmount:     in.TotalAmount,
		PrincipalAmount: in.PrincipalAmount,
		InitialAmount:   in.InitialAmount,
		// HoldbackAmount:   in.HoldbackAmount,
		// LienPosition:     "1st",
		Rate: in.Rate,
		// ClosingFund:      in.ClosingFund,
		// MonthlyPayment:   in.MonthlyPayment,
		BusinessStrategy: pq.StringArray(in.BusinessStrategy),
		Source:           constants.LoanSourceBorrower,
	}
	property = &model.Property{
		BorrowerID: in.BorrowerId,
		Type:       constants.PropertyType(in.Property.PropertyType),
		Unit:       in.Property.Unit,
		Address: types.Address{
			State:    in.Property.Address.State,
			City:     in.Property.Address.City,
			Address1: in.Property.Address.Address1,
			Address2: in.Property.Address.Address2,
			ZipCode:  in.Property.Address.ZipCode,
		},
	}
	applicant = &model.Applicant{
		BorrowerID: in.BorrowerId,
		LoanID:     0,
		FirstName:  firstName,
		LastName:   lastName,
		FullName:   in.ApplicantName,
		IsPrimary:  true,
	}
	if loan.TotalAmount < loan.InitialAmount {
		return nil, nil, nil, errors.New("Total amount cannot be less than initial amount")
	}
	switch constants.LoanStage(in.Stage) {
	case constants.LoanStagePending:
		closingDate, err := timex.FormateTimeToStrWithUsStyle(in.ClosingDate)

		if err != nil {
			return nil, nil, nil, errors.Wrap(err, "invalid closing date")
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
			return nil, nil, nil, errors.Wrap(err, "invalid maturity date")
		}
		loan.MaturityDate = &sql.NullTime{
			Time:  maturityDate,
			Valid: true,
		}
		loan.PrincipalAmount = in.PrincipalAmount
		loan.MonthlyPayment = in.MonthlyPayment
		loan.LienPosition = in.LienPosition
		loan.ServiceCompany = in.ServiceCompany
		loan.PrepaymentPenalty = in.PrepaymentPenalty
	case constants.LoanStageArchived:
		paidOffDate, err := timex.FormateTimeToStrWithUsStyle(in.PaidOffDate)
		if err != nil {
			return nil, nil, nil, errors.Wrap(err, "invalid paid off date")
		}
		loan.PaidOffDate = &sql.NullTime{
			Time:  paidOffDate,
			Valid: true,
		}
		loan.PrincipalAmount = in.PrincipalAmount
	default:
		return nil, nil, nil, errors.WithStack(errors.New("invalid loan stage"))
	}
	return
}
