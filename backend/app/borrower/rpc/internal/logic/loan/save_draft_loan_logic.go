package loanlogic

import (
	"context"

	"github.com/duke-git/lancet/v2/condition"
	"github.com/duke-git/lancet/v2/mathutil"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/timex"
	"zeebee-backend/common/utils"
	"zeebee-backend/common/uuid"
	"zeebee-backend/domains/zeebee/constants"
	data2 "zeebee-backend/domains/zeebee/data"
	"zeebee-backend/domains/zeebee/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveDraftLoanLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewSaveDraftLoanLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveDraftLoanLogic {
	return &SaveDraftLoanLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *SaveDraftLoanLogic) SaveDraftLoan(in *borrower.LoanApplyRequest) (resp *borrower.LoanApplyResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	var loan *model.Loan
	var oldApplicants []*model.LoanApplicant
	var loanBorrower *model.LoanBorrower
	var loanProperty *model.LoanProperty
	var oldBorrowerPropertyId int64
	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	// get original loan
	if in.Data.Id > 0 {
		loan, err = l.svcCtx.LoanRepo.GetBorrowerLoan(l.ctx, in.Data.Id, claim.OrganizationId)
		if err != nil {
			return nil, err
		}
		if loan.Source != constants.LoanSourceSystem && loan.Source != constants.LoanSourcePreapproval {
			return nil, errorx.WithStack(errorx.New("Invalid loan source"))
		}

		stages := []constants.LoanStage{
			constants.LoanStageDraft,
		}
		if !slice.Contain(stages, loan.Stage) && loan.PendingStatus != constants.LoanPendingStatusNone {
			return nil, errorx.New("Loan's stage should be draft")
		}
		oldApplicants = loan.Applicants
		if loan.LoanProperty != nil {
			oldBorrowerPropertyId = loan.LoanProperty.BorrowerPropertyId
		}
	}
	// build data
	updateLogic := NewUpdateApplyLoanLogic(l.ctx, l.svcCtx)
	data, err := l.buildModelData(in.Data, loan, claim)
	if err != nil {
		return nil, err
	}
	if loan != nil {
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
						in.Data.LenderCompanyId = loan.LenderOrganizationId
					}
					_, err = l.svcCtx.LoanRepo.UpdateOneWithTx(l.ctx, tx, in.Data.Id, data.loan)
				} else {
					var uid string
					uid, err = uuid.GetSonyflakeIDString()
					if err != nil {
						return err
					}
					data.loan.Code = uid
					savaLoan, err := l.svcCtx.LoanRepo.CreateWithTx(l.ctx, tx, data.loan, false)
					if err != nil {
						return err
					}
					loan, err = l.svcCtx.LoanRepo.GetBorrowerLoanWithTx(l.ctx, tx, savaLoan.Id, claim.OrganizationId)
				}
				return err
			},
			// handle property
			func() error {
				var err error
				if data.loanProperty == nil {
					return nil
				}
				data.loanProperty.LoanId = loan.Id
				loanProperty, err = l.updateProperty(tx, data.loanProperty, loan)
				return err
			},
			// loan borrower
			func() error {
				var err error
				loanBorrower, err = l.updateLoanBorrower(tx, data.loanBorrower, loan)
				return err
			},
			// handle applicant
			func() error {
				slice.ForEach(data.applicants, func(index int, item *model.LoanApplicant) {
					if loanBorrower != nil {
						item.LoanBorrowerId = loanBorrower.Id
					}
					if loanProperty != nil {
						item.LoanPropertyId = loanProperty.Id
					}
				})
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
			// sync BorrowerProperty loan count
			func() error {
				var err error
				if in.Data.LoanProperty == nil {
					return nil
				}
				err = l.svcCtx.BorrowerPropertyRepo.SyncTotalAmountAndLoanCountWithTx(l.ctx, tx, in.Data.LoanProperty.BorrowerPropertyId)

				if err != nil {
					return err
				}

				if oldBorrowerPropertyId != in.Data.LoanProperty.BorrowerPropertyId {
					err = l.svcCtx.BorrowerPropertyRepo.SyncTotalAmountAndLoanCountWithTx(l.ctx, tx, oldBorrowerPropertyId)
				}

				return err
			},
		)

		innerErr = eg.Wait()
		return innerErr
	})

	return &borrower.LoanApplyResponse{
		Data: &borrower.LoanApplyResponseData{
			Id: loan.Id,
		},
	}, nil
}

// updateProperty update loan property and property base
func (l *SaveDraftLoanLogic) updateProperty(tx *gorm.DB, property *model.LoanProperty, loan *model.Loan) (newProperty *model.LoanProperty, err error) {
	if property == nil {
		return nil, nil
	}
	if loan.LoanProperty != nil && loan.LoanProperty.Id > 0 {
		newProperty, err = l.svcCtx.LoanPropertyRepo.UpdateOneWithTx(l.ctx, tx, loan.LoanProperty.Id, property)
	} else {
		newProperty, err = l.svcCtx.LoanPropertyRepo.CreateWithTx(l.ctx, tx, property, false)
	}
	// _, err = l.svcCtx.LoanPropertyRepo.UpdateOneWithTx(l.ctx, tx, property.Id, property)
	return newProperty, err
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
func (l *SaveDraftLoanLogic) buildModelData(reqData *borrower.LoanApplyRequestData, oldLoan *model.Loan, claim utils.JwtUser) (data *applyLoanModelData, err error) {

	if reqData == nil {
		return nil, errorx.New("Invalid request data")
	}

	if reqData.LenderCompanyId == 0 {
		return nil, errorx.New("Invalid lender id")
	}

	// loan
	loan := &model.Loan{}
	if oldLoan != nil && oldLoan.Id > 0 {
		_ = copier.CopyWithOption(loan, oldLoan, copier.Option{
			CaseSensitive: true,
			DeepCopy:      true,
			Converters: []copier.TypeConverter{
				copierx.StrToSqlTimePtr(timex.LayoutUsStyle),
			},
		})
		loan.RemoveRelationships()
	}
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
	if loan.Id == 0 {
		loan.CreatorId = claim.OrganizationUserId
		loan.CreatorName = claim.Username
		loan.BorrowerOrganizationId = claim.OrganizationId
	}
	loan.Stage = constants.LoanStageDraft
	loan.PendingStatus = constants.LoanPendingStatusNone
	// loan.Source = constants.LoanSourceSystem
	loan.Source = condition.Ternary(reqData.IsPreapproveLoan, constants.LoanSourcePreapproval, constants.LoanSourceSystem)
	if reqData.IsPreapproveLoan {
		preapprovalInfo, _ := copierx.Trans[borrower.LoanPreapproveInfo, data2.LoanPreapprovalInfo](reqData.LoanPreapproveInfo)
		loan.Extend.LoanPreapprovalInfo = preapprovalInfo
	}

	if reqData.LenderCompanyId > 0 {
		lender, err := l.svcCtx.OrganizationRepo.FindOneById(l.ctx, reqData.LenderCompanyId)
		if err != nil {
			return nil, err
		}
		loan.LenderOrganizationId = lender.Id
		loan.LenderCompany = lender.Name
	}

	// LoanBorrower
	var loanBorrower *model.LoanBorrower
	if reqData.LoanBorrower != nil && reqData.LoanBorrower.BorrowerEntityId > 0 {
		loanBorrower, err = l.svcCtx.BorrowerEntityRepo.TranToLoanBorrowerWithId(l.ctx, reqData.LoanBorrower.BorrowerEntityId, claim.OrganizationId)
		if err != nil {
			return nil, err
		}
		loanBorrower.LoanId = reqData.Id
	}

	// LoanProperty
	var loanProperty *model.LoanProperty
	if reqData.LoanProperty != nil {
		if reqData.LoanProperty.BorrowerPropertyId > 0 {
			loanProperty, err = l.svcCtx.BorrowerPropertyRepo.TranToLoanPropertyWithId(l.ctx, reqData.LoanProperty.BorrowerPropertyId, claim.OrganizationId)
			if err != nil {
				return nil, err
			}
		} else {
			loanProperty = &model.LoanProperty{}
		}
		err = copier.CopyWithOption(loanProperty, reqData.LoanProperty, copier.Option{
			CaseSensitive: false,
			DeepCopy:      true,
			Converters: []copier.TypeConverter{
				copierx.StrToSqlTimePtr(timex.LayoutUsStyle),
			},
		})
		if err != nil {
			return nil, errorx.GetFailedTransErr(err)
		}
		if loanProperty.YearlyLiabilityExpenseAmount > 0 {
			loanProperty.Dscr = mathutil.RoundToFloat((loanProperty.EstYearlyRentalIncomeAmount/loanProperty.YearlyLiabilityExpenseAmount)*100, 2)
		}
	}

	// LoanApplicants
	var applicants []*model.LoanApplicant
	var totalOwnPercent int32
	for _, applicant := range reqData.Applicants {
		totalOwnPercent += applicant.OwnPercent
		loanApplicant, err := l.svcCtx.BorrowerApplicantRepo.TranToLoanApplicantWithId(l.ctx, applicant.BorrowerApplicantId, claim.OrganizationId)
		if err != nil {
			return nil, err
		}
		loanApplicant.Id = applicant.Id
		loanApplicant.LoanId = reqData.Id
		loanApplicant.BorrowerOrganizationId = claim.OrganizationId
		loanApplicant.IsPrimary = applicant.IsPrimary
		loanApplicant.OwnPercent = applicant.OwnPercent
		applicants = append(applicants, loanApplicant)
	}
	if totalOwnPercent > 100 {
		return nil, errorx.New("Total own percent is greater than 100")
	}

	// transaction experience
	var transactionExperiences []*model.LoanExperience
	if len(reqData.TransactionExperiences) > 0 {
		transactionExperiences = make([]*model.LoanExperience, 0)
		err = copier.CopyWithOption(&transactionExperiences, &reqData.TransactionExperiences, copier.Option{
			IgnoreEmpty: false,
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
			CaseSensitive: false,
			DeepCopy:      true,
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
			CaseSensitive: false,
			DeepCopy:      true,
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
			CaseSensitive: false,
			DeepCopy:      true,
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
