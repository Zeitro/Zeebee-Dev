package loanlogic

import (
	"context"
	"database/sql"
	"strings"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"zeebee-backend/app/borrower/rpc/internal/repo/model"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/constants"
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

func (l *CreateBorrowerLoanLogic) CreateBorrowerLoan(in *borrower.LoanCreateRequest) (resp *borrower.LoanCreateResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	loan, property, applicants, err := l.buildCreateModels(in.Data)
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

		loan.PropertyId = property.Id
		loan, innerErr = l.svcCtx.LoanRepo.CreateWithTx(l.ctx, tx, loan)
		if innerErr != nil {
			return innerErr
		}

		applicants = slice.Map(applicants, func(index int, item *model.Applicant) *model.Applicant {
			item.LoanId = loan.Id
			item.PropertyId = property.Id
			return item
		})
		_, innerErr = l.svcCtx.ApplicantRepo.BatchCreateWithTx(l.ctx, tx, applicants)
		if innerErr != nil {
			return innerErr
		}
		innerErr = l.svcCtx.PropertyRepo.SyncTotalAmountWithTx(l.ctx, tx, property.Id)
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
func (l *CreateBorrowerLoanLogic) buildCreateModels(in *borrower.LoanCreateRequestData) (loan *model.Loan, property *model.Property, applicants []*model.Applicant, err error) {
	// build applicants
	for _, applicant := range in.Applicants {
		if len(applicant.Name) == 0 {
			return nil, nil, nil, errors.New("applicant name is required")
		}
		applicantName := strings.TrimSpace(applicant.Name)
		if len(applicantName) == 0 {
			return nil, nil, nil, errors.New("applicant name is empty")
		}
		var applicantItem *model.Applicant
		switch constants.ApplicantType(applicant.Type) {
		case constants.ApplicantTypeIndividual:
			applicantItem = &model.Applicant{
				Name:      applicantName,
				IsPrimary: applicant.IsPrimary,
				Type:      constants.ApplicantTypeIndividual,
			}
		case constants.ApplicantTypeCompany:
			applicantItem = &model.Applicant{
				Name:      applicantName,
				IsPrimary: applicant.IsPrimary,
				Type:      constants.ApplicantTypeCompany,
			}
		default:
			return nil, nil, nil, errors.New("invalid applicant type")
		}
		applicantItem.BorrowerId = in.BorrowerId
		applicantItem.InvitationStatus = constants.ApplicantInvitationStatusAccepted
		applicants = append(applicants, applicantItem)
	}
	slice.SortBy(applicants, func(a, b *model.Applicant) bool {
		return b.IsPrimary
	})
	applicants = slice.UniqueByComparator(applicants, func(item *model.Applicant, other *model.Applicant) bool {
		return item.Name == other.Name
	})
	if len(applicants) == 0 {
		return nil, nil, nil, errors.New("no applicant")
	}
	primaryCount := slice.CountBy(applicants, func(index int, item *model.Applicant) bool {
		return item.IsPrimary
	})
	if primaryCount == 0 {
		return nil, nil, nil, errors.New("no primary applicant")
	}
	if primaryCount > 1 {
		return nil, nil, nil, errors.New("multiple primary applicant")
	}
	// update loan
	loan = &model.Loan{
		BorrowerId:       in.BorrowerId,
		LenderCompany:    in.LenderCompany,
		ServiceCompany:   in.ServiceCompany,
		Stage:            constants.LoanStage(in.Stage),
		PendingStatus:    constants.LoanPendingStatusClosing,
		LoanPurpose:      constants.LoanPurpose(in.LoanPurpose),
		LoanType:         constants.LoanType(in.LoanType),
		TotalAmount:      in.TotalAmount,
		PrincipalAmount:  in.PrincipalAmount,
		InitialAmount:    in.InitialAmount,
		Rate:             in.Rate,
		BusinessStrategy: pq.StringArray(in.BusinessStrategy),
		Source:           constants.LoanSourceBorrower,
	}
	property = &model.Property{
		BorrowerId:   in.BorrowerId,
		PropertyType: constants.PropertyType(in.Property.PropertyType),
		Unit:         in.Property.Unit,
		Address: types.Address{
			State:    in.Property.Address.State,
			City:     in.Property.Address.City,
			Address1: in.Property.Address.Address1,
			Address2: in.Property.Address.Address2,
			ZipCode:  in.Property.Address.ZipCode,
		},
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
