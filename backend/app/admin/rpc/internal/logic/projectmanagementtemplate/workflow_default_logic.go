package projectmanagementtemplatelogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/proto"
)

type WorkflowDefaultLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWorkflowDefaultLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WorkflowDefaultLogic {
	return &WorkflowDefaultLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WorkflowDefaultLogic) WorkflowDefault(in *admin.ProjectManagementTemplateWorkflowDefaultRequest) (resp *proto.EmptyResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	base := logic.GetBaseLogic(l.svcCtx)
	claim, err := base.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	if err = base.ProjectManagementTemplateLogic.SetWorkflowDefault(l.ctx, claim, in.Data.TemplateId, in.Data.WorkflowId); err != nil {
		return nil, err
	}
	return &proto.EmptyResponse{}, nil
}
