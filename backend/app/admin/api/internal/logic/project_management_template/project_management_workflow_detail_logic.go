package project_management_template

import (
	"context"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProjectManagementWorkflowDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// detail
func NewProjectManagementWorkflowDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectManagementWorkflowDetailLogic {
	return &ProjectManagementWorkflowDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProjectManagementWorkflowDetailLogic) ProjectManagementWorkflowDetail(req *types.ProjectManagementWorkflowDetailRequest) (resp *types.ProjectManagementWorkflowEntity, err error) {
	data, err := copierx.Trans[types.ProjectManagementWorkflowDetailRequest, admin.ProjectManagementTemplateWorkflowDetailRequestData](req)
	if err != nil {
		return nil, err
	}
	rpcResp, err := l.svcCtx.ProjectManagementTemplateRpcClient.WorkflowDetail(l.ctx, &admin.ProjectManagementTemplateWorkflowDetailRequest{Data: data})
	if err != nil {
		return nil, err
	}
	return copierx.Trans[proto.ProjectManagementWorkflowEntity, types.ProjectManagementWorkflowEntity](rpcResp.Data)
}
