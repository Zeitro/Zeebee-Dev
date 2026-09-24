package project_management_workflow

import (
	"context"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProjectManagementTemplateWorkflowsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// template's workflow
func NewProjectManagementTemplateWorkflowsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectManagementTemplateWorkflowsLogic {
	return &ProjectManagementTemplateWorkflowsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProjectManagementTemplateWorkflowsLogic) ProjectManagementTemplateWorkflows(req *types.ProjectManagementTemplateWorflowRequest) (resp *types.ProjectManagementTemplateWorflowResonse, err error) {
	rpcResp, err := l.svcCtx.ProjectManagementWorkflowRpcClient.Workflows(l.ctx, &proto.PathIdRequest{Data: &proto.PathIdRequestData{Id: req.TemplateId}})
	if err != nil {
		return nil, err
	}
	list, err := copierx.Trans[[]*proto.ProjectManagementWorkflowEntity, []types.ProjectManagementTemplateWorflowEntity](&rpcResp.List)
	if err != nil {
		return nil, err
	}
	return &types.ProjectManagementTemplateWorflowResonse{List: *list}, nil
}
