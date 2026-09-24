package projectmanagementworkflowlogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/proto"
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

func (l *WorkflowsLogic) Workflows(in *proto.PathIdRequest) (resp *admin.ProjectManagementWorkflowResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	claim, err := l.baseLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	workflows, err := l.baseLogic.ProjectManagementTemplateLogic.Workflows(l.ctx, claim, in.Data.Id)
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
	return &admin.ProjectManagementWorkflowResponse{List: list}, nil
}
