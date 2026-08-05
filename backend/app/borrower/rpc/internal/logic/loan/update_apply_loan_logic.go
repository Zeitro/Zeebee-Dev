package loanlogic

import (
	"context"

	"github.com/duke-git/lancet/v2/slice"
	"gorm.io/gorm"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/proto"
)

type UpdateApplyLoanLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewUpdateApplyLoanLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateApplyLoanLogic {
	return &UpdateApplyLoanLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *UpdateApplyLoanLogic) UpdateApplyLoan(in *borrower.LoanApplyRequest) (resp *proto.EmptyResponse, err error) {

	var loan *model.Loan
	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx, constants.UserTypeBorrower)
	if err != nil {
		return nil, err
	}

	// var borrowerProperty *model.BorrowerProperty

	loan, err = l.svcCtx.LoanRepo.GetBorrowerLoan(l.ctx, in.Data.Id, claim.OrganizationId)
	if err != nil {
		return nil, err
	}

	applyLogic := NewApplyLoanLogic(l.ctx, l.svcCtx)
	data, err := applyLogic.buildModelData(in.Data, loan, claim)
	if err != nil {
		return nil, err
	}

	sources := []constants.LoanSource{
		constants.LoanSourceSystem,
		constants.LoanSourcePreapproval,
	}
	if !slice.Contain(sources, loan.Source) {
		return nil, errorx.WithStack(errorx.New("invalid loan source"))
	}
	stages := []constants.LoanStage{
		constants.LoanStagePending,
		constants.LoanStageDraft,
	}
	pendingStatus := []constants.LoanPendingStatus{
		constants.LoanPendingStatusQuoting,
		constants.LoanPendingStatusNone,
	}
	if !slice.Contain(stages, loan.Stage) && !slice.Contain(pendingStatus, loan.PendingStatus) {
		return nil, errorx.New("loan's stage should be pending or draft")
	}

	if loan.BorrowerOrganizationId != claim.OrganizationId {
		return nil, errorx.WithStack(errorx.New("loan does not belong to borrower"))
	}

	oldApplicants := loan.Applicants
	// if len(oldApplicants) == 0 {
	// 	return nil, errors.WithStack(errors.New("Applicant Is Not Found"))
	// }

	// update loan
	err = l.svcCtx.LoanPropertyRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
		var innerErr error
		var previousBorrowerPropertyId int64
		if loan.LoanProperty != nil {
			previousBorrowerPropertyId = loan.LoanProperty.BorrowerPropertyId
		}
		// previousPropertyBaseId := data.loanProperty.PropertyBaseId

		eg := errorx.NewErrorGroup()
		eg.Do(
			// update loan
			func() error {
				data.loan.RemoveRelationships()
				if in.Data.LenderCompanyId == 0 {
					in.Data.LenderCompanyId = loan.LenderOrganizationId
				}
				_, err := l.svcCtx.LoanRepo.UpdateOneWithTx(l.ctx, tx, in.Data.Id, data.loan)
				return err
			},
			// update loan property
			func() error {
				property := data.loanProperty
				property.LoanId = loan.Id
				// property.BorrowerPropertyId = borrowerProperty.Id

				if property.Id == 0 {
					property, err = l.svcCtx.LoanPropertyRepo.CreateWithTx(l.ctx, tx, property, false)
				} else {
					property, err = l.svcCtx.LoanPropertyRepo.UpdateOneWithTx(l.ctx, tx, property.Id, property)
				}
				if err == nil {
					loan.LoanProperty = property
				}
				return err
			},
			// loan borrower
			func() error {
				var err error
				var loanBorrower *model.LoanBorrower
				data.loanBorrower.LoanId = loan.Id
				if loan.LoanBorrower == nil {
					loanBorrower, err = l.svcCtx.LoanBorrowerRepo.CreateWithTx(l.ctx, tx, data.loanBorrower, false)
				} else {
					loanBorrower, err = l.svcCtx.LoanBorrowerRepo.UpdateOneWithTx(l.ctx, tx, loan.LoanBorrower.Id, data.loanBorrower)
				}
				if err == nil {
					loan.LoanBorrower = loanBorrower
				}
				return err
			},
			// handle applicant
			func() error {
				return l.updateApplicants(tx, oldApplicants, data.applicants, loan)
			},

			// loan demographic
			func() error {
				var err error
				data.loanDemographic.LoanId = loan.Id
				if loan.LoanDemographic == nil {
					_, err = l.svcCtx.LoanDemographicRepo.CreateWithTx(l.ctx, tx, data.loanDemographic, false)
				} else {
					_, err = l.svcCtx.LoanDemographicRepo.UpdateOneWithTx(l.ctx, tx, loan.LoanDemographic.Id, data.loanDemographic)
				}
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
				var err error
				if len(in.Data.ApplyFormFile) > 0 {
					err = l.svcCtx.LoanDocumentRepo.SaveApplyFormDocumentWithTx(l.ctx, tx, loan.Id, in.Data.ApplyFormFile)
				}
				return err
			},
			// Create loan documents with pending upload status
			func() error {
				var err error
				switch loan.Stage {
				case constants.LoanStagePending:
					return err
				case constants.LoanStageDraft:
					// change stage from draft to pending
				default:
					err = errorx.WithStack(errorx.New("Invalid loan stage"))
					return err
				}
				loanDocumentTypes := constants.GetLoanDocumentTypes()
				var loanDocuments []*model.LoanDocument
				for category, fileTypes := range loanDocumentTypes {
					documentStatus := constants.LoanDocumentStatusPendingUpload
					if category == constants.LoanDocumentCategoryBorrower || category == constants.LoanDocumentCategoryExtend {
						documentStatus = constants.LoanDocumentStatusPendingRequested
					}
					for _, fileType := range fileTypes {
						loanDocument := &model.LoanDocument{
							LoanId:   loan.Id,
							Category: category,
							FileType: fileType,
							Status:   documentStatus,
						}
						loanDocuments = append(loanDocuments, loanDocument)
					}
				}

				if len(loanDocuments) > 0 {
					_, err = l.svcCtx.LoanDocumentRepo.BatchCreateWithTx(l.ctx, tx, loanDocuments)
					if err != nil {
						return err
					}
				}
				return err
			},
			// sync property total amount and loan count
			func() error {
				var err error
				err = l.svcCtx.BorrowerPropertyRepo.SyncTotalAmountAndLoanCountWithTx(l.ctx, tx, previousBorrowerPropertyId)
				if err != nil {
					return err
				}
				if in.Data.LoanProperty.BorrowerPropertyId != previousBorrowerPropertyId {
					err = l.svcCtx.BorrowerPropertyRepo.SyncTotalAmountAndLoanCountWithTx(l.ctx, tx, in.Data.LoanProperty.BorrowerPropertyId)
				}
				return err
			},
		)

		innerErr = eg.Wait()
		return innerErr
	})

	return &proto.EmptyResponse{}, err
}

func (l *UpdateApplyLoanLogic) updateApplicants(tx *gorm.DB, oldItems []*model.LoanApplicant, curItems []*model.LoanApplicant, loan *model.Loan) (err error) {
	slice.ForEach(curItems, func(index int, item *model.LoanApplicant) {
		curItems[index].LoanId = loan.Id
		if loan.LoanProperty != nil && loan.LoanProperty.Id > 0 {
			curItems[index].LoanPropertyId = loan.LoanProperty.Id
		}
		if loan.LoanBorrower != nil && loan.LoanBorrower.Id > 0 {
			curItems[index].LoanBorrowerId = loan.LoanBorrower.Id
		}
	})

	removedItems := slice.DifferenceWith(oldItems, curItems, func(item1, item2 *model.LoanApplicant) bool {
		return item1.Id == item2.Id
	})
	removedIds := slice.Map(removedItems, func(index int, item *model.LoanApplicant) int64 {
		return item.Id
	})
	newItems := slice.DifferenceWith(curItems, oldItems, func(item1, item2 *model.LoanApplicant) bool {
		return item1.Id == item2.Id
	})

	updatedItems := slice.Filter(curItems, func(index int, newApplicant *model.LoanApplicant) bool {
		_, ok := slice.FindBy(oldItems, func(index int, oldApplicant *model.LoanApplicant) bool {
			return newApplicant.Id == oldApplicant.Id
		})
		return ok
	})

	if len(removedIds) > 0 {
		_, err = l.svcCtx.LoanApplicantRepo.DeleteByIDsWithTx(l.ctx, tx, removedIds)
		if err != nil {
			return err
		}
	}

	if len(newItems) > 0 {
		_, err = l.svcCtx.LoanApplicantRepo.BatchCreateWithTx(l.ctx, tx, newItems, false)
		if err != nil {
			return err
		}
	}

	for _, updatedItem := range updatedItems {
		_, err = l.svcCtx.LoanApplicantRepo.UpdateOneWithTx(l.ctx, tx, updatedItem.Id, updatedItem)
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
