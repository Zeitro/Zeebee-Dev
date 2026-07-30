package loan

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
)

type UpdateStageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// update loan stage
func NewUpdateStageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateStageLogic {
	return &UpdateStageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateStageLogic) UpdateStage(req *types.UpdateLoanStageRequest) (resp *types.EmptyResponse, err error) {
	rpcReq := &lender.LoanUpdateStageRequest{
		Data: &lender.LoanUpdateStageRequestData{},
	}
	resp = &types.EmptyResponse{}

	rpcReq.Data, err = copierx.Trans[types.UpdateLoanStageRequest, lender.LoanUpdateStageRequestData](req)
	if err != nil {
		return nil, err
	}
	_, err = l.svcCtx.LoanRpcClient.UpdateLoanStage(l.ctx, rpcReq)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	return resp, nil
}
