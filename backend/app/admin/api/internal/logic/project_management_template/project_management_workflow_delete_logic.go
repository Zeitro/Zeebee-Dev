package project_management_template

import (
	"context"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProjectManagementWorkflowDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// delete
func NewProjectManagementWorkflowDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectManagementWorkflowDeleteLogic {
	return &ProjectManagementWorkflowDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProjectManagementWorkflowDeleteLogic) ProjectManagementWorkflowDelete(req *types.ProjectManagementWorkflowDetailRequest) (resp *types.EmptyResponse, err error) {
	data, err := copierx.Trans[types.ProjectManagementWorkflowDetailRequest, admin.ProjectManagementTemplateWorkflowDeleteRequestData](req)
	if err != nil {
		return nil, err
	}
	if _, err = l.svcCtx.ProjectManagementTemplateRpcClient.WorkflowDelete(l.ctx, &admin.ProjectManagementTemplateWorkflowDeleteRequest{Data: data}); err != nil {
		return nil, err
	}
	return &types.EmptyResponse{}, nil
}
