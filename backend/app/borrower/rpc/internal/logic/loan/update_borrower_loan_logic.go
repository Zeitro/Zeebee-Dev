package loanlogic

import (
	"context"
	"database/sql"
	"strings"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/rpc/internal/repo/model"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/constants"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/stringx"
	"zeebee-backend/common/timex"

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

	loanModel, err := l.svcCtx.LoanRepo.GetBorrowerLoan(l.ctx, in.Data.Id, in.Data.BorrowerId)
	if err != nil {
		return nil, err
	}

	loan, propertyBase, loanProperty, applicants, loanBorrower, err := l.buildUpdateModels(in.Data)
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
		_, innerErr = l.svcCtx.LoanRepo.UpdateOneWithTx(l.ctx, tx, in.Data.Id, loan, false)
		if innerErr != nil {
			return innerErr
		}

		// update loan property
		loanProperty.PropertyBaseId = propertyBase.Id
		_, innerErr = l.svcCtx.LoanPropertyRepo.UpdateOneWithTx(l.ctx, tx, loanPropertyModel.Id, loanProperty)
		// innerErr = l.svcCtx.LoanPropertyRepo.UpdateColumnsByIdWithTx(l.ctx, tx, loanPropertyModel.Id, map[string]interface{}{
		// 	"property_base_id": propertyBase.Id,
		// 	"property_type":    loanProperty.PropertyType,
		// })
		if innerErr != nil {
			return innerErr
		}

		// loan borrower
		if loanModel.LoanBorrower == nil {
			loanBorrower.LoanId = loan.Id
			loanBorrower, innerErr = l.svcCtx.LoanBorrowerRepo.CreateWithTx(l.ctx, tx, loanBorrower, false)
			if innerErr != nil {
				return innerErr
			}
		} else {
			loanBorrower.LoanId = loan.Id
			loanBorrower.Id = loan.LoanBorrower.Id
			_, innerErr = l.svcCtx.LoanBorrowerRepo.UpdateOneWithTx(l.ctx, tx, loanBorrower.Id, loanBorrower)
			if innerErr != nil {
				return innerErr
			}

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
			newApplicant.LoanBorrowerId = loanBorrower.Id
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
			updateApplicant.LoanBorrowerId = loanBorrower.Id
			return updateApplicant
		})

		if len(removedApplicantIds) > 0 {
			_, innerErr = l.svcCtx.ApplicantRepo.DeleteByIDsWithTx(l.ctx, tx, removedApplicantIds)
			if innerErr != nil {
				return innerErr
			}
		}

		if len(newApplicants) > 0 {
			_, innerErr = l.svcCtx.ApplicantRepo.BatchCreateWithTx(l.ctx, tx, newApplicants, false)
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

func (l *UpdateBorrowerLoanLogic) buildUpdateModels(in *borrower.LoanUpdateRequestData) (loan *model.Loan, propertyBase *model.PropertyBase, loanProperty *model.LoanProperty, applicants []*model.Applicant, loanBorrower *model.LoanBorrower, err error) {
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
			// BorrowerId: in.BorrowerId,
			// // Name:      applicantName,
			// FirstName: applicant.FirstName,
			// LastName:  applicant.LastName,
			// MidName:   applicant.MidName,
			// FullName:  stringx.FormatFullName(applicant.FirstName, applicant.MidName, applicant.LastName),
			// IsPrimary: applicant.IsPrimary,
			// Type:      constants.ApplicantType(applicant.Type),
		}

		err = copier.CopyWithOption(&applicantItem, &applicant, copier.Option{
			IgnoreEmpty:   true,
			CaseSensitive: true,
			DeepCopy:      true,
		})
		if err != nil {
			return nil, nil, nil, nil, nil, errorx.GetFailedTransErr(err)
		}

		if !applicantItem.Type.IsValid() {
			return nil, nil, nil, nil, nil, errors.New("invalid applicant type")
		}

		applicantItem.Id = applicant.Id
		applicantItem.FullName = stringx.FormatFullName(applicantItem.FirstName, applicantItem.MidName, applicantItem.LastName)
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
	// loan
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
	// }
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
	loan.RemoveRelationships()

	// property base
	propertyBase = &model.PropertyBase{}
	err = copier.CopyWithOption(propertyBase, in.LoanProperty.PropertyBase, copier.Option{
		IgnoreEmpty:   false,
		CaseSensitive: true,
		DeepCopy:      true,
		Converters:    nil,
	})
	if err != nil {
		return nil, nil, nil, nil, nil, errorx.GetFailedTransErr(err)
	}
	propertyBase.BorrowerId = in.BorrowerId

	// propertyBase = &model.PropertyBase{
	// 	BorrowerId: in.BorrowerId,
	// 	// PropertyType: constants.PropertyType(in.Property.PropertyType),
	// 	Unit: in.LoanProperty.PropertyBase.Unit,
	// 	Address: types.Address{
	// 		State:    in.LoanProperty.PropertyBase.Address.State,
	// 		City:     in.LoanProperty.PropertyBase.Address.City,
	// 		Address1: in.LoanProperty.PropertyBase.Address.Address1,
	// 		Address2: in.LoanProperty.PropertyBase.Address.Address2,
	// 		ZipCode:  in.LoanProperty.PropertyBase.Address.ZipCode,
	// 	},
	// }
	loanProperty = &model.LoanProperty{
		// Owner:        in.LoanProperty.Ower,
		// PropertyType: constants.PropertyType(in.LoanProperty.PropertyType),
	}
	err = copier.CopyWithOption(loanProperty, in.LoanProperty, copier.Option{
		CaseSensitive: true,
		DeepCopy:      true,
	})
	if err != nil {
		return nil, nil, nil, nil, nil, errorx.GetFailedTransErr(err)
	}
	loanProperty.RemoveRelationships()
	loanProperty.Unit = propertyBase.Unit

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
		loan.MonthlyPayment = in.MonthlyPayment
		loan.LienPosition = in.LienPosition
		loan.ServiceCompany = in.ServiceCompany
		loan.PrepaymentPenalty = in.PrepaymentPenalty
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
