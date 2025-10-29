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
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/stringx"
	"zeebee-backend/common/timex"
	"zeebee-backend/common/uuid"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type importPendingLoanModelData struct {
	loan         *model.Loan
	loanBorrower *model.LoanBorrower
	loanProperty *model.LoanProperty
	propertyBase *model.PropertyBase
	applicants   []*model.Applicant
	// loanExtend   *model.LoanExtend
}

type ImportPendingLoansLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewImportPendingLoansLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ImportPendingLoansLogic {
	return &ImportPendingLoansLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ImportPendingLoansLogic) ImportPendingLoans(in *borrower.ImportPendingLoansRequest) (resp *borrower.ImportPendingLoansResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	// err = errors.Errorf("no pending loan")
	// return nil, errors.WithStack(err)

	// get save data
	models, err := l.buildImportModels(in.Data)
	if err != nil {
		return nil, err
	}

	// 创建贷款
	err = l.svcCtx.LoanRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
		var innerErr error

		for i := range models {
			var propertyBase *model.PropertyBase
			var uid string
			data := models[i]
			// get property base
			propertyBase, innerErr = l.svcCtx.PropertyBaseRepo.GetOrCreateWithTx(l.ctx, tx, data.propertyBase)
			if innerErr != nil {
				return innerErr
			}

			// loan
			uid, innerErr = uuid.GetSonyflakeIDString()
			if innerErr != nil {
				return innerErr
			}
			loan := data.loan
			loan.Code = uid
			loan.RemoveRelationships()
			loan, innerErr = l.svcCtx.LoanRepo.CreateWithTx(l.ctx, tx, loan, false)
			if innerErr != nil {
				return innerErr
			}
			// loan property
			loanProperty := data.loanProperty
			loanProperty.LoanId = loan.Id
			loanProperty.PropertyBaseId = propertyBase.Id
			loanProperty, innerErr = l.svcCtx.LoanPropertyRepo.CreateWithTx(l.ctx, tx, loanProperty, false)
			if innerErr != nil {
				return innerErr
			}

			// loan borrower
			loanBorrower := data.loanBorrower
			loanBorrower.LoanId = loan.Id
			loanBorrower, innerErr = l.svcCtx.LoanBorrowerRepo.CreateWithTx(l.ctx, tx, loanBorrower, false)
			if innerErr != nil {
				return innerErr
			}

			// applicants
			applicants := data.applicants
			applicants = slice.Map(applicants, func(index int, item *model.Applicant) *model.Applicant {
				item.LoanId = loan.Id
				item.LoanPropertyId = loanProperty.Id
				item.LoanBorrowerId = loanBorrower.Id
				item.FullName = stringx.FormatFullName(item.FirstName, item.MidName, item.LastName)

				return item
			})
			_, innerErr = l.svcCtx.ApplicantRepo.BatchCreateWithTx(l.ctx, tx, applicants, false)
			if innerErr != nil {
				return innerErr
			}
		}

		return innerErr
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to import pending loan")
	}

	return &borrower.ImportPendingLoansResponse{
		Data: &borrower.ImportPendingLoansResponseData{
			RecordCount: int64(len(models)),
		},
	}, nil
}

// buildImportModels gets the models for importing a pending loan.
func (l *ImportPendingLoansLogic) buildImportModels(in []*borrower.ImportPendingLoansRequestData) (data []*importPendingLoanModelData, err error) {
	data = make([]*importPendingLoanModelData, 0, len(in))
	for item := range in {
		item, err := l.buildImportModel(in[item])
		if err != nil {
			return nil, err
		}
		data = append(data, item)
	}
	return data, nil
}

func (l *ImportPendingLoansLogic) buildImportModel(in *borrower.ImportPendingLoansRequestData) (data *importPendingLoanModelData, err error) {
	// loan
	loan := &model.Loan{}

	err = copier.CopyWithOption(loan, &in, copier.Option{
		IgnoreEmpty: true,
		DeepCopy:    true,
		Converters: []copier.TypeConverter{
			copierx.StrToSqlTimePtr(timex.LayoutUsStyle),
		},
	})
	if err != nil {
		return nil, errorx.GetFailedTransErr(err)
	}
	loan.BorrowerId = in.BorrowerId
	loan.LoanType = constants.LoanType(in.LoanType)
	loan.Stage = constants.LoanStagePending
	loan.PendingStatus = constants.LoanPendingStatusQuoting
	loan.Source = constants.LoanSourceBorrower

	// loan borrower
	loanBorrower := &model.LoanBorrower{}
	err = copier.CopyWithOption(loanBorrower, in.LoanBorrower, copier.Option{
		IgnoreEmpty:   true,
		CaseSensitive: true,
		DeepCopy:      true,
		Converters: []copier.TypeConverter{
			copierx.StrToSqlTimePtr(timex.LayoutUsStyle),
		},
	})
	if err != nil {
		return nil, errorx.GetFailedTransErr(err)
	}

	// applicants
	var applicants []*model.Applicant

	err = copier.CopyWithOption(&applicants, in.Applicants, copier.Option{
		IgnoreEmpty:   true,
		CaseSensitive: true,
		DeepCopy:      true,
		Converters: []copier.TypeConverter{
			copierx.StrToSqlTimePtr(timex.LayoutUsStyle),
		},
	})
	if err != nil {
		return nil, errorx.GetFailedTransErr(err)
	}
	if len(applicants) == 0 {
		return nil, errors.New("no applicant")
	}
	primaryCount := slice.CountBy(applicants, func(index int, item *model.Applicant) bool {
		return item.IsPrimary
	})
	if primaryCount == 0 {
		return nil, errors.New("no primary applicant")
	}
	slice.ForEach(applicants, func(index int, item *model.Applicant) {
		item.BorrowerId = in.BorrowerId
	})

	// loan property
	loanProperty := &model.LoanProperty{}
	err = copier.CopyWithOption(loanProperty, in.LoanProperty, copier.Option{
		IgnoreEmpty:   true,
		CaseSensitive: true,
		DeepCopy:      true,
		Converters: []copier.TypeConverter{
			copierx.StrToSqlTimePtr(timex.LayoutUsStyle),
		},
	})
	if err != nil {
		return nil, errorx.GetFailedTransErr(err)
	}
	loanProperty.Unit = loanProperty.PropertyBase.Unit

	// property base
	propertyBase := loanProperty.PropertyBase

	return &importPendingLoanModelData{
		loan:         loan,
		loanBorrower: loanBorrower,
		loanProperty: loanProperty,
		propertyBase: propertyBase,
		applicants:   applicants,
	}, nil
}
