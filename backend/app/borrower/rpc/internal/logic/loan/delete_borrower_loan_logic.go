package loanlogic

import (
	"context"

	"github.com/pkg/errors"
	"gorm.io/gorm"

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
		// delete applicants
		for _, applicant := range loanModel.Applicants {
			_, innerErr = l.svcCtx.ApplicantRepo.DeleteByIDWithTx(l.ctx, tx, applicant.Id)
			if innerErr != nil {
				return innerErr
			}
		}
		// delete property
		property := loanModel.Property
		if property != nil {
			var loanCount int64
			loanCount, innerErr = l.svcCtx.LoanRepo.GetLoanCountByPropertyId(l.ctx, property.Id)
			if innerErr != nil {
				return innerErr
			}
			if loanCount == 1 {
				_, innerErr = l.svcCtx.PropertyRepo.DeleteByIDWithTx(l.ctx, tx, property.Id)
				if innerErr != nil {
					return innerErr
				}
			}
		}
		// delete loan
		_, innerErr = l.svcCtx.LoanRepo.DeleteByIDWithTx(l.ctx, tx, loanModel.Id)

		if innerErr != nil {
			return innerErr
		}

		if property != nil {
			innerErr = l.svcCtx.PropertyRepo.SyncTotalAmountWithTx(l.ctx, tx, property.Id)
		}

		return innerErr
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed to delete loan")
	}
	return &borrower.EmptyResponse{}, nil
}
