package loanlogic

import (
	"context"

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
	"zeebee-backend/common/types"

	"zeebee-backend/common/errorx"
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
	applicants             []*model.Applicant
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
	data, err := l.buildModelData(in.Data)
	if err != nil {
		return nil, err
	}
	// apply loan
	var propertyBase *model.PropertyBase
	var loan *model.Loan
	var loanBorrower *model.LoanBorrower
	var loanProperty *model.LoanProperty
	err = l.svcCtx.LoanRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
		var innerErr error

		// get or create property base
		propertyBase, innerErr = l.svcCtx.PropertyBaseRepo.GetOrCreateWithTx(l.ctx, tx, data.loanProperty.PropertyBase)
		if innerErr != nil {
			return innerErr
		}

		// check loan
		if len(data.loan.Code) == 0 {
			innerErr = errors.WithStack(errors.New("code is empty"))
		}
		isExisted, innerErr := l.svcCtx.LoanRepo.IsExistedByColumnsWithTx(l.ctx, tx, map[string]interface{}{"code": data.loan.Code})
		if innerErr != nil {
			return innerErr
		}
		if isExisted {
			return errors.WithStack(errors.New("code already exists"))
		}

		// uid, innerErr := uuid.GetSonyflakeIDString()
		// if innerErr != nil {
		// 	return innerErr
		// }
		// data.loan.Code = uid
		// save loan
		loan, innerErr = l.svcCtx.LoanRepo.CreateWithTx(l.ctx, tx, data.loan, false)
		if innerErr != nil {
			return innerErr
		}

		// save loan borrower
		data.loanBorrower.LoanId = loan.Id
		data.loanBorrower.Applicants = nil
		loanBorrower, innerErr = l.svcCtx.LoanBorrowerRepo.CreateWithTx(l.ctx, tx, data.loanBorrower, false)
		if innerErr != nil {
			return innerErr
		}

		// save loan property
		data.loanProperty.LoanId = loan.Id
		data.loanProperty.PropertyBaseId = propertyBase.Id
		loanProperty, innerErr = l.svcCtx.LoanPropertyRepo.CreateWithTx(l.ctx, tx, data.loanProperty, false)
		if innerErr != nil {
			return innerErr
		}

		// save applicants
		slice.ForEach(data.applicants, func(index int, item *model.Applicant) {
			item.BorrowerId = in.Data.BorrowerId
			item.LoanId = loan.Id
			item.LoanBorrowerId = loanBorrower.Id
			item.LoanPropertyId = loanProperty.Id
		})
		_, innerErr = l.svcCtx.ApplicantRepo.BatchCreateWithTx(l.ctx, tx, data.applicants, false)
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
		innerErr = l.svcCtx.PropertyBaseRepo.SyncTotalAmountAndLoanCountWithTx(l.ctx, tx, propertyBase.Id)

		return innerErr
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to apply loan")
	}

	return &borrower.LoanApplyResponse{
		Data: &borrower.LoanApplyResponseData{
			Id: loan.Id,
		},
	}, nil
}

// buildModelData return save model data
func (l *ApplyLoanLogic) buildModelData(reqData *borrower.LoanApplyRequestData) (data *applyLoanModelData, err error) {

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
	loan := &model.Loan{
		// BorrowerId:    reqData.BorrowerId,
		// Stage:         constants.LoanStagePending,
		// PendingStatus: constants.LoanPendingStatusQuoting,
	}
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
	loan.Stage = constants.LoanStagePending
	loan.PendingStatus = constants.LoanPendingStatusQuoting
	loan.Source = constants.LoanSourceSystem
	// get lender
	lender, err := l.svcCtx.LenderCompanyRepo.FindOneById(l.ctx, reqData.LenderCompanyId)
	if err != nil {
		return nil, err
	}
	loan.LenderCompanyId = lender.Id
	loan.LenderCompany = lender.Name

	// loan_borrower
	loanBorrower := &model.LoanBorrower{
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
		item.BorrowerId = reqData.BorrowerId
		item.FullName = stringx.FormatFullName(item.FirstName, item.MidName, item.LastName)
	})

	// loan_property
	loanProperty := &model.LoanProperty{
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

	// transaction experience
	transactionExperiences := make([]*model.LoanExperience, 0)
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

	// real_estates experience
	realEstatesExperiences := make([]*model.LoanExperience, 0)
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

	// third services
	loanThirdServices := make([]*model.LoanThirdService, 0)
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

	// demographic
	loanDemographic := &model.LoanDemographic{}
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

	// apply form file
	// var applyFormFile *model.File
	// if len(reqData.ApplyFormFile) > 0 {
	// 	applyFormFile = &model.File{
	// 		File:     reqData.ApplyFormFile,
	// 		FileType: constants.FileTypeApplyForm,
	// 		BizType:  constants.FileBizTypeLoan,
	// 	}
	// }

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
