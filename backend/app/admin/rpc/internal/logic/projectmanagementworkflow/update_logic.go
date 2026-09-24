package projectmanagementworkflowlogic

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

type UpdateLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *UpdateLogic) Update(in *admin.ProjectManagementWorkflowUpdateRequest) (resp *admin.ProjectManagementWorkflowDetailResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	claim, err := l.baseLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	transitions, err := copierx.Trans[[]*admin.ProjectManagementWorkflowTransitionInput, []dto.ProjectManagementWorkflowTransitionDto](&in.Data.Transitions)
	if err != nil {
		return nil, err
	}
	workflow, err := l.baseLogic.ProjectManagementTemplateLogic.UpdateWorkflow(l.ctx, claim, in.Data.TemplateId, in.Data.Id, dto.ProjectManagementWorkflowDto{Name: in.Data.Name, Description: in.Data.Description, Transitions: *transitions})
	if err != nil {
		return nil, err
	}
	data, err := copierx.Trans[model.ProjectManagementWorkflow, proto.ProjectManagementWorkflowEntity](workflow, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}
	return &admin.ProjectManagementWorkflowDetailResponse{Data: data}, nil
}
