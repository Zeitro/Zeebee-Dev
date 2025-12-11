package loan

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatePendingStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// update pending status
func NewUpdatePendingStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePendingStatusLogic {
	return &UpdatePendingStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdatePendingStatusLogic) UpdatePendingStatus(req *types.LoanUpdatePendingStatusRequest) (resp *types.EmptyResponse, err error) {
	_, err = l.svcCtx.LoanRpcClient.UpdateBorrowerLoanPendingStatus(l.ctx, &borrower.LoanUpdatePendingStatusRequest{
		Data: &borrower.LoanUpdatePendingStatusRequestData{
			Id:            req.Id,
			PendingStatus: req.PendingStatus,
		},
	})

	return &types.EmptyResponse{}, err
}
