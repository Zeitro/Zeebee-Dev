package project_management_workflow

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProjectManagementProjectWorkflowTemplateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// discover project templates
func NewProjectManagementProjectWorkflowTemplateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectManagementProjectWorkflowTemplateLogic {
	return &ProjectManagementProjectWorkflowTemplateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProjectManagementProjectWorkflowTemplateLogic) ProjectManagementProjectWorkflowTemplate(req *types.ProjectManagementWorkflowTemplateRequest) (resp *types.ProjectManagementWorkflowTemplateResponse, err error) {
	rpcResp, err := l.svcCtx.ProjectManagementWorkflowRpcClient.ProjectWorkflowTemplates(l.ctx, &proto.EmptyRequest{})
	if err != nil {
		return nil, err
	}
	list, err := copierx.Trans[[]*proto.ProjectManagementTemplateEntity, []*types.ProjectManagementTemplateEntity](&rpcResp.List)
	if err != nil {
		return nil, err
	}
	return &types.ProjectManagementWorkflowTemplateResponse{List: *list}, nil
}
