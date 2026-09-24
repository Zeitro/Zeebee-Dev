package project_management_workflow

import (
	"context"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProjectManagementWorkflowDefaultLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// set default template workflow
func NewProjectManagementWorkflowDefaultLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectManagementWorkflowDefaultLogic {
	return &ProjectManagementWorkflowDefaultLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProjectManagementWorkflowDefaultLogic) ProjectManagementWorkflowDefault(req *types.ProjectManagementWorkflowDetailRequest) (resp *types.EmptyResponse, err error) {
	data, err := copierx.Trans[types.ProjectManagementWorkflowDetailRequest, admin.ProjectManagementWorkflowDefaultRequestData](req)
	if err != nil {
		return nil, err
	}
	if _, err = l.svcCtx.ProjectManagementWorkflowRpcClient.Default(l.ctx, &admin.ProjectManagementWorkflowDefaultRequest{Data: data}); err != nil {
		return nil, err
	}
	return &types.EmptyResponse{}, nil
}
