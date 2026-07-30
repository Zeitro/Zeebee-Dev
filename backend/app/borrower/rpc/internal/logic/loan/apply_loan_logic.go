package loanlogic

import (
	"context"

	"github.com/duke-git/lancet/v2/mathutil"
	"github.com/duke-git/lancet/v2/slice"
	"gorm.io/gorm"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/timex"
	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/proto"
)

type ApplyLoanLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

type applyLoanModelData struct {
	loan                   *model.Loan
	loanBorrower           *model.LoanBorrower
	loanProperty           *model.LoanProperty
	loanDemographic        *model.LoanDemographic
	loanThirdServices      []*model.LoanThirdService
	applicants             []*model.LoanApplicant
	transactionExperiences []*model.LoanExperience
	realEstatesExperiences []*model.LoanExperience
	// applyFormFile          *model.File
}

func NewApplyLoanLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApplyLoanLogic {
	return &ApplyLoanLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *ApplyLoanLogic) ApplyLoan(in *borrower.LoanApplyRequest) (resp *borrower.LoanApplyResponse, err error) {

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, errorx.Wrap(err, "fail to get user claim")
	}

	data, err := l.buildModelData(in.Data, nil, claim)
	if err != nil {
		return nil, err
	}
	// apply loan
	// var propertyBase *model.BorrowerProperty
	var loan *model.Loan
	var loanBorrower *model.LoanBorrower
	var loanProperty *model.LoanProperty
	err = l.svcCtx.LoanRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
		var innerErr error

		// check loan
		if len(data.loan.Code) == 0 {
			innerErr = errorx.WithStack(errorx.New("code is empty"))
		}
		isExisted, innerErr := l.svcCtx.LoanRepo.IsExistedByColumnsWithTx(l.ctx, tx, map[string]interface{}{"code": data.loan.Code})
		if innerErr != nil {
			return innerErr
		}
		if isExisted {
			return errorx.WithStack(errorx.New("code already exists"))
		}

		// save loan
		loan, innerErr = l.svcCtx.LoanRepo.CreateWithTx(l.ctx, tx, data.loan, false)
		if innerErr != nil {
			return innerErr
		}

		// save loan borrower
		data.loanBorrower.Id = 0
		data.loanBorrower.LoanId = loan.Id
		data.loanBorrower.Applicants = nil
		loanBorrower, innerErr = l.svcCtx.LoanBorrowerRepo.CreateWithTx(l.ctx, tx, data.loanBorrower, false)
		if innerErr != nil {
			return innerErr
		}

		// save loan property
		data.loanProperty.Id = 0
		data.loanProperty.LoanId = loan.Id
		loanProperty, innerErr = l.svcCtx.LoanPropertyRepo.CreateWithTx(l.ctx, tx, data.loanProperty, false)
		if innerErr != nil {
			return innerErr
		}

		// save applicants
		slice.ForEach(data.applicants, func(index int, item *model.LoanApplicant) {
			item.Id = 0
			item.BorrowerOrganizationId = claim.OrganizationId
			item.LoanId = loan.Id
			item.LoanBorrowerId = loanBorrower.Id
			item.LoanPropertyId = loanProperty.Id
		})
		_, innerErr = l.svcCtx.LoanApplicantRepo.BatchCreateWithTx(l.ctx, tx, data.applicants, false)
		if innerErr != nil {
			return innerErr
		}

		// save demographic
		data.loanDemographic.LoanId = loan.Id
		_, innerErr = l.svcCtx.LoanDemographicRepo.CreateWithTx(l.ctx, tx, data.loanDemographic, false)
		if innerErr != nil {
			return innerErr
		}

		// save experience service
		if len(data.transactionExperiences) > 0 {
			slice.ForEach(data.transactionExperiences, func(index int, item *model.LoanExperience) {
				item.LoanId = loan.Id
			})
			_, innerErr = l.svcCtx.LoanExperienceRepo.BatchCreateWithTx(l.ctx, tx, data.transactionExperiences, false)
			if innerErr != nil {
				return innerErr
			}
		}

		if len(data.realEstatesExperiences) > 0 {
			slice.ForEach(data.realEstatesExperiences, func(index int, item *model.LoanExperience) {
				item.LoanId = loan.Id
			})
			_, innerErr = l.svcCtx.LoanExperienceRepo.BatchCreateWithTx(l.ctx, tx, data.realEstatesExperiences, false)
			if innerErr != nil {
				return innerErr
			}
		}

		// third services
		if len(data.loanThirdServices) > 0 {
			slice.ForEach(data.loanThirdServices, func(index int, item *model.LoanThirdService) {
				item.LoanId = loan.Id
			})
			_, innerErr = l.svcCtx.LoanThirdServiceRepo.BatchCreateWithTx(l.ctx, tx, data.loanThirdServices, false)
			if innerErr != nil {
				return innerErr
			}
		}

		// apply form
		if len(in.Data.ApplyFormFile) > 0 {
			innerErr = l.svcCtx.LoanDocumentRepo.SaveApplyFormDocumentWithTx(l.ctx, tx, loan.Id, in.Data.ApplyFormFile)
			if innerErr != nil {
				return innerErr
			}
		}

		// Create loan documents with pending upload status
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
			_, innerErr = l.svcCtx.LoanDocumentRepo.BatchCreateWithTx(l.ctx, tx, loanDocuments)
			if innerErr != nil {
				return innerErr
			}
		}

		// sync total amount
		innerErr = l.svcCtx.BorrowerPropertyRepo.SyncTotalAmountAndLoanCountWithTx(l.ctx, tx, data.loanProperty.BorrowerPropertyId)

		return innerErr
	})
	if err != nil {
		return nil, errorx.Wrap(err, "fail to apply loan")
	}

	return &borrower.LoanApplyResponse{
		Data: &borrower.LoanApplyResponseData{
			Id: loan.Id,
		},
	}, nil
}

// buildModelData return save model data
func (l *ApplyLoanLogic) buildModelData(reqData *borrower.LoanApplyRequestData, oldLoan *model.Loan, claim utils.JwtUser) (data *applyLoanModelData, err error) {
	var borrowerProperty *model.BorrowerProperty

	if reqData == nil {
		return nil, errorx.New("invalid request data")
	}

	// When BorrowerPropertyTransactionId is provided, derive property from transaction
	if reqData.BorrowerPropertyTransactionId > 0 {
		transaction, txErr := l.svcCtx.BorrowerPropertyTransactionRepo.GetOrganizationTransaction(l.ctx, reqData.BorrowerPropertyTransactionId, claim.OrganizationId)
		if txErr != nil {
			return nil, errorx.Wrap(txErr, "invalid transaction")
		}
		// Override loan property to use transaction's property
		if reqData.LoanProperty == nil {
			reqData.LoanProperty = &borrower.LoanPropertyApplyEntity{}
		}
		reqData.LoanProperty.BorrowerPropertyId = transaction.BorrowerPropertyId
	}

	if reqData.LoanProperty == nil {
		return nil, errorx.New("invalid loan property")
	}

	if reqData.LoanProperty.BorrowerPropertyId == 0 {
		return nil, errorx.New("invalid borrower property")
	}

	// loan
	loan := &model.Loan{}
	if oldLoan != nil {
		loan, err = copierx.Trans[model.Loan, model.Loan](oldLoan)

		if err != nil {
			return nil, err
		}

		loan.RemoveRelationships()
	}
	err = copierx.Assign[borrower.LoanApplyRequestData, model.Loan](reqData, loan, copierx.StrToSqlTimePtr(timex.LayoutUsStyle))

	if err != nil {
		return nil, err
	}
	loan.CreatorId = claim.OrganizationUserId
	loan.CreatorName = claim.Username
	loan.BorrowerOrganizationId = claim.OrganizationId
	loan.Stage = constants.LoanStagePending
	loan.PendingStatus = constants.LoanPendingStatusQuoting
	loan.Source = constants.LoanSourceSystem
	// get lender
	lender, err := l.svcCtx.OrganizationRepo.FindOneById(l.ctx, reqData.LenderCompanyId)
	if err != nil {
		return nil, err
	}
	loan.LenderOrganizationId = lender.Id
	loan.LenderCompany = lender.Name

	// LoanBorrower
	loanBorrowerData := reqData.LoanBorrower
	if loanBorrowerData == nil {
		return nil, errorx.WithStack(errorx.New("invalid loan borrower"))
	}
	loanBorrower, err := l.svcCtx.BorrowerEntityRepo.TranToLoanBorrowerWithId(l.ctx, loanBorrowerData.BorrowerEntityId, claim.OrganizationId)
	if err != nil {
		return nil, err
	}
	if oldLoan != nil && oldLoan.LoanBorrower != nil {
		loanBorrower.Id = oldLoan.LoanBorrower.Id
	}
	loanBorrower.BorrowerEntityId = loanBorrowerData.BorrowerEntityId
	// applicants
	applicants := make([]*model.LoanApplicant, 0)
	var totalOwnPercent int32 = 0
	for _, applicant := range reqData.Applicants {
		totalOwnPercent += applicant.OwnPercent
		borrowerApplicant, err := l.svcCtx.BorrowerApplicantRepo.GetBorrowerApplicantById(l.ctx, applicant.BorrowerApplicantId, claim.OrganizationId)
		if err != nil {
			return nil, err
		}
		applicantItem := &model.LoanApplicant{}
		applicantItem, err = copierx.Trans[model.BorrowerApplicant, model.LoanApplicant](borrowerApplicant, copierx.StrToSqlTimePtr(timex.LayoutUsStyle))
		if err != nil {
			return nil, err
		}
		applicantItem.Id = applicant.Id
		applicantItem.IsPrimary = applicant.IsPrimary
		applicantItem.BorrowerApplicantId = applicant.BorrowerApplicantId
		applicantItem.OwnPercent = applicant.OwnPercent
		applicantItem.InvitationStatus = constants.ApplicantInvitationStatusAccepted
		applicants = append(applicants, applicantItem)
	}
	if totalOwnPercent > 100 {
		return nil, errorx.New("Total own percent is greater than 100")
	}

	slice.SortBy(applicants, func(a, b *model.LoanApplicant) bool {
		return b.IsPrimary
	})
	primaryCount := slice.CountBy(applicants, func(index int, item *model.LoanApplicant) bool {
		return item.IsPrimary
	})
	if primaryCount == 0 {
		return nil, errorx.New("No primary applicant")
	}
	if primaryCount > 1 {
		return nil, errorx.New("Multiple primary applicants")
	}

	// loan_property
	borrowerProperty, err = l.svcCtx.BorrowerPropertyRepo.GetOneByIdForOrganization(l.ctx, reqData.LoanProperty.BorrowerPropertyId, claim.OrganizationId)
	if err != nil {
		return nil, err
	}
	loanProperty := &model.LoanProperty{}
	loanProperty, err = copierx.Trans[borrower.LoanPropertyApplyEntity, model.LoanProperty](reqData.LoanProperty, copierx.StrToSqlTimePtr(timex.LayoutUsStyle))
	if err != nil {
		return nil, err
	}
	// loanProperty.Id = 0
	loanProperty.PropertyType = borrowerProperty.PropertyType
	loanProperty.Unit = borrowerProperty.Unit
	loanProperty.Bedrooms = borrowerProperty.Bedrooms
	loanProperty.Bathrooms = borrowerProperty.Bathrooms
	loanProperty.LivingAreaSqft = borrowerProperty.LivingAreaSqft
	loanProperty.Address = borrowerProperty.Address
	loanProperty.BorrowerPropertyId = reqData.LoanProperty.BorrowerPropertyId
	if oldLoan != nil && oldLoan.LoanProperty != nil {
		loanProperty.Id = oldLoan.LoanProperty.Id
	}

	if loanProperty.YearlyLiabilityExpenseAmount > 0 {
		loanProperty.Dscr = mathutil.RoundToFloat((loanProperty.EstYearlyRentalIncomeAmount/loanProperty.YearlyLiabilityExpenseAmount)*100, 2)
	}

	// transaction experience
	transactionExperiences := make([]*model.LoanExperience, 0)
	err = copierx.Assign[[]*proto.LoanExperienceEntity, []*model.LoanExperience](&reqData.TransactionExperiences, &transactionExperiences, copierx.StrToSqlTimePtr(timex.LayoutUsStyle))
	if err != nil {
		return nil, err
	}

	// real_estates experience
	realEstatesExperiences := make([]*model.LoanExperience, 0)
	err = copierx.Assign[[]*proto.LoanExperienceEntity, []*model.LoanExperience](&reqData.RealEstateExperiences, &realEstatesExperiences, copierx.StrToSqlTimePtr(timex.LayoutUsStyle))
	if err != nil {
		return nil, err
	}

	// third services
	loanThirdServices := make([]*model.LoanThirdService, 0)
	err = copierx.Assign[[]*proto.LoanThirdServiceEntity, []*model.LoanThirdService](&reqData.LoanThirdServices, &loanThirdServices, copierx.StrToSqlTimePtr(timex.LayoutUsStyle))
	if err != nil {
		return nil, err
	}

	// demographic
	loanDemographic := &model.LoanDemographic{}
	loanDemographic, err = copierx.Trans[proto.LoanDemographicEntity, model.LoanDemographic](reqData.LoanDemographic, copierx.StrToSqlTimePtr(timex.LayoutUsStyle))
	if err != nil {
		return nil, err
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
		// applyFormFile:          applyFormFile,
	}
	return
}
