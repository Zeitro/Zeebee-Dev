package projectmanagementtemplatelogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/domains/zeebee/proto"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/model"
)

type WorkflowsLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewWorkflowsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WorkflowsLogic {
	return &WorkflowsLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *WorkflowsLogic) Workflows(in *proto.PathIdRequest) (resp *admin.ProjectManagementTemplateWorkflowResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	base := l.baseLogic
	claim, err := base.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	workflows, err := base.ProjectManagementTemplateLogic.Workflows(l.ctx, claim, in.Data.Id)
	if err != nil {
		return nil, err
	}
	list := make([]*proto.ProjectManagementWorkflowEntity, 0, len(workflows))
	for _, workflow := range workflows {
		entity, convertErr := copierx.Trans[model.ProjectManagementWorkflow, proto.ProjectManagementWorkflowEntity](workflow, copierx.TimeToRFC3339Str())
		if convertErr != nil {
			return nil, convertErr
		}
		list = append(list, entity)
	}
	return &admin.ProjectManagementTemplateWorkflowResponse{List: list}, nil
}
