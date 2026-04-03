package loan

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/copierx"
)

type UpdatePendingStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// update loan pending status
func NewUpdatePendingStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePendingStatusLogic {
	return &UpdatePendingStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdatePendingStatusLogic) UpdatePendingStatus(req *types.UpdateLoanPendingStatusRequest) (resp *types.EmptyResponse, err error) {
	rpcReq := &lender.LoanUpdatePendingStatusRequest{
		Data: &lender.LoanUpdatePendingStatusRequestData{},
	}
	resp = &types.EmptyResponse{}

	rpcReq.Data, err = copierx.Trans[types.UpdateLoanPendingStatusRequest, lender.LoanUpdatePendingStatusRequestData](req)
	if err != nil {
		return nil, err
	}
	_, err = l.svcCtx.LoanRpcClient.UpdateLoanPendingStatus(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
