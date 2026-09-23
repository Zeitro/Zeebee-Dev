package projectmanagementtemplatelogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/common/errorx"
)

type WorkflowDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWorkflowDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WorkflowDetailLogic {
	return &WorkflowDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WorkflowDetailLogic) WorkflowDetail(in *admin.ProjectManagementTemplateWorkflowDetailRequest) (resp *admin.ProjectManagementTemplateWorkflowDetailResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	base := logic.GetBaseLogic(l.svcCtx)
	claim, err := base.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	workflow, err := base.ProjectManagementTemplateLogic.WorkflowDetail(l.ctx, claim, in.Data.TemplateId, in.Data.WorkflowId)
	if err != nil {
		return nil, err
	}
	return templateWorkflowResponse(workflow)
}
