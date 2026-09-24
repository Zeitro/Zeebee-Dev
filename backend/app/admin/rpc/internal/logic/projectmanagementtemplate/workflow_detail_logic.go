package projectmanagementtemplatelogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/proto"
)

type WorkflowDetailLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewWorkflowDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WorkflowDetailLogic {
	return &WorkflowDetailLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *WorkflowDetailLogic) WorkflowDetail(in *admin.ProjectManagementTemplateWorkflowDetailRequest) (resp *admin.ProjectManagementTemplateWorkflowDetailResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	base := l.baseLogic
	claim, err := base.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	workflow, err := base.ProjectManagementTemplateLogic.WorkflowDetail(l.ctx, claim, in.Data.TemplateId, in.Data.WorkflowId)
	if err != nil {
		return nil, err
	}
	data, err := copierx.Trans[model.ProjectManagementWorkflow, proto.ProjectManagementWorkflowEntity](workflow, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}
	return &admin.ProjectManagementTemplateWorkflowDetailResponse{Data: data}, nil
}
