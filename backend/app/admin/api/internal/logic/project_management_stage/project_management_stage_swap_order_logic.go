package project_management_stage

import (
	"context"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProjectManagementStageSwapOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// swap stage order
func NewProjectManagementStageSwapOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectManagementStageSwapOrderLogic {
	return &ProjectManagementStageSwapOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProjectManagementStageSwapOrderLogic) ProjectManagementStageSwapOrder(req *types.ProjectManagementTemplateStageSwapOrderRequest) (resp *types.EmptyResponse, err error) {
	data, err := copierx.Trans[types.ProjectManagementTemplateStageSwapOrderRequest, admin.ProjectManagementStageSwapOrderRequestData](req)
	if err != nil {
		return nil, err
	}
	if _, err = l.svcCtx.ProjectManagementStageRpcClient.SwapOrder(l.ctx, &admin.ProjectManagementStageSwapOrderRequest{Data: data}); err != nil {
		return nil, err
	}
	return &types.EmptyResponse{}, nil
}
