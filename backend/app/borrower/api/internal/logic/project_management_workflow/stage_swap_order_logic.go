package project_management_workflow

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type StageSwapOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// swap stage order
func NewStageSwapOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StageSwapOrderLogic {
	return &StageSwapOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StageSwapOrderLogic) StageSwapOrder(req *types.ProjectManagementWorkflowStageSwapOrderRequest) (resp *types.EmptyResponse, err error) {
	data, err := copierx.Trans[types.ProjectManagementWorkflowStageSwapOrderRequest, borrower.ProjectManagementWorkflowStageSwapOrderRequestData](req)
	if err != nil {
		return nil, err
	}
	if _, err = l.svcCtx.ProjectManagementWorkflowStageRpcClient.SwapOrder(l.ctx, &borrower.ProjectManagementWorkflowStageSwapOrderRequest{Data: data}); err != nil {
		return nil, errorx.WithStack(err)
	}
	return &types.EmptyResponse{}, nil
}
