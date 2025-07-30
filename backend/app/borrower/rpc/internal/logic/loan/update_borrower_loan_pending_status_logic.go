package loanlogic

import (
	"context"

	"github.com/pkg/errors"

	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/constants"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type UpdateBorrowerLoanPendingStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateBorrowerLoanPendingStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateBorrowerLoanPendingStatusLogic {
	return &UpdateBorrowerLoanPendingStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateBorrowerLoanPendingStatusLogic) UpdateBorrowerLoanPendingStatus(in *borrower.UpdateLoanPendingStatusRequest) (resp *borrower.EmptyResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	loanModel, err := l.svcCtx.LoanRepo.FindOneByColumns(l.ctx, map[string]interface{}{
		"id":          in.Data.Id,
		"borrower_id": in.Data.BorrowerId,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get loan")
	}
	if loanModel.Stage != constants.LoanStagePending {
		return nil, errors.New("loan's stage should be pending")
	}

	if loanModel.Source != constants.LoanSourceBorrower {
		return nil, errors.WithStack(errors.New("loan source is not 3rd party"))
	}
	err = l.svcCtx.LoanRepo.UpdateColumnsById(l.ctx, in.Data.Id, map[string]interface{}{"pending_status": in.Data.PendingStatus})

	return &borrower.EmptyResponse{}, err
}
