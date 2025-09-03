package loanlogic

import (
	"context"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"zeebee-backend/app/borrower/rpc/internal/repo/model"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/constants"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type UpdateApplyLoanLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateApplyLoanLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateApplyLoanLogic {
	return &UpdateApplyLoanLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateApplyLoanLogic) UpdateApplyLoan(in *borrower.LoanApplyRequest) (resp *borrower.EmptyResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	var loan *model.Loan
	var propertyBase *model.PropertyBase
	applyLogic := NewApplyLoanLogic(l.ctx, l.svcCtx)

	data, err := applyLogic.buildModelData(in.Data)
	if err != nil {
		return nil, err
	}
	loan, err = l.svcCtx.LoanRepo.GetBorrowerLoan(l.ctx, data.loan.Id, in.Data.BorrowerId)
	if err != nil {
		return nil, err
	}
	if loan.Source != constants.LoanSourceSystem {
		return nil, errors.WithStack(errors.New("invalid loan source"))
	}
	if loan.Stage != constants.LoanStagePending && loan.PendingStatus != constants.LoanPendingStatusQuoting {
		return nil, errors.New("loan's stage should be pending")
	}
	oldApplicants := loan.Applicants
	if len(oldApplicants) == 0 {
		return nil, errors.WithStack(errors.New("applicant is not found"))
	}

	// update loan
	err = l.svcCtx.LoanPropertyRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
		var innerErr error
		previousPropertyBaseId := data.loanProperty.PropertyBaseId

		eg := errorx.NewErrorGroup()
		eg.Do(
			// handle property
			func() error {
				updateLogic := NewUpdateBorrowerLoanLogic(l.ctx, l.svcCtx)
				propertyBase, err = updateLogic.handlePropertyBase(tx, loan.LoanProperty.PropertyBaseId, data.loanProperty.PropertyBase)
				return err
			},
			// update loan
			func() error {
				data.loan.RemoveRelationships()
				_, err := l.svcCtx.LoanRepo.UpdateOneWithTx(l.ctx, tx, in.Data.Id, data.loan, true)
				return err
			},
			// update loan property
			func() error {
				return l.svcCtx.LoanPropertyRepo.UpdateColumnsByIdWithTx(l.ctx, tx, data.loanProperty.Id, map[string]interface{}{
					"property_base_id": propertyBase.Id,
					"property_type":    data.loanProperty.PropertyType,
					"living_area_sqft": data.loanProperty.LivingAreaSqft,
					"bedrooms":         data.loanProperty.Bedrooms,
					"bathrooms":        data.loanProperty.Bathrooms,
				})
			},
			// handle applicant
			func() error {
				return l.updateApplicants(tx, oldApplicants, data.applicants, loan)
			},
			// loan borrower
			func() error {
				_, err := l.svcCtx.LoanBorrowerRepo.UpdateOneWithTx(l.ctx, tx, loan.LoanBorrower.Id, data.loanBorrower)
				return err
			},
			// loan demographic
			func() error {
				_, err := l.svcCtx.LoanDemographicRepo.UpdateOneWithTx(l.ctx, tx, loan.LoanDemographic.Id, data.loanDemographic)
				return err
			},
			// loan experience service
			func() error {
				oldExperiences := loan.LoanExperiences
				transactionExperiences := data.transactionExperiences
				realEstatesExperiences := data.realEstatesExperiences
				allExperiences := append(transactionExperiences, realEstatesExperiences...)
				return l.updateLoanExperiences(tx, oldExperiences, allExperiences, loan)
			},
			// loan third service
			func() error {
				return l.updateThirdServices(tx, loan.LoanThirdServices, data.loanThirdServices, loan)
			},
			// apply form
			func() error {
				if data.applyFormFile == nil {
					return nil
				}
				oldApplyForm := loan.ApplyFormFile
				if oldApplyForm == nil {
					data.applyFormFile.RefId = loan.Id
					_, err := l.svcCtx.FileRepo.CreateWithTx(l.ctx, tx, data.applyFormFile, false)
					return err
				}
				data.applyFormFile.RefId = loan.Id
				_, err := l.svcCtx.FileRepo.UpdateOneWithTx(l.ctx, tx, oldApplyForm.Id, data.applyFormFile)
				return err
			},
			// sync property total amount
			func() error {
				err := l.svcCtx.PropertyBaseRepo.SyncTotalAmountWithTx(l.ctx, tx, previousPropertyBaseId)
				if err != nil {
					return err
				}
				if propertyBase.Id != previousPropertyBaseId {
					err = l.svcCtx.PropertyBaseRepo.SyncTotalAmountWithTx(l.ctx, tx, propertyBase.Id)
				}
				return err
			},
		)

		innerErr = eg.Error()
		return innerErr
	})

	return &borrower.EmptyResponse{}, err
}

func (l *UpdateApplyLoanLogic) updateApplicants(tx *gorm.DB, oldItems []*model.Applicant, curItems []*model.Applicant, loan *model.Loan) (err error) {
	removedItems := slice.DifferenceWith(oldItems, curItems, func(item1, item2 *model.Applicant) bool {
		return item1.Id == item2.Id
	})
	removedIds := slice.Map(removedItems, func(index int, item *model.Applicant) int64 {
		return item.Id
	})
	newItems := slice.DifferenceWith(curItems, oldItems, func(item1, item2 *model.Applicant) bool {
		return item1.Id == item2.Id
	})
	newItems = slice.Map(newItems, func(index int, newApplicant *model.Applicant) *model.Applicant {
		newApplicant.LoanId = loan.Id
		newApplicant.LoanPropertyId = loan.LoanProperty.Id
		return newApplicant
	})

	updatedItems := slice.Filter(curItems, func(index int, newApplicant *model.Applicant) bool {
		_, ok := slice.FindBy(oldItems, func(index int, oldApplicant *model.Applicant) bool {
			return newApplicant.Id == oldApplicant.Id
		})
		return ok
	})

	updatedItems = slice.Map(updatedItems, func(index int, updateApplicant *model.Applicant) *model.Applicant {
		updateApplicant.LoanId = loan.Id
		updateApplicant.LoanPropertyId = loan.LoanProperty.Id
		return updateApplicant
	})

	if len(removedIds) > 0 {
		_, err = l.svcCtx.ApplicantRepo.DeleteByIDsWithTx(l.ctx, tx, removedIds)
		if err != nil {
			return err
		}
	}

	if len(newItems) > 0 {
		_, err = l.svcCtx.ApplicantRepo.BatchCreateWithTx(l.ctx, tx, newItems, false)
		if err != nil {
			return err
		}
	}

	for _, updatedItem := range updatedItems {
		_, err = l.svcCtx.ApplicantRepo.UpdateOneWithTx(l.ctx, tx, updatedItem.Id, updatedItem)
		if err != nil {
			return err
		}
	}

	return err
}

func (l *UpdateApplyLoanLogic) updateLoanExperiences(tx *gorm.DB, oldItems []*model.LoanExperience, curItems []*model.LoanExperience, loan *model.Loan) (err error) {
	removedItems := slice.DifferenceWith(oldItems, curItems, func(item1, item2 *model.LoanExperience) bool {
		return item1.Id == item2.Id
	})
	removedIds := slice.Map(removedItems, func(index int, item *model.LoanExperience) int64 {
		return item.Id
	})
	newItems := slice.DifferenceWith(curItems, oldItems, func(item1, item2 *model.LoanExperience) bool {
		return item1.Id == item2.Id
	})
	newItems = slice.Map(newItems, func(index int, item *model.LoanExperience) *model.LoanExperience {
		item.LoanId = loan.Id
		return item
	})

	updatedItems := slice.Filter(curItems, func(index int, newItem *model.LoanExperience) bool {
		_, ok := slice.FindBy(oldItems, func(index int, oldItem *model.LoanExperience) bool {
			return newItem.Id == oldItem.Id
		})
		return ok
	})

	updatedItems = slice.Map(updatedItems, func(index int, updateApplicant *model.LoanExperience) *model.LoanExperience {
		updateApplicant.LoanId = loan.Id
		return updateApplicant
	})

	if len(removedIds) > 0 {
		_, err = l.svcCtx.LoanExperienceRepo.DeleteByIDsWithTx(l.ctx, tx, removedIds)
		if err != nil {
			return err
		}
	}

	if len(newItems) > 0 {
		_, err = l.svcCtx.LoanExperienceRepo.BatchCreateWithTx(l.ctx, tx, newItems, false)
		if err != nil {
			return err
		}
	}

	for _, updatedItem := range updatedItems {
		_, err = l.svcCtx.LoanExperienceRepo.UpdateOneWithTx(l.ctx, tx, updatedItem.Id, updatedItem)
		if err != nil {
			return err
		}
	}

	return err
}

func (l *UpdateApplyLoanLogic) updateThirdServices(tx *gorm.DB, oldItems []*model.LoanThirdService, curItems []*model.LoanThirdService, loan *model.Loan) (err error) {
	removedItems := slice.DifferenceWith(oldItems, curItems, func(item1, item2 *model.LoanThirdService) bool {
		return item1.Id == item2.Id
	})
	removedIds := slice.Map(removedItems, func(index int, item *model.LoanThirdService) int64 {
		return item.Id
	})
	newItems := slice.DifferenceWith(curItems, oldItems, func(item1, item2 *model.LoanThirdService) bool {
		return item1.Id == item2.Id
	})
	newItems = slice.Map(newItems, func(index int, item *model.LoanThirdService) *model.LoanThirdService {
		item.LoanId = loan.Id
		return item
	})

	updatedItems := slice.Filter(curItems, func(index int, newItem *model.LoanThirdService) bool {
		_, ok := slice.FindBy(oldItems, func(index int, oldItem *model.LoanThirdService) bool {
			return newItem.Id == oldItem.Id
		})
		return ok
	})

	updatedItems = slice.Map(updatedItems, func(index int, updateApplicant *model.LoanThirdService) *model.LoanThirdService {
		updateApplicant.LoanId = loan.Id
		return updateApplicant
	})

	if len(removedIds) > 0 {
		_, err = l.svcCtx.LoanThirdServiceRepo.DeleteByIDsWithTx(l.ctx, tx, removedIds)
		if err != nil {
			return err
		}
	}

	if len(newItems) > 0 {
		_, err = l.svcCtx.LoanThirdServiceRepo.BatchCreateWithTx(l.ctx, tx, newItems, false)
		if err != nil {
			return err
		}
	}

	for _, updatedItem := range updatedItems {
		_, err = l.svcCtx.LoanThirdServiceRepo.UpdateOneWithTx(l.ctx, tx, updatedItem.Id, updatedItem)
		if err != nil {
			return err
		}
	}

	return err
}
