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
	loan, propertyBase, loanProperty, applicants, err := l.buildCreateModels(in.Data)
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
	loanPropertyModel := loanModel.LoanProperty
	if loanPropertyModel == nil {
		return nil, errors.WithStack(errors.New("loan property is not found"))
	}

	// create loan
	err = l.svcCtx.LoanRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
		var innerErr error
		// handle property
		propertyBase, innerErr = l.handlePropertyBase(tx, loanPropertyModel.PropertyBaseId, propertyBase)
		if innerErr != nil {
			return innerErr
		}
		previousPropertyBaseId := loanPropertyModel.PropertyBaseId

		// update loan
		loan.Id = in.Data.Id
		// loan.PropertyId = property.Id
		_, innerErr = l.svcCtx.LoanRepo.UpdateOneWithTx(l.ctx, tx, in.Data.Id, loan)
		if innerErr != nil {
			return innerErr
		}

		// update loan property
		innerErr = l.svcCtx.LoanPropertyRepo.UpdateColumnsByIdWithTx(l.ctx, tx, loanPropertyModel.Id, map[string]interface{}{
			"property_base_id": propertyBase.Id,
			"property_type":    loanProperty.PropertyType,
		})
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
			newApplicant.LoanPropertyId = loanPropertyModel.Id
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
			updateApplicant.LoanPropertyId = loanPropertyModel.Id
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
		_ = l.svcCtx.PropertyBaseRepo.SyncTotalAmountWithTx(l.ctx, tx, previousPropertyBaseId)
		if propertyBase.Id != previousPropertyBaseId {
			_ = l.svcCtx.PropertyBaseRepo.SyncTotalAmountWithTx(l.ctx, tx, propertyBase.Id)
		}
		return innerErr
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to update loan")
	}

	return &borrower.EmptyResponse{}, nil
}

// handle property
func (l *UpdateBorrowerLoanLogic) handlePropertyBase(tx *gorm.DB, propertyId int64, propertyBase *model.PropertyBase) (newPropertyBase *model.PropertyBase, err error) {
	// get existed property
	newPropertyBase, err = l.svcCtx.PropertyBaseRepo.GetOneWithTxByProperty(l.ctx, tx, propertyBase)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get property base")
	}
	if newPropertyBase != nil {
		return newPropertyBase, nil
	}

	// get previous property
	prevPropertyBase, err := l.svcCtx.PropertyBaseRepo.GetOneWithTx(l.ctx, tx, propertyId, []string{"LoanProperties"})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get property")
	}

	// create new property if relate zero or more than one loan property
	if len(prevPropertyBase.LoanProperties) == 0 || len(prevPropertyBase.LoanProperties) > 1 {
		return l.svcCtx.PropertyBaseRepo.GetOrCreateWithTx(l.ctx, tx, propertyBase)
	}

	// update property when relate one loan property
	propertyBase.Id = prevPropertyBase.Id
	newPropertyBase, err = l.svcCtx.PropertyBaseRepo.UpdateOneWithTx(l.ctx, tx, propertyBase.Id, propertyBase)
	return
}

func (l *UpdateBorrowerLoanLogic) buildCreateModels(in *borrower.LoanUpdateRequestData) (loan *model.Loan, propertyBase *model.PropertyBase, loanProperty *model.LoanProperty, applicants []*model.Applicant, err error) {
	// build applicants
	for _, applicant := range in.Applicants {
		if len(applicant.Name) == 0 {
			return nil, nil, nil, nil, errors.New("applicant name is required")
		}
		applicantName := strings.TrimSpace(applicant.Name)
		if len(applicantName) == 0 {
			return nil, nil, nil, nil, errors.New("applicant name is empty")
		}
		var applicantItem *model.Applicant
		switch constants.ApplicantType(applicant.Type) {
		case constants.ApplicantTypeIndividual:
			applicantItem = &model.Applicant{
				Name:      applicantName,
				IsPrimary: applicant.IsPrimary,
				Type:      constants.ApplicantTypeIndividual,
			}
		case constants.ApplicantTypeCorporation:
			applicantItem = &model.Applicant{
				Name:      applicantName,
				IsPrimary: applicant.IsPrimary,
				Type:      constants.ApplicantTypeCorporation,
			}
		default:
			return nil, nil, nil, nil, errors.New("invalid applicant type")
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
		return item.Name == other.Name
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
		ClosingFund:      in.ClosingFund,
	}
	propertyBase = &model.PropertyBase{
		BorrowerId: in.BorrowerId,
		// PropertyType: constants.PropertyType(in.Property.PropertyType),
		Unit: in.LoanProperty.PropertyBase.Unit,
		Address: types.Address{
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
		return nil, nil, nil, nil, errors.New("Total amount cannot be less than initial amount")
	}
	switch constants.LoanStage(in.Stage) {
	case constants.LoanStagePending:
		closingDate, err := timex.FormateTimeToStrWithUsStyle(in.ClosingDate)

		if err != nil {
			return nil, nil, nil, nil, errors.Wrap(err, "invalid closing date")
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
			return nil, nil, nil, nil, errors.Wrap(err, "invalid maturity date")
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
			return nil, nil, nil, nil, errors.Wrap(err, "invalid paid off date")
		}
		loan.PaidOffDate = &sql.NullTime{
			Time:  paidOffDate,
			Valid: true,
		}
		loan.PrincipalAmount = in.PrincipalAmount
	default:
		return nil, nil, nil, nil, errors.WithStack(errors.New("invalid loan stage"))
	}
	return
}
