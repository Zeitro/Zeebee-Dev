package loanlogic

import (
	"context"
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
	if loanModel.Source != constants.LoanSourceBorrower && loanModel.Source != constants.LoanSourcePreApproval {
		return nil, errors.WithStack(errors.New("Loan Source Is Not 3rd Party"))
	}
	loanPropertyModel := loanModel.LoanProperty
	if loanPropertyModel == nil {
		return nil, errors.WithStack(errors.New("Loan Property Is Not Found"))
	}

	// create loan
	err = l.svcCtx.LoanRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
		var innerErr error
		previousPropertyBaseId := loanPropertyModel.BorrowerPropertyId

		// update loan
		loan.Id = in.Data.Id
		loan.RemoveRelationships()
		// loan.PropertyId = property.Id
		_, innerErr = l.svcCtx.LoanRepo.UpdateOneWithTx(l.ctx, tx, in.Data.Id, loan, false)
		if innerErr != nil {
			return innerErr
		}

		// update loan property
		loanProperty.BorrowerPropertyId = in.Data.LoanProperty.BorrowerPropertyId
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
		removedApplicants := slice.DifferenceWith(oldApplicants, applicants, func(item1, item2 *model.LoanApplicant) bool {
			return item1.Id == item2.Id
		})
		removedApplicantIds := slice.Map(removedApplicants, func(index int, item *model.LoanApplicant) int64 {
			return item.Id
		})
		newApplicants := slice.DifferenceWith(applicants, oldApplicants, func(item1, item2 *model.LoanApplicant) bool {
			return item1.Id == item2.Id
		})
		newApplicants = slice.Map(newApplicants, func(index int, newApplicant *model.LoanApplicant) *model.LoanApplicant {
			newApplicant.LoanId = loan.Id
			newApplicant.LoanPropertyId = loanPropertyModel.Id
			newApplicant.LoanBorrowerId = loanBorrower.Id
			return newApplicant
		})

		updatedApplicants := slice.Filter(applicants, func(index int, newApplicant *model.LoanApplicant) bool {
			_, ok := slice.FindBy(oldApplicants, func(index int, oldApplicant *model.LoanApplicant) bool {
				return newApplicant.Id == oldApplicant.Id
			})
			return ok
		})

		updatedApplicants = slice.Map(updatedApplicants, func(index int, updateApplicant *model.LoanApplicant) *model.LoanApplicant {
			updateApplicant.LoanId = loan.Id
			updateApplicant.LoanPropertyId = loanPropertyModel.Id
			updateApplicant.LoanBorrowerId = loanBorrower.Id
			return updateApplicant
		})

		if len(removedApplicantIds) > 0 {
			_, innerErr = l.svcCtx.LoanApplicantRepo.DeleteByIDsWithTx(l.ctx, tx, removedApplicantIds)
			if innerErr != nil {
				return innerErr
			}
		}

		if len(newApplicants) > 0 {
			_, innerErr = l.svcCtx.LoanApplicantRepo.BatchCreateWithTx(l.ctx, tx, newApplicants, false)
			if innerErr != nil {
				return innerErr
			}
		}

		for _, updatedApplicant := range updatedApplicants {
			_, innerErr = l.svcCtx.LoanApplicantRepo.UpdateOneWithTx(l.ctx, tx, updatedApplicant.Id, updatedApplicant)
			if innerErr != nil {
				return innerErr
			}
		}

		// sync property total amount
		_ = l.svcCtx.BorrowerPropertyRepo.SyncTotalAmountAndLoanCountWithTx(l.ctx, tx, previousPropertyBaseId)
		if in.Data.LoanProperty.BorrowerPropertyId != previousPropertyBaseId {
			_ = l.svcCtx.BorrowerPropertyRepo.SyncTotalAmountAndLoanCountWithTx(l.ctx, tx, in.Data.LoanProperty.BorrowerPropertyId)
		}
		return innerErr
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to update loan")
	}

	return &borrower.EmptyResponse{}, nil
}

// handle property
func (l *UpdateBorrowerLoanLogic) handleBorrowerProperty(tx *gorm.DB, propertyId int64, property *model.BorrowerProperty) (newProperty *model.BorrowerProperty, err error) {
	// get existed property
	newProperty, err = l.svcCtx.BorrowerPropertyRepo.GetOneWithTxByProperty(l.ctx, tx, property)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get property base")
	}
	if newProperty != nil {
		return newProperty, nil
	}

	// get previous property
	prevProperty, err := l.svcCtx.BorrowerPropertyRepo.GetOneWithTx(l.ctx, tx, propertyId, []string{"LoanProperties"})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get property")
	}

	// create new property if relate zero or more than one loan property
	if len(prevProperty.LoanProperties) == 0 || len(prevProperty.LoanProperties) > 1 {
		return l.svcCtx.BorrowerPropertyRepo.GetOrCreateWithTx(l.ctx, tx, property)
	}

	// update property when relate one loan property
	property.Id = prevProperty.Id
	newProperty, err = l.svcCtx.BorrowerPropertyRepo.UpdateOneWithTx(l.ctx, tx, property.Id, property)
	return
}

func (l *UpdateBorrowerLoanLogic) buildUpdateModels(in *borrower.LoanUpdateRequestData) (loan *model.Loan, loanProperty *model.LoanProperty, applicants []*model.LoanApplicant, loanBorrower *model.LoanBorrower, err error) {
	// build applicants
	for _, applicant := range in.Applicants {
		if len(applicant.FirstName) == 0 {
			return nil, nil, nil, nil, errors.New("Applicant Name Is Required")
		}
		applicantName := strings.TrimSpace(applicant.FirstName)
		if len(applicantName) == 0 {
			return nil, nil, nil, nil, errors.New("Applicant Name Is Empty")
		}
		var applicantItem *model.LoanApplicant
		applicantItem = &model.LoanApplicant{}

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
	slice.SortBy(applicants, func(a, b *model.LoanApplicant) bool {
		return b.IsPrimary
	})
	applicants = slice.UniqueByComparator(applicants, func(item *model.LoanApplicant, other *model.LoanApplicant) bool {
		return item.FullName == other.FullName
	})
	if len(applicants) == 0 {
		return nil, nil, nil, nil, errors.New("No Applicant")
	}
	primaryCount := slice.CountBy(applicants, func(index int, item *model.LoanApplicant) bool {
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
	if err != nil {
		return nil, nil, nil, nil, errorx.GetFailedTransErr(err)
	}
	// loan property
	borrowerProperty, err := l.svcCtx.BorrowerPropertyRepo.FindOneByColumns(l.ctx, map[string]interface{}{
		"id":          in.LoanProperty.BorrowerPropertyId,
		"borrower_id": in.BorrowerId,
	})
	if err != nil {
		return nil, nil, nil, nil, err
	}

	loanProperty = &model.LoanProperty{}

	err = copier.CopyWithOption(loanProperty, borrowerProperty, copier.Option{
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
	loanProperty.BorrowerPropertyId = in.LoanProperty.BorrowerPropertyId
	loanProperty.Owner = in.LoanProperty.Owner
	loanProperty.PurchaseAmount = in.LoanProperty.PurchaseAmount

	// loanProperty = &model.LoanProperty{
	// 	BorrowerPropertyId:   in.LoanProperty.BorrowerPropertyId,
	// 	Unit:                 propertyBase.Unit,
	// 	PropertyType:         propertyBase.PropertyType,
	// 	LivingAreaSqft:       propertyBase.LivingAreaSqft,
	// 	Bedrooms:             propertyBase.Bedrooms,
	// 	Bathrooms:            propertyBase.Bathrooms,
	// 	FuturePropertyType:   propertyBase.PropertyType,
	// 	FutureUnit:           propertyBase.Unit,
	// 	FutureLivingAreaSqft: propertyBase.LivingAreaSqft,
	// 	FutureBedrooms:       propertyBase.Bedrooms,
	// 	FutureBathrooms:      propertyBase.Bathrooms,
	// 	Address:              propertyBase.Address,
	// 	Owner:                in.LoanProperty.Owner,
	// 	PurchaseAmount:       in.LoanProperty.PurchaseAmount,
	// }
	// if len(in.LoanProperty.PurchaseDate) > 0 {
	// 	purchaseDate, err := timex.FormateTimeToStrWithUsStyle(in.LoanProperty.PurchaseDate)
	// 	if err != nil {
	// 		return nil, nil, nil, nil, errors.Wrap(err, "Invalid Purchase Date")
	// 	}
	// 	loanProperty.PurchaseDate = &sql.NullTime{
	// 		Time:  purchaseDate,
	// 		Valid: true,
	// 	}
	//
	// }

	if loan.TotalAmount < loan.InitialAmount {
		return nil, nil, nil, nil, errors.New("Total Amount Cannot Be Less Than Initial Amount")
	}
	switch constants.LoanStage(in.Stage) {
	case constants.LoanStagePending:
		if loan.ClosingDate == nil {
			return nil, nil, nil, nil, errors.New("Closing Date Is Required")
		}
		loan.HoldbackAmount = loan.TotalAmount - loan.InitialAmount
		loan.PrincipalAmount = loan.TotalAmount
		// loan.PendingStatus = constants.LoanPendingStatus(in.PendingStatus)

	case constants.LoanStageActive:

		// loan.PrincipalAmount = in.PrincipalAmount
		if loan.MaturityDate == nil {
			return nil, nil, nil, nil, errors.New("Maturity Date Is Required")
		}
	case constants.LoanStageArchived:
		// loan.PrincipalAmount = in.PrincipalAmount
		if loan.PaidOffDate == nil {
			return nil, nil, nil, nil, errors.New("Paid Off Date Is Required")
		}
	default:
		return nil, nil, nil, nil, errors.WithStack(errors.New("Invalid Loan Stage"))
	}
	return
}
