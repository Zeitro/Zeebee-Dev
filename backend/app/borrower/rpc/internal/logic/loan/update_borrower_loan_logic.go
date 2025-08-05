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
	loan, property, applicants, err := l.buildCreateModels(in.Data)
	if err != nil {
		return nil, err
	}
	loanModel, err := l.svcCtx.LoanRepo.GetBorrowerLoan(l.ctx, in.Data.Id, in.Data.BorrowerId)
	if err != nil {
		return nil, err
	}
	oldApplicants := loanModel.Applicants
	if len(oldApplicants) == 0 {
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

		// handle applicant
		removedApplicants := slice.DifferenceWith(oldApplicants, applicants, func(item1, item2 *model.Applicant) bool {
			return item1.Id == item2.Id
		})
		removedApplicantIds := slice.Map(removedApplicants, func(index int, item *model.Applicant) int64 {
			return item.Id
		})
		newApplicants := slice.DifferenceWith(applicants, oldApplicants, func(item1, item2 *model.Applicant) bool {
			return item1.Id == item2.Id
		})
		newApplicants = slice.Map(newApplicants, func(index int, newApplicant *model.Applicant) *model.Applicant {
			newApplicant.LoanId = loan.Id
			newApplicant.PropertyId = property.Id
			return newApplicant
		})

		updatedApplicants := slice.Filter(applicants, func(index int, newApplicant *model.Applicant) bool {
			_, ok := slice.FindBy(oldApplicants, func(index int, oldApplicant *model.Applicant) bool {
				return newApplicant.Id == oldApplicant.Id
			})
			return ok
		})

		updatedApplicants = slice.Map(updatedApplicants, func(index int, updateApplicant *model.Applicant) *model.Applicant {
			updateApplicant.LoanId = loan.Id
			updateApplicant.PropertyId = property.Id
			return updateApplicant
		})

		if len(removedApplicantIds) > 0 {
			_, innerErr = l.svcCtx.ApplicantRepo.DeleteByIDsWithTx(l.ctx, tx, removedApplicantIds)
			if innerErr != nil {
				return innerErr
			}
		}

		if len(newApplicants) > 0 {
			_, innerErr = l.svcCtx.ApplicantRepo.BatchCreateWithTx(l.ctx, tx, newApplicants)
			if innerErr != nil {
				return innerErr
			}
		}

		for _, updatedApplicant := range updatedApplicants {
			_, innerErr = l.svcCtx.ApplicantRepo.UpdateOneWithTx(l.ctx, tx, updatedApplicant.Id, updatedApplicant)
			if innerErr != nil {
				return innerErr
			}
		}
		// sync property total amount
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

func (l *UpdateBorrowerLoanLogic) buildCreateModels(in *borrower.LoanUpdateRequestData) (loan *model.Loan, property *model.Property, applicants []*model.Applicant, err error) {
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
			firstName, lastName := stringx.ParseUserName(applicantName)
			applicantItem = &model.Applicant{
				FirstName: firstName,
				LastName:  lastName,
				FullName:  applicantName,
				IsPrimary: applicant.IsPrimary,
				Type:      constants.ApplicantTypeIndividual,
			}
		case constants.ApplicantTypeCompany:
			applicantItem = &model.Applicant{
				FullName:  applicantName,
				IsPrimary: applicant.IsPrimary,
				Type:      constants.ApplicantTypeCompany,
			}
		default:
			return nil, nil, nil, errors.New("invalid applicant type")
		}
		applicantItem.Id = applicant.Id
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
