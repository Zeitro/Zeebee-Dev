package loanlogic

import (
	"context"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"zeebee-backend/app/borrower/rpc/internal/repo/model"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/constants"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type DeleteBorrowerLoanLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteBorrowerLoanLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteBorrowerLoanLogic {
	return &DeleteBorrowerLoanLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// DeleteBorrowerLoan delete borrower loan
func (l *DeleteBorrowerLoanLogic) DeleteBorrowerLoan(in *borrower.LoanDeleteRequest) (resp *borrower.EmptyResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	loanModel, err := l.svcCtx.LoanRepo.GetBorrowerLoan(l.ctx, in.Data.Id, in.Data.BorrowerId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get loan")
	}

	if loanModel.Source != constants.LoanSourceBorrower {
		return nil, errors.WithStack(errors.New("loan source is not 3rd party"))
	}

	err = l.svcCtx.LoanRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
		var innerErr error
		var loanProperty *model.LoanProperty
		var propertyBase *model.PropertyBase
		loanProperty = loanModel.LoanProperty
		if loanProperty != nil {
			propertyBase = loanProperty.PropertyBase
		}

		// delete property_base
		if propertyBase != nil {
			var loanCount int64
			loanCount, innerErr = l.svcCtx.LoanPropertyRepo.GetLoanPropertyCountByPropertyBaseId(l.ctx, propertyBase.Id)
			if innerErr != nil {
				return innerErr
			}
			if loanCount == 1 {
				_, innerErr = l.svcCtx.PropertyBaseRepo.DeleteByIDWithTx(l.ctx, tx, propertyBase.Id)
				if innerErr != nil {
					return innerErr
				}
			}
		}

		// delete loan_property
		_, innerErr = l.svcCtx.LoanPropertyRepo.DeleteByIDWithTx(l.ctx, tx, loanProperty.Id)
		if innerErr != nil {
			return innerErr
		}

		// delete applicants
		for _, applicant := range loanModel.Applicants {
			_, innerErr = l.svcCtx.ApplicantRepo.DeleteByIDWithTx(l.ctx, tx, applicant.Id)
			if innerErr != nil {
				return innerErr
			}
		}

		// delete loan
		_, innerErr = l.svcCtx.LoanRepo.DeleteByIDWithTx(l.ctx, tx, loanModel.Id)

		if innerErr != nil {
			return innerErr
		}

		if propertyBase != nil {
			innerErr = l.svcCtx.PropertyBaseRepo.SyncTotalAmountWithTx(l.ctx, tx, propertyBase.Id)
		}

		return innerErr
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed to delete loan")
	}
	return &borrower.EmptyResponse{}, nil
}
