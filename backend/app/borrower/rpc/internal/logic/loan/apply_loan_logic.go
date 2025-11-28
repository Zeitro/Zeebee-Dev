package loanlogic

import (
	"context"

	"github.com/duke-git/lancet/v2/mathutil"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/constants"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/timex"
	"zeebee-backend/domains/zeebee/model"
)

type ApplyLoanLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
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
		Logger: logx.WithContext(ctx),
	}
}

func (l *ApplyLoanLogic) ApplyLoan(in *borrower.LoanApplyRequest) (resp *borrower.LoanApplyResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	data, err := l.buildModelData(in.Data, nil)
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
			innerErr = errorx.WithStack(errorx.New("Code Is Empty"))
		}
		isExisted, innerErr := l.svcCtx.LoanRepo.IsExistedByColumnsWithTx(l.ctx, tx, map[string]interface{}{"code": data.loan.Code})
		if innerErr != nil {
			return innerErr
		}
		if isExisted {
			return errorx.WithStack(errorx.New("Code Already Exists"))
		}
		// todo check code whether generate by system

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
			item.BorrowerOrganizationId = in.Data.OrganizationId
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
		return nil, errorx.Wrap(err, "Failed To Apply Loan")
	}

	return &borrower.LoanApplyResponse{
		Data: &borrower.LoanApplyResponseData{
			Id: loan.Id,
		},
	}, nil
}

// buildModelData return save model data
func (l *ApplyLoanLogic) buildModelData(reqData *borrower.LoanApplyRequestData, oldLoan *model.Loan) (data *applyLoanModelData, err error) {
	var borrowerProperty *model.BorrowerProperty

	if reqData == nil {
		return nil, errorx.New("Invalid Request Data")
	}
	if reqData.LoanProperty == nil {
		return nil, errorx.New("Invalid Loan Property")
	}

	if reqData.LoanProperty.BorrowerPropertyId == 0 {
		return nil, errorx.New("Invalid Borrower Property")
	}

	// loan
	loan := &model.Loan{}
	if oldLoan != nil {
		err = copier.CopyWithOption(loan, oldLoan, copier.Option{
			IgnoreEmpty:   true,
			DeepCopy:      true,
			CaseSensitive: true,
		})
		if err != nil {
			return nil, errorx.GetFailedTransErr(err)
		}

		loan.RemoveRelationships()
	}
	err = copier.CopyWithOption(loan, reqData, copier.Option{
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
	creator, err := l.svcCtx.UserRepo.GetUser(l.ctx, reqData.UserId, constants.UserTypeBorrower)
	if err != nil {
		return nil, err
	}

	loan.CreatorId = reqData.UserId
	loan.CreatorName = creator.FullName
	loan.BorrowerOrganizationId = reqData.OrganizationId
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
		return nil, errorx.WithStack(errorx.New("Invalid Loan Borrower"))
	}
	loanBorrower, err := l.svcCtx.BorrowerEntityRepo.TranToLoanBorrowerWithId(l.ctx, loanBorrowerData.BorrowerEntityId, reqData.OrganizationId)
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
		borrowerApplicant, err := l.svcCtx.BorrowerApplicantRepo.GetBorrowerApplicantById(l.ctx, applicant.BorrowerApplicantId, reqData.OrganizationId)
		if err != nil {
			return nil, err
		}
		applicantItem := &model.LoanApplicant{}
		err = copier.CopyWithOption(applicantItem, borrowerApplicant, copier.Option{
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
		applicantItem.Id = applicant.Id
		applicantItem.IsPrimary = applicant.IsPrimary
		applicantItem.BorrowerApplicantId = applicant.BorrowerApplicantId
		applicantItem.OwnPercent = applicant.OwnPercent
		applicantItem.InvitationStatus = constants.ApplicantInvitationStatusAccepted
		applicants = append(applicants, applicantItem)
	}
	if totalOwnPercent > 100 {
		return nil, errorx.New("Total Own Percent Is Greater Than 100")
	}

	slice.SortBy(applicants, func(a, b *model.LoanApplicant) bool {
		return b.IsPrimary
	})
	primaryCount := slice.CountBy(applicants, func(index int, item *model.LoanApplicant) bool {
		return item.IsPrimary
	})
	if primaryCount == 0 {
		return nil, errorx.New("No Primary Applicant")
	}
	if primaryCount > 1 {
		return nil, errorx.New("Multiple Primary Applicant")
	}

	// loan_property
	borrowerProperty, err = l.svcCtx.BorrowerPropertyRepo.GetOneByIdForBorrower(l.ctx, reqData.LoanProperty.BorrowerPropertyId, reqData.UserId)
	if err != nil {
		return nil, err
	}
	loanProperty := &model.LoanProperty{}
	err = copier.CopyWithOption(loanProperty, reqData.LoanProperty, copier.Option{
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
	// loanProperty.Id = 0
	loanProperty.PropertyType = borrowerProperty.PropertyType
	loanProperty.Unit = borrowerProperty.Unit
	loanProperty.Bedrooms = borrowerProperty.Bedrooms
	loanProperty.Bathrooms = borrowerProperty.Bathrooms
	loanProperty.LivingAreaSqft = borrowerProperty.LivingAreaSqft
	loanProperty.Address = borrowerProperty.Address
	loanProperty.BorrowerPropertyId = reqData.LoanProperty.BorrowerPropertyId
	if loanProperty.YearlyLiabilityExpenseAmount > 0 {
		loanProperty.Dscr = mathutil.RoundToFloat((loanProperty.EstYearlyRentalIncomeAmount/loanProperty.YearlyLiabilityExpenseAmount)*100, 2)
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
		DeepCopy:      true,
		CaseSensitive: true,
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
		DeepCopy:      true,
		CaseSensitive: true,
		Converters: []copier.TypeConverter{
			copierx.StrToSqlTimePtr(timex.LayoutUsStyle),
		},
	})
	if err != nil {
		return nil, errorx.GetFailedTransErr(err)
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
