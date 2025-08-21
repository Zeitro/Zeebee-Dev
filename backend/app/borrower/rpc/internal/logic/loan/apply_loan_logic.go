package loanlogic

import (
	"context"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"zeebee-backend/app/borrower/rpc/internal/repo/model"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/constants"
	"zeebee-backend/common/stringx"
	"zeebee-backend/common/types"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type ApplyLoanLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

type applyLoanModelData struct {
	loan                     *model.Loan
	loanBorrower             *model.LoanBorrower
	loanProperty             *model.LoanProperty
	loanDemographic          *model.LoanDemographic
	loanThirdServices        []*model.LoanThirdService
	applicants               []*model.Applicant
	experienceThirdServices  []*model.LoanThirdService
	realEstatesThirdServices []*model.LoanThirdService
}

func NewApplyLoanLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApplyLoanLogic {
	return &ApplyLoanLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ApplyLoanLogic) ApplyLoan(in *borrower.ApplyLoanRequest) (resp *borrower.ApplyLoanResponse, err error) {
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
	err = l.svcCtx.LoanRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
		var innerErr error
		var propertyBase *model.PropertyBase
		var loan *model.Loan
		var loanBorrower *model.LoanBorrower

		// get or create property base
		propertyBase, innerErr = l.svcCtx.PropertyBaseRepo.GetOrCreateWithTx(l.ctx, tx, data.loanProperty.PropertyBase)
		if innerErr != nil {
			return innerErr
		}

		// save loan
		loan, innerErr = l.svcCtx.LoanRepo.CreateWithTx(l.ctx, tx, data.loan)
		if innerErr != nil {
			return innerErr
		}

		// save loan borrower
		data.loanBorrower.LoanId = loan.Id
		loanBorrower, innerErr = l.svcCtx.LoanBorrowerRepo.CreateWithTx(l.ctx, tx, data.loanBorrower)
		if innerErr != nil {
			return innerErr
		}

		// save applicants
		slice.ForEach(data.applicants, func(index int, item *model.Applicant) {
			item.BorrowerId = in.Data.BorrowerId
			item.LoanId = loan.Id
			item.LoanBorrowerId = loanBorrower.Id
		})
		_, innerErr = l.svcCtx.ApplicantRepo.BatchCreateWithTx(l.ctx, tx, data.applicants)
		if innerErr != nil {
			return innerErr
		}

		// save loan property
		data.loanProperty.LoanId = loan.Id
		data.loanProperty.PropertyBaseId = propertyBase.Id
		_, innerErr = l.svcCtx.LoanPropertyRepo.CreateWithTx(l.ctx, tx, data.loanProperty)
		if innerErr != nil {
			return innerErr
		}

		// save demographic
		data.loanDemographic.LoanId = loan.Id
		_, innerErr = l.svcCtx.LoanDemographicRepo.CreateWithTx(l.ctx, tx, data.loanDemographic)
		if innerErr != nil {
			return innerErr
		}

		// save third service
		slice.ForEach(data.experienceThirdServices, func(index int, item *model.LoanThirdService) {
			item.LoanId = loan.Id
		})
		_, innerErr = l.svcCtx.LoanThirdServiceRepo.BatchCreateWithTx(l.ctx, tx, data.experienceThirdServices)
		if innerErr != nil {
			return innerErr
		}
		slice.ForEach(data.realEstatesThirdServices, func(index int, item *model.LoanThirdService) {
			item.LoanId = loan.Id
		})
		_, innerErr = l.svcCtx.LoanThirdServiceRepo.BatchCreateWithTx(l.ctx, tx, data.realEstatesThirdServices)
		if innerErr != nil {
			return innerErr
		}

		// sync total amount
		innerErr = l.svcCtx.PropertyBaseRepo.SyncTotalAmountWithTx(l.ctx, tx, propertyBase.Id)

		return innerErr
	})

	return &borrower.ApplyLoanResponse{}, nil
}

// buildModelData return save model data
func (l *ApplyLoanLogic) buildModelData(reqData *borrower.ApplyLoanRequestData) (data *applyLoanModelData, err error) {

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
	err = copier.CopyWithOption(loan, reqData, copier.Option{IgnoreEmpty: true, DeepCopy: true})
	if err != nil {
		return nil, errorx.GetFailedTransErr(err)
	}
	loan.BorrowerId = reqData.BorrowerId
	loan.Stage = constants.LoanStagePending
	loan.PendingStatus = constants.LoanPendingStatusQuoting
	lender, err := l.svcCtx.LenderRepo.FindOneById(l.ctx, reqData.LenderId)
	if err != nil {
		return nil, err
	}
	loan.LenderId = lender.Id
	loan.LenderCompany = lender.Name

	// loan_borrower
	loanBorrower := &model.LoanBorrower{
		MailingAddress: types.Address{},
	}
	err = copier.CopyWithOption(loanBorrower, reqData.LoanBorrower, copier.Option{IgnoreEmpty: true, DeepCopy: true, CaseSensitive: true})
	if err != nil {
		return nil, errorx.GetFailedTransErr(err)
	}

	// applicants
	applicants := make([]*model.Applicant, 0)
	err = copier.CopyWithOption(applicants, reqData.Applicants, copier.Option{IgnoreEmpty: true, DeepCopy: true})
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
	err = copier.CopyWithOption(loanProperty, reqData.LoanProperty, copier.Option{IgnoreEmpty: true, DeepCopy: true})
	if err != nil {
		return nil, errorx.GetFailedTransErr(err)
	}
	loanProperty.PropertyBase.BorrowerId = reqData.BorrowerId

	// experience_third_services
	experienceThirdServices := make([]*model.LoanThirdService, 0)
	err = copier.CopyWithOption(experienceThirdServices, reqData.ExperienceServices, copier.Option{IgnoreEmpty: true, DeepCopy: true})
	if err != nil {
		return nil, errorx.GetFailedTransErr(err)
	}

	// real_estates_third_services
	realEstatesThirdServices := make([]*model.LoanThirdService, 0)
	err = copier.CopyWithOption(realEstatesThirdServices, reqData.RealEstateServices, copier.Option{IgnoreEmpty: true, DeepCopy: true})
	if err != nil {
		return nil, errorx.GetFailedTransErr(err)
	}

	// demographic
	loanDemographic := &model.LoanDemographic{}
	err = copier.CopyWithOption(loanDemographic, reqData.LoanDemographic, copier.Option{IgnoreEmpty: true, DeepCopy: true})
	if err != nil {
		return nil, errorx.GetFailedTransErr(err)
	}
	data = &applyLoanModelData{
		loan:                     loan,
		loanProperty:             loanProperty,
		loanBorrower:             loanBorrower,
		applicants:               applicants,
		loanThirdServices:        experienceThirdServices,
		realEstatesThirdServices: realEstatesThirdServices,
		loanDemographic:          loanDemographic,
	}
	return
}
