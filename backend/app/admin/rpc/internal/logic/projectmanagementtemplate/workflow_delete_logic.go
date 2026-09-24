package projectmanagementtemplatelogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/proto"
)

type WorkflowDeleteLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewWorkflowDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WorkflowDeleteLogic {
	return &WorkflowDeleteLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *WorkflowDeleteLogic) WorkflowDelete(in *admin.ProjectManagementTemplateWorkflowDeleteRequest) (resp *proto.EmptyResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	base := l.baseLogic
	claim, err := base.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	if err = base.ProjectManagementTemplateLogic.DeleteWorkflow(l.ctx, claim, in.Data.TemplateId, in.Data.WorkflowId); err != nil {
		return nil, err
	}
	return &proto.EmptyResponse{}, nil
}
