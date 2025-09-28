package loanlogic

import (
	"context"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"zeebee-backend/app/lender/rpc/internal/repo/model"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/constants"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/stringx"
	"zeebee-backend/common/timex"
	"zeebee-backend/common/types"
	"zeebee-backend/common/utils"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type updateLoanModelData struct {
	loan                   *model.Loan
	loanBorrower           *model.LoanBorrower
	loanProperty           *model.LoanProperty
	loanDemographic        *model.LoanDemographic
	loanThirdServices      []*model.LoanThirdService
	applicants             []*model.Applicant
	transactionExperiences []*model.LoanExperience
	realEstatesExperiences []*model.LoanExperience
	// applyFormFile          *model.File
}

type UpdateLoanLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateLoanLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLoanLogic {
	return &UpdateLoanLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
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
	var propertyBase *model.PropertyBase
	// old loan
	oldLoan, err = l.svcCtx.LoanRepo.GetLenderLoan(l.ctx, in.Data.Id, in.Data.LenderCompanyId, true)
	if err != nil {
		return nil, err
	}
	if oldLoan.Source != constants.LoanSourceSystem {
		return nil, errors.WithStack(errors.New("invalid loan source"))
	}

	oldApplicants := oldLoan.Applicants

	oldLoanCopy, err = utils.Clone(oldLoan)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get copy loan")
	}

	data, err := l.buildModelData(in.Data, oldLoanCopy)
	if err != nil {
		return nil, err
	}

	if len(data.applicants) == 0 {
		return nil, errors.WithStack(errors.New("applicant is not found"))
	}

	// update loan
	err = l.svcCtx.LoanPropertyRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
		var innerErr error
		previousPropertyBaseId := data.loanProperty.PropertyBaseId

		eg := errorx.NewErrorGroup()
		eg.Do(
			// handle property base
			func() error {
				propertyBase, err = l.handlePropertyBase(tx, oldLoan.LoanProperty.PropertyBaseId, data.loanProperty.PropertyBase)
				return err
			},
			// update loan
			func() error {
				data.loan.RemoveRelationships()
				_, err := l.svcCtx.LoanRepo.UpdateOneWithTx(l.ctx, tx, in.Data.Id, data.loan)
				return err
			},
			// update loan property
			func() error {
				property := data.loanProperty
				property.LoanId = oldLoan.Id
				property.PropertyBaseId = propertyBase.Id

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
			// handle applicant
			func() error {
				return l.updateApplicants(tx, oldApplicants, data.applicants, oldLoan)
			},
			// loan borrower
			func() error {
				var err error
				data.loanBorrower.LoanId = oldLoan.Id
				if oldLoan.LoanBorrower == nil {
					_, err = l.svcCtx.LoanBorrowerRepo.CreateWithTx(l.ctx, tx, data.loanBorrower, false)
				} else {
					_, err = l.svcCtx.LoanBorrowerRepo.UpdateOneWithTx(l.ctx, tx, oldLoan.LoanBorrower.Id, data.loanBorrower)
				}
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
				err = l.svcCtx.PropertyBaseRepo.SyncTotalAmountWithTx(l.ctx, tx, previousPropertyBaseId)
				if propertyBase.Id != previousPropertyBaseId {
					err = l.svcCtx.PropertyBaseRepo.SyncTotalAmountWithTx(l.ctx, tx, propertyBase.Id)
				}
				return err
			},

			// loan operation log
			func() error {
				operatorId := in.Data.LenderEmployeeId

				fromStatus := types.LoanOperationStatus{
					Stage:         oldLoan.Stage,
					PendingStatus: oldLoan.PendingStatus,
				}

				toStatus := types.LoanOperationStatus{
					Stage:         data.loan.Stage,
					PendingStatus: data.loan.PendingStatus,
				}

				loanBaseLogic := NewLoanBaseLogic(l.ctx, l.svcCtx)
				return loanBaseLogic.CreateLoanOperationLog(tx, oldLoan.Id, operatorId, fromStatus, toStatus, constants.LoanOperationActionUpdate, nil)
			},
		)

		innerErr = eg.Error()
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

func (l *UpdateLoanLogic) handlePropertyBase(tx *gorm.DB, propertyId int64, propertyBase *model.PropertyBase) (newPropertyBase *model.PropertyBase, err error) {
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

func (l *UpdateLoanLogic) updateApplicants(tx *gorm.DB, oldItems []*model.Applicant, curItems []*model.Applicant, loan *model.Loan) (err error) {
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
		if loan.LoanProperty != nil && loan.LoanProperty.Id > 0 {
			newApplicant.LoanPropertyId = loan.LoanProperty.Id
		}
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
		if loan.LoanProperty != nil && loan.LoanProperty.Id > 0 {
			updateApplicant.LoanPropertyId = loan.LoanProperty.Id
		}
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

func (l *UpdateLoanLogic) buildModelData(reqData *lender.LoanUpdateRequestData, oldLoan *model.Loan) (data *updateLoanModelData, err error) {
	if oldLoan == nil {
		return nil, errors.WithStack(errors.New("invalid loan"))
	}

	if reqData == nil {
		return nil, errors.WithStack(errors.New("invalid request data"))
	}
	if reqData.LoanProperty == nil {
		return nil, errors.WithStack(errors.New("invalid loan property"))
	}
	if reqData.LoanProperty.PropertyBase == nil {
		return nil, errors.WithStack(errors.New("invalid base property"))
	}
	// loan
	loan := oldLoan
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
	loan.BusinessStrategy = reqData.BusinessStrategy

	// loan_borrower
	loanBorrower := &model.LoanBorrower{
		MailingAddress: types.Address{},
	}
	err = copier.CopyWithOption(loanBorrower, reqData.LoanBorrower, copier.Option{
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

	// applicants
	applicants := make([]*model.Applicant, 0)
	err = copier.CopyWithOption(&applicants, &reqData.Applicants, copier.Option{
		IgnoreEmpty: true,
		DeepCopy:    true,
		Converters: []copier.TypeConverter{
			copierx.StrToSqlTimePtr(timex.LayoutUsStyle),
		},
	})
	if err != nil {
		return nil, errorx.GetFailedTransErr(err)
	}
	slice.ForEach(applicants, func(index int, item *model.Applicant) {
		item.BorrowerId = oldLoan.BorrowerId
		item.FullName = stringx.FormatFullName(item.FirstName, item.MidName, item.LastName)
	})

	// loan_property
	loanProperty := &model.LoanProperty{
		PropertyBase: &model.PropertyBase{},
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
	loanProperty.PropertyBase.BorrowerId = oldLoan.BorrowerId
	loanProperty.Unit = loanProperty.PropertyBase.Unit

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

	// apply form file
	// var applyFormFile *model.File
	// if len(reqData.ApplyFormFile) > 0 {
	// 	applyFormFile = &model.File{
	// 		File:     reqData.ApplyFormFile,
	// 		FileType: constants.FileTypeApplyForm,
	// 		BizType:  constants.FileBizTypeLoan,
	// 	}
	// }

	data = &updateLoanModelData{
		loan:                   loan,
		loanProperty:           loanProperty,
		loanBorrower:           loanBorrower,
		applicants:             applicants,
		transactionExperiences: transactionExperiences,
		realEstatesExperiences: realEstatesExperiences,
		loanDemographic:        loanDemographic,
		loanThirdServices:      loanThirdServices,
		// applyFormFile:          applyFormFile,
	}
	return
}
