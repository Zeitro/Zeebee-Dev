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

	loan, loanProperty, applicants, loanBorrower, err := l.buildUpdateModels(in.Data)
	if err != nil {
		return nil, err
	}

	oldApplicants := loanModel.Applicants
	if len(oldApplicants) == 0 {
		return nil, errors.WithStack(errors.New("Applicant Is Not Found"))
	}
	if loanModel.Source != constants.LoanSourceBorrower {
		return nil, errors.WithStack(errors.New("Loan Source Is Not 3rd Party"))
	}
	loanPropertyModel := loanModel.LoanProperty
	if loanPropertyModel == nil {
		return nil, errors.WithStack(errors.New("Loan Property Is Not Found"))
	}

	// create loan
	err = l.svcCtx.LoanRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
		var innerErr error
		previousPropertyBaseId := loanPropertyModel.PropertyBaseId

		// update loan
		loan.Id = in.Data.Id
		loan.RemoveRelationships()
		// loan.PropertyId = property.Id
		_, innerErr = l.svcCtx.LoanRepo.UpdateOneWithTx(l.ctx, tx, in.Data.Id, loan, false)
		if innerErr != nil {
			return innerErr
		}

		// update loan property
		loanProperty.PropertyBaseId = in.Data.LoanProperty.PropertyBaseId
		_, innerErr = l.svcCtx.LoanPropertyRepo.UpdateOneWithTx(l.ctx, tx, loanPropertyModel.Id, loanProperty)
		if innerErr != nil {
			return innerErr
		}

		// loan borrower
		if loanModel.LoanBorrower == nil {
			loanBorrower.LoanId = loanModel.Id
			loanBorrower, innerErr = l.svcCtx.LoanBorrowerRepo.CreateWithTx(l.ctx, tx, loanBorrower, false)
			if innerErr != nil {
				return innerErr
			}
		} else {
			loanBorrower.LoanId = loanModel.Id
			loanBorrower.Id = loanModel.LoanBorrower.Id
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
		_ = l.svcCtx.PropertyBaseRepo.SyncTotalAmountAndLoanCountWithTx(l.ctx, tx, previousPropertyBaseId)
		if in.Data.LoanProperty.PropertyBaseId != previousPropertyBaseId {
			_ = l.svcCtx.PropertyBaseRepo.SyncTotalAmountAndLoanCountWithTx(l.ctx, tx, in.Data.LoanProperty.PropertyBaseId)
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

func (l *UpdateBorrowerLoanLogic) buildUpdateModels(in *borrower.LoanUpdateRequestData) (loan *model.Loan, loanProperty *model.LoanProperty, applicants []*model.Applicant, loanBorrower *model.LoanBorrower, err error) {
	// build applicants
	for _, applicant := range in.Applicants {
		if len(applicant.FirstName) == 0 {
			return nil, nil, nil, nil, errors.New("Applicant Name Is Required")
		}
		applicantName := strings.TrimSpace(applicant.FirstName)
		if len(applicantName) == 0 {
			return nil, nil, nil, nil, errors.New("Applicant Name Is Empty")
		}
		var applicantItem *model.Applicant
		applicantItem = &model.Applicant{}

		err = copier.CopyWithOption(&applicantItem, &applicant, copier.Option{
			IgnoreEmpty:   false,
			CaseSensitive: true,
			DeepCopy:      true,
		})
		if err != nil {
			return nil, nil, nil, nil, errorx.GetFailedTransErr(err)
		}

		if !applicantItem.Type.IsValid() {
			return nil, nil, nil, nil, errors.New("Invalid Applicant Type")
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
		return nil, nil, nil, nil, errors.New("No Applicant")
	}
	primaryCount := slice.CountBy(applicants, func(index int, item *model.Applicant) bool {
		return item.IsPrimary
	})
	if primaryCount == 0 {
		return nil, nil, nil, nil, errors.New("No Primary Applicant")
	}
	if primaryCount > 1 {
		return nil, nil, nil, nil, errors.New("Multiple Primary Applicant")
	}

	// loan borrower
	loanBorrower = &model.LoanBorrower{
		Name:       applicants[0].FullName,
		Type:       applicants[0].Type,
		OwnPercent: 100,
	}

	loan = &model.Loan{}
	err = copier.CopyWithOption(loan, in, copier.Option{
		IgnoreEmpty:   false,
		DeepCopy:      true,
		CaseSensitive: true,
		Converters: []copier.TypeConverter{
			copierx.StrToSqlTimePtr(timex.LayoutUsStyle),
		},
	})
	// loan property
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
		loan.MonthlyPayment = in.MonthlyPayment
		loan.LienPosition = in.LienPosition
		loan.ServiceCompany = in.ServiceCompany
		loan.PrepaymentPenalty = in.PrepaymentPenalty
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
