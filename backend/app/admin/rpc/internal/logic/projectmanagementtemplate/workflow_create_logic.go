package projectmanagementtemplatelogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/data/dto"
)

type WorkflowCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWorkflowCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WorkflowCreateLogic {
	return &WorkflowCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WorkflowCreateLogic) WorkflowCreate(in *admin.ProjectManagementTemplateWorkflowCreateRequest) (resp *admin.ProjectManagementTemplateWorkflowDetailResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	base := logic.GetBaseLogic(l.svcCtx)
	claim, err := base.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	transitions, err := copierx.Trans[[]*admin.ProjectManagementTemplateWorkflowTransitionInput, []dto.ProjectManagementWorkflowTransitionDto](&in.Data.Transitions)
	if err != nil {
		return nil, err
	}
	workflow, err := base.ProjectManagementTemplateLogic.CreateWorkflow(l.ctx, claim, in.Data.TemplateId, dto.ProjectManagementWorkflowDto{Name: in.Data.Name, Description: in.Data.Description, Transitions: *transitions})
	if err != nil {
		return nil, err
	}
	return templateWorkflowResponse(workflow)
}
