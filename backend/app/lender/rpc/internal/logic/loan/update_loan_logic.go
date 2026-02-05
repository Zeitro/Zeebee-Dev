package loanlogic

import (
	"context"

	"github.com/duke-git/lancet/v2/mathutil"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/timex"
	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type updateLoanModelData struct {
	loan *model.Loan
	// loanBorrower      *model.LoanBorrower
	loanProperty           *model.LoanProperty
	loanDemographic        *model.LoanDemographic
	loanThirdServices      []*model.LoanThirdService
	applicants             []*model.LoanApplicant
	transactionExperiences []*model.LoanExperience
	realEstatesExperiences []*model.LoanExperience
	// applyFormFile          *model.File
}

type UpdateLoanLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewUpdateLoanLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLoanLogic {
	return &UpdateLoanLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.NewBaseLogic(ctx, svcCtx),
	}
}

// UpdateLoan updates a loan
func (l *UpdateLoanLogic) UpdateLoan(in *lender.LoanUpdateRequest) (resp *lender.LoanUpdateResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	var oldLoan *model.Loan
	var oldLoanCopy *model.Loan
	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	// old loan
	oldLoan, err = l.svcCtx.LoanRepo.GetLenderLoan(l.ctx, in.Data.Id, claim.OrganizationId, true)
	if err != nil {
		return nil, err
	}
	if oldLoan.Source != constants.LoanSourceSystem && oldLoan.Source != constants.LoanSourcePreapproval {
		return nil, errorx.WithStack(errorx.New("Invalid loan source"))
	}

	// oldApplicants := oldLoan.Applicants

	oldLoanCopy, err = utils.Clone(oldLoan)
	if err != nil {
		return nil, errorx.Wrap(err, "Failed To Get Copy Loan")
	}

	data, err := l.buildModelData(in.Data, oldLoanCopy)
	if err != nil {
		return nil, err
	}

	// update loan
	err = l.svcCtx.LoanPropertyRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
		var innerErr error
		// previousPropertyBaseId := data.loanProperty.PropertyBaseId

		eg := errorx.NewErrorGroup()
		eg.Do(
			// update loan
			func() error {
				data.loan.RemoveRelationships()
				// if in.Data.OrganizationId == 0 {
				// 	in.Data.OrganizationId = oldLoan.LenderOrganizationId
				// }
				_, err := l.svcCtx.LoanRepo.UpdateOneWithTx(l.ctx, tx, in.Data.Id, data.loan)
				return err
			},
			// update loan property
			func() error {
				property := data.loanProperty
				property.LoanId = oldLoan.Id
				// property.PropertyBaseId = propertyBase.Id

				if property.Id == 0 {
					property, err = l.svcCtx.LoanPropertyRepo.CreateWithTx(l.ctx, tx, property, false)
				} else {
					property, err = l.svcCtx.LoanPropertyRepo.UpdateOneWithTx(l.ctx, tx, property.Id, property)
				}
				if err == nil {
					oldLoan.LoanProperty = property
				}
				return err

			},
			// update Applicants
			func() error {
				var err error
				err = l.updateApplicants(tx, data.applicants)
				return err
			},
			// loan demographic
			func() error {
				var err error
				data.loanDemographic.LoanId = oldLoan.Id
				if oldLoan.LoanDemographic == nil {
					_, err = l.svcCtx.LoanDemographicRepo.CreateWithTx(l.ctx, tx, data.loanDemographic, false)
				} else {
					_, err = l.svcCtx.LoanDemographicRepo.UpdateOneWithTx(l.ctx, tx, oldLoan.LoanDemographic.Id, data.loanDemographic)
				}
				return err
			},
			// loan experience service
			func() error {
				oldExperiences := oldLoan.LoanExperiences
				transactionExperiences := data.transactionExperiences
				realEstatesExperiences := data.realEstatesExperiences
				allExperiences := append(transactionExperiences, realEstatesExperiences...)
				return l.updateLoanExperiences(tx, oldExperiences, allExperiences, oldLoan)
			},
			// loan third service
			func() error {
				return l.updateThirdServices(tx, oldLoan.LoanThirdServices, data.loanThirdServices, oldLoan)
			},
			// sync loan total amount
			func() error {
				var err error
				if oldLoan.LoanProperty != nil {
					err = l.svcCtx.BorrowerPropertyRepo.SyncTotalAmountAndLoanCountWithTx(l.ctx, tx, oldLoan.LoanProperty.BorrowerPropertyId)
				}
				return err
			},

			// loan operation log
			func() error {
				operatorId := claim.OrganizationUserId

				fromStatus := model.LoanOperationStatus{
					Stage:         oldLoan.Stage,
					PendingStatus: oldLoan.PendingStatus,
				}

				toStatus := model.LoanOperationStatus{
					Stage:         data.loan.Stage,
					PendingStatus: data.loan.PendingStatus,
				}

				loanBaseLogic := NewLoanBaseLogic(l.ctx, l.svcCtx)
				return loanBaseLogic.CreateLoanOperationLog(tx, oldLoan.Id, operatorId, fromStatus, toStatus, constants.LoanOperationActionUpdate, nil)
			},
		)

		innerErr = eg.Wait()
		return innerErr
	})
	if err != nil {
		return nil, err
	}

	return &lender.LoanUpdateResponse{
		Data: &lender.LoanUpdateResponseData{
			Id: oldLoan.Id,
		},
	}, nil
}

// updateApplicants updates loan applicants
func (l *UpdateLoanLogic) updateApplicants(tx *gorm.DB, applicants []*model.LoanApplicant) (err error) {
	for _, item := range applicants {
		_, err = l.svcCtx.LoanApplicantRepo.UpdateOneWithTx(l.ctx, tx, item.Id, item)
		if err != nil {
			return err
		}
	}

	return err
}

// updateLoanExperiences updates loan experiences
func (l *UpdateLoanLogic) updateLoanExperiences(tx *gorm.DB, oldItems []*model.LoanExperience, curItems []*model.LoanExperience, loan *model.Loan) (err error) {
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

// updateThirdServices updates loan third services
func (l *UpdateLoanLogic) updateThirdServices(tx *gorm.DB, oldItems []*model.LoanThirdService, curItems []*model.LoanThirdService, loan *model.Loan) (err error) {
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

// buildModelData builds model data
func (l *UpdateLoanLogic) buildModelData(reqData *lender.LoanUpdateRequestData, oldLoan *model.Loan) (data *updateLoanModelData, err error) {
	if oldLoan == nil {
		return nil, errorx.WithStack(errorx.New("Invalid loan"))
	}

	if reqData == nil {
		return nil, errorx.WithStack(errorx.New("Invalid request data"))
	}
	if reqData.LoanProperty == nil {
		return nil, errorx.WithStack(errorx.New("Invalid loan property"))
	}

	// loan
	loan := &model.Loan{}
	_ = copier.CopyWithOption(loan, oldLoan, copier.Option{
		IgnoreEmpty:   false,
		CaseSensitive: true,
		DeepCopy:      true,
	})

	err = copier.CopyWithOption(loan, reqData, copier.Option{
		CaseSensitive: true,
		DeepCopy:      true,
		Converters: []copier.TypeConverter{
			copierx.StrToSqlTimePtr(timex.LayoutUsStyle),
		},
	})
	if err != nil {
		return nil, errorx.GetFailedTransErr(err)
	}
	loan.RemoveRelationships()

	// LoanProperty
	loanProperty := &model.LoanProperty{}
	if oldLoan.LoanProperty != nil {
		_ = copier.CopyWithOption(loanProperty, oldLoan.LoanProperty, copier.Option{
			IgnoreEmpty:   false,
			DeepCopy:      true,
			CaseSensitive: true,
		})
	}

	err = copier.CopyWithOption(loanProperty, reqData.LoanProperty, copier.Option{
		IgnoreEmpty:   false,
		CaseSensitive: true,
		DeepCopy:      true,
		Converters: []copier.TypeConverter{
			copierx.StrToSqlTimePtr(timex.LayoutUsStyle),
		},
	})
	if err != nil {
		return nil, errorx.GetFailedTransErr(err)
	}
	if loanProperty.YearlyLiabilityExpenseAmount > 0 {
		dscr := (loanProperty.EstYearlyRentalIncomeAmount / loanProperty.YearlyLiabilityExpenseAmount) * 100
		loanProperty.Dscr = mathutil.RoundToFloat(dscr, 2)
	}

	// applicants
	applicants := make([]*model.LoanApplicant, 0)
	var totalOwnPercent int32
	for _, applicant := range reqData.Applicants {
		applicantModel, err := l.svcCtx.LoanApplicantRepo.FindOneById(l.ctx, applicant.Id)
		if err != nil {
			return nil, err
		}
		totalOwnPercent += applicant.OwnPercent
		applicantModel.OwnPercent = applicant.OwnPercent
		applicants = append(applicants, applicantModel)
	}
	if totalOwnPercent > 100 {
		return nil, errorx.WithStack(errorx.New("Total own percent greater than 100"))
	}

	// transaction experience
	transactionExperiences := make([]*model.LoanExperience, 0)
	err = copier.CopyWithOption(&transactionExperiences, &reqData.TransactionExperiences, copier.Option{
		IgnoreEmpty:   false,
		DeepCopy:      true,
		CaseSensitive: true,
		Converters: []copier.TypeConverter{
			copierx.StrToSqlTimePtr(timex.LayoutUsStyle),
		},
	})
	if err != nil {
		return nil, errorx.GetFailedTransErr(err)
	}

	// real_estates experience
	realEstatesExperiences := make([]*model.LoanExperience, 0)
	err = copier.CopyWithOption(&realEstatesExperiences, &reqData.RealEstateExperiences, copier.Option{
		IgnoreEmpty:   false,
		DeepCopy:      true,
		CaseSensitive: true,
		Converters: []copier.TypeConverter{
			copierx.StrToSqlTimePtr(timex.LayoutUsStyle),
		},
	})
	if err != nil {
		return nil, errorx.GetFailedTransErr(err)
	}

	// third services
	loanThirdServices := make([]*model.LoanThirdService, 0)
	err = copier.CopyWithOption(&loanThirdServices, &reqData.LoanThirdServices, copier.Option{
		IgnoreEmpty:   false,
		CaseSensitive: true,
		DeepCopy:      true,
		Converters: []copier.TypeConverter{
			copierx.StrToSqlTimePtr(timex.LayoutUsStyle),
		},
	})
	if err != nil {
		return nil, errorx.GetFailedTransErr(err)
	}

	// demographic
	loanDemographic := &model.LoanDemographic{}
	err = copier.CopyWithOption(&loanDemographic, &reqData.LoanDemographic, copier.Option{
		IgnoreEmpty:   false,
		CaseSensitive: true,
		DeepCopy:      true,
		Converters: []copier.TypeConverter{
			copierx.StrToSqlTimePtr(timex.LayoutUsStyle),
		},
	})
	if err != nil {
		return nil, errorx.GetFailedTransErr(err)
	}
	data = &updateLoanModelData{
		loan:                   loan,
		loanProperty:           loanProperty,
		applicants:             applicants,
		transactionExperiences: transactionExperiences,
		realEstatesExperiences: realEstatesExperiences,
		loanDemographic:        loanDemographic,
		loanThirdServices:      loanThirdServices,
	}
	return
}
