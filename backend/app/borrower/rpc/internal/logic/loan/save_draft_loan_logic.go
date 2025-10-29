package loanlogic

import (
	"context"

	"github.com/duke-git/lancet/v2/condition"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"zeebee-backend/app/borrower/rpc/internal/repo/model"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/constants"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/stringx"
	"zeebee-backend/common/timex"
	"zeebee-backend/common/types"
	"zeebee-backend/common/uuid"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type SaveDraftLoanLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSaveDraftLoanLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveDraftLoanLogic {
	return &SaveDraftLoanLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SaveDraftLoanLogic) SaveDraftLoan(in *borrower.LoanApplyRequest) (resp *borrower.LoanApplyResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	var loan *model.Loan
	// var propertyBase *model.PropertyBase
	var oldApplicants []*model.Applicant
	var loanBorrower *model.LoanBorrower
	updateLogic := NewUpdateApplyLoanLogic(l.ctx, l.svcCtx)
	data, err := l.buildModelData(in.Data)
	if err != nil {
		return nil, err
	}
	if data.loan.Id > 0 {
		loan, err = l.svcCtx.LoanRepo.GetBorrowerLoan(l.ctx, data.loan.Id, in.Data.BorrowerId)
		if err != nil {
			return nil, err
		}
		if loan.Source != constants.LoanSourceSystem && loan.Source != constants.LoanSourcePreApproval {
			return nil, errors.WithStack(errors.New("invalid loan source"))
		}

		stages := []constants.LoanStage{
			constants.LoanStageDraft,
			// constants.LoanStagePredraft,
		}
		if !slice.Contain(stages, loan.Stage) && loan.PendingStatus != constants.LoanPendingStatusNone {
			return nil, errors.New("loan's stage should be draft")
		}
		oldApplicants = loan.Applicants
		data.loan.Stage = loan.Stage
	}
	// update loan
	err = l.svcCtx.LoanPropertyRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
		var innerErr error
		eg := errorx.NewErrorGroup()
		eg.Do(
			// update loan
			func() error {
				var err error
				data.loan.RemoveRelationships()

				if in.Data.Id > 0 {
					if in.Data.LenderCompanyId == 0 {
						in.Data.LenderCompanyId = loan.LenderCompanyId
					}
					_, err = l.svcCtx.LoanRepo.UpdateOneWithTx(l.ctx, tx, in.Data.Id, data.loan, true)
				} else {
					var uid string
					uid, err = uuid.GetSonyflakeIDString()
					if innerErr != nil {
						return innerErr
					}
					data.loan.Code = uid
					savaLoan, err := l.svcCtx.LoanRepo.CreateWithTx(l.ctx, tx, data.loan, false)
					if err != nil {
						return err
					}
					loan, err = l.svcCtx.LoanRepo.GetBorrowerLoanWithTx(l.ctx, tx, savaLoan.Id, in.Data.BorrowerId)
				}
				return err
			},
			// handle property
			func() error {
				if data.loanProperty == nil {
					return nil
				}
				return l.updateProperty(tx, data.loanProperty, loan)
			},
			// loan borrower
			func() error {
				var err error
				loanBorrower, err = l.updateLoanBorrower(tx, data.loanBorrower, loan)
				return err
			},
			// handle applicant
			func() error {
				if loanBorrower != nil {
					slice.ForEach(data.applicants, func(index int, item *model.Applicant) {
						item.LoanBorrowerId = loanBorrower.Id
					})
				}

				return updateLogic.updateApplicants(tx, oldApplicants, data.applicants, loan)
			},
			// loan demographic
			func() error {
				return l.updateLoanDemographic(tx, data.loanDemographic, loan)
			},
			// loan experience service
			func() error {
				oldExperiences := loan.LoanExperiences
				transactionExperiences := data.transactionExperiences
				realEstatesExperiences := data.realEstatesExperiences
				allExperiences := append(transactionExperiences, realEstatesExperiences...)
				return updateLogic.updateLoanExperiences(tx, oldExperiences, allExperiences, loan)
			},
			// loan third service
			func() error {
				return updateLogic.updateThirdServices(tx, loan.LoanThirdServices, data.loanThirdServices, loan)
			},
			// apply form
			func() error {
				var err error
				if len(in.Data.ApplyFormFile) > 0 {
					err = l.svcCtx.LoanDocumentRepo.SaveApplyFormDocumentWithTx(l.ctx, tx, loan.Id, in.Data.ApplyFormFile)
				}
				return err
			},
		)

		innerErr = eg.Error()
		return innerErr
	})

	return &borrower.LoanApplyResponse{
		Data: &borrower.LoanApplyResponseData{
			Id: loan.Id,
		},
	}, nil
}

// updateProperty update loan property and property base
func (l *SaveDraftLoanLogic) updateProperty(tx *gorm.DB, property *model.LoanProperty, loan *model.Loan) (err error) {
	if property == nil {
		return nil
	}
	propertyBase := property.PropertyBase
	if propertyBase == nil {
		return errors.WithStack(errors.New("property base is required"))
	}
	propertyBase, err = l.svcCtx.PropertyBaseRepo.GetOrCreateWithTx(l.ctx, tx, propertyBase)
	if err != nil {
		return err
	}

	// save loan property
	property.LoanId = loan.Id
	property.PropertyBaseId = propertyBase.Id
	property, err = l.svcCtx.LoanPropertyRepo.CreateWithTx(l.ctx, tx, property, false)
	return err
}

// updateLoanBorrower update loan borrower
func (l *SaveDraftLoanLogic) updateLoanBorrower(tx *gorm.DB, loanBorrower *model.LoanBorrower, loan *model.Loan) (newLoanBorrower *model.LoanBorrower, err error) {

	if loanBorrower == nil {
		return nil, nil
	}
	if loan.LoanBorrower != nil && loan.LoanBorrower.Id > 0 {
		newLoanBorrower, err = l.svcCtx.LoanBorrowerRepo.UpdateOneWithTx(l.ctx, tx, loan.LoanBorrower.Id, loanBorrower)
	} else {
		loanBorrower.LoanId = loan.Id
		loanBorrower.Applicants = nil
		newLoanBorrower, err = l.svcCtx.LoanBorrowerRepo.CreateWithTx(l.ctx, tx, loanBorrower, false)

	}
	return newLoanBorrower, err
}

// updateLoanDemographic update loan demographic
func (l *SaveDraftLoanLogic) updateLoanDemographic(tx *gorm.DB, loanDemographic *model.LoanDemographic, loan *model.Loan) (err error) {

	if loanDemographic == nil {
		return nil
	}
	if loan.LoanDemographic != nil && loan.LoanDemographic.Id > 0 {
		_, err = l.svcCtx.LoanDemographicRepo.UpdateOneWithTx(l.ctx, tx, loan.LoanDemographic.Id, loanDemographic)
	} else {
		loanDemographic.LoanId = loan.Id
		_, err = l.svcCtx.LoanDemographicRepo.CreateWithTx(l.ctx, tx, loanDemographic, false)
	}
	return err
}

// buildModelData return save model data
func (l *SaveDraftLoanLogic) buildModelData(reqData *borrower.LoanApplyRequestData) (data *applyLoanModelData, err error) {

	if reqData == nil {
		return nil, errors.WithStack(errors.New("invalid request data"))
	}

	if reqData.LenderCompanyId == 0 {
		return nil, errors.WithStack(errors.New("invalid lender id"))
	}

	// loan
	loan := &model.Loan{}
	err = copier.CopyWithOption(loan, reqData, copier.Option{
		IgnoreEmpty: true,
		DeepCopy:    true,
		Converters: []copier.TypeConverter{
			copierx.StrToSqlTimePtr(timex.LayoutUsStyle),
		},
	})
	if err != nil {
		return nil, errorx.GetFailedTransErr(err)
	}
	loan.BorrowerId = reqData.BorrowerId
	loan.Stage = constants.LoanStageDraft
	loan.PendingStatus = constants.LoanPendingStatusNone
	// loan.Source = constants.LoanSourceSystem
	loan.Source = condition.Ternary(reqData.IsPreapproveLoan, constants.LoanSourcePreApproval, constants.LoanSourceSystem)
	if reqData.LenderCompanyId > 0 {
		lender, err := l.svcCtx.LenderCompanyRepo.FindOneById(l.ctx, reqData.LenderCompanyId)
		if err != nil {
			return nil, err
		}
		loan.LenderCompanyId = lender.Id
		loan.LenderCompany = lender.Name
	}

	// loan_borrower
	var loanBorrower *model.LoanBorrower
	if reqData.LoanBorrower != nil {
		loanBorrower = &model.LoanBorrower{
			MailingAddress: types.Address{},
		}
		err = copier.CopyWithOption(loanBorrower, reqData.LoanBorrower, copier.Option{
			IgnoreEmpty:   true,
			DeepCopy:      true,
			CaseSensitive: true,
			Converters: []copier.TypeConverter{
				copierx.StrToSqlTimePtr(timex.LayoutUsStyle),
			},
		})
		if err != nil {
			return nil, errorx.GetFailedTransErr(err)
		}
	}

	// applicants
	var applicants []*model.Applicant
	if len(reqData.Applicants) > 0 {
		applicants = make([]*model.Applicant, 0)
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
	}

	slice.ForEach(applicants, func(index int, item *model.Applicant) {
		item.BorrowerId = reqData.BorrowerId
		item.FullName = stringx.FormatFullName(item.FirstName, item.MidName, item.LastName)
	})

	// loan_property
	var loanProperty *model.LoanProperty
	if reqData.LoanProperty != nil {
		if reqData.LoanProperty.PropertyBase == nil || reqData.LoanProperty.PropertyBase.Address == nil {
			return nil, errors.WithStack(errors.New("property base is required"))
		}
		addr := reqData.LoanProperty.PropertyBase.Address
		if len(addr.Address1) == 0 || len(addr.City) == 0 || len(addr.State) == 0 || len(addr.ZipCode) == 0 {
			return nil, errors.WithStack(errors.New("property address is required"))
		}

		loanProperty = &model.LoanProperty{
			PropertyBase: &model.PropertyBase{},
		}
		err = copier.CopyWithOption(loanProperty, reqData.LoanProperty, copier.Option{
			IgnoreEmpty: true,
			DeepCopy:    true,
			Converters: []copier.TypeConverter{
				copierx.StrToSqlTimePtr(timex.LayoutUsStyle),
			},
		})
		if err != nil {
			return nil, errorx.GetFailedTransErr(err)
		}
		loanProperty.PropertyBase.BorrowerId = reqData.BorrowerId
		loanProperty.Unit = loanProperty.PropertyBase.Unit
	}

	// transaction experience
	var transactionExperiences []*model.LoanExperience
	if len(reqData.TransactionExperiences) > 0 {
		transactionExperiences = make([]*model.LoanExperience, 0)
		err = copier.CopyWithOption(&transactionExperiences, &reqData.TransactionExperiences, copier.Option{
			IgnoreEmpty: true,
			DeepCopy:    true,
			Converters: []copier.TypeConverter{
				copierx.StrToSqlTimePtr(timex.LayoutUsStyle),
			},
		})
		if err != nil {
			return nil, errorx.GetFailedTransErr(err)
		}
	}

	// real_estates experience
	var realEstatesExperiences []*model.LoanExperience
	if len(reqData.RealEstateExperiences) > 0 {
		realEstatesExperiences = make([]*model.LoanExperience, 0)
		err = copier.CopyWithOption(&realEstatesExperiences, &reqData.RealEstateExperiences, copier.Option{
			IgnoreEmpty: true,
			DeepCopy:    true,
			Converters: []copier.TypeConverter{
				copierx.StrToSqlTimePtr(timex.LayoutUsStyle),
			},
		})
		if err != nil {
			return nil, errorx.GetFailedTransErr(err)
		}
	}

	// third services
	var loanThirdServices []*model.LoanThirdService
	if len(reqData.LoanThirdServices) > 0 {
		loanThirdServices = make([]*model.LoanThirdService, 0)
		err = copier.CopyWithOption(&loanThirdServices, &reqData.LoanThirdServices, copier.Option{
			IgnoreEmpty: true,
			DeepCopy:    true,
			Converters: []copier.TypeConverter{
				copierx.StrToSqlTimePtr(timex.LayoutUsStyle),
			},
		})
		if err != nil {
			return nil, errorx.GetFailedTransErr(err)
		}
	}

	// demographic
	var loanDemographic *model.LoanDemographic
	if reqData.LoanDemographic != nil {
		loanDemographic = &model.LoanDemographic{}
		err = copier.CopyWithOption(&loanDemographic, &reqData.LoanDemographic, copier.Option{
			IgnoreEmpty: true,
			DeepCopy:    true,
			Converters: []copier.TypeConverter{
				copierx.StrToSqlTimePtr(timex.LayoutUsStyle),
			},
		})
		if err != nil {
			return nil, errorx.GetFailedTransErr(err)
		}
	}

	data = &applyLoanModelData{
		loan:                   loan,
		loanProperty:           loanProperty,
		loanBorrower:           loanBorrower,
		applicants:             applicants,
		transactionExperiences: transactionExperiences,
		realEstatesExperiences: realEstatesExperiences,
		loanDemographic:        loanDemographic,
		loanThirdServices:      loanThirdServices,
	}
	return
}
