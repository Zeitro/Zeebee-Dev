package projectmanagementtemplatelogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/proto"
)

type WorkflowUpdateLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewWorkflowUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WorkflowUpdateLogic {
	return &WorkflowUpdateLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *WorkflowUpdateLogic) WorkflowUpdate(in *admin.ProjectManagementTemplateWorkflowUpdateRequest) (resp *admin.ProjectManagementTemplateWorkflowDetailResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	base := l.baseLogic
	claim, err := base.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	transitions, err := copierx.Trans[[]*admin.ProjectManagementTemplateWorkflowTransitionInput, []dto.ProjectManagementWorkflowTransitionDto](&in.Data.Transitions)
	if err != nil {
		return nil, err
	}
	workflow, err := base.ProjectManagementTemplateLogic.UpdateWorkflow(l.ctx, claim, in.Data.TemplateId, in.Data.WorkflowId, dto.ProjectManagementWorkflowDto{Name: in.Data.Name, Description: in.Data.Description, Transitions: *transitions})
	if err != nil {
		return nil, err
	}
	data, err := copierx.Trans[model.ProjectManagementWorkflow, proto.ProjectManagementWorkflowEntity](workflow, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}
	return &admin.ProjectManagementTemplateWorkflowDetailResponse{Data: data}, nil
}
