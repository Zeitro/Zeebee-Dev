package loanlogic

import (
	"context"
	"database/sql"
	"strings"

	"github.com/duke-git/lancet/v2/datetime"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"zeebee-backend/app/borrower/rpc/internal/repo/model"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/constants"
	"zeebee-backend/common/stringx"
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

	err = l.svcCtx.LoanRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
		var innerErr error
		property, innerErr = l.svcCtx.PropertyRepo.CreateWithTx(l.ctx, tx, property)
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
		BorrowerID:       in.BorrowerId,
		LenderCompany:    in.LenderCompany,
		ServiceCompany:   in.ServiceCompany,
		Stage:            constants.LoanStage(in.Stage),
		PendingStatus:    constants.LoanPendingStatusClosing,
		Purpose:          constants.LoanPurpose(in.LoanPurpose),
		LoanType:         constants.LoanType(in.LoanType),
		TotalAmount:      in.TotalAmount,
		PrincipalAmount:  in.PrincipalAmount,
		InitialAmount:    in.InitialAmount,
		HoldbackAmount:   in.HoldbackAmount,
		LienPosition:     in.LienPosition,
		Rate:             in.Rate,
		ClosingFund:      in.ClosingFund,
		MonthlyPayment:   in.MonthlyPayment,
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
	switch constants.LoanStage(in.Stage) {
	case constants.LoanStagePending:
		closingDate, err := datetime.FormatStrToTime(in.ClosingDate, "mm/dd/yyyy")
		if err != nil {
			return nil, nil, nil, errors.Wrap(err, "invalid closing date")
		}
		loan.ClosingDate = &sql.NullTime{
			Time:  closingDate,
			Valid: true,
		}
	case constants.LoanStageActive:
		maturityDate, err := datetime.FormatStrToTime(in.MaturityDate, "mm/dd/yyyy")
		if err != nil {
			return nil, nil, nil, errors.Wrap(err, "invalid maturity date")
		}
		loan.MaturityDate = &sql.NullTime{
			Time:  maturityDate,
			Valid: true,
		}
	case constants.LoanStageArchived:
		paidOffDate, err := datetime.FormatStrToTime(in.PaidOffDate, "mm/dd/yyyy")
		if err != nil {
			return nil, nil, nil, errors.Wrap(err, "invalid paid off date")
		}
		loan.PaidOffDate = &sql.NullTime{
			Time:  paidOffDate,
			Valid: true,
		}
	default:
		return nil, nil, nil, errors.WithStack(errors.New("invalid loan stage"))
	}
	return
}
