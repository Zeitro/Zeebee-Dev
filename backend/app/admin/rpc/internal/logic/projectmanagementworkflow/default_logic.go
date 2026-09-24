package projectmanagementworkflowlogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/proto"
)

type DefaultLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewDefaultLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DefaultLogic {
	return &DefaultLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *DefaultLogic) Default(in *admin.ProjectManagementWorkflowDefaultRequest) (resp *proto.EmptyResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	claim, err := l.baseLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	if err = l.baseLogic.ProjectManagementTemplateLogic.SetWorkflowDefault(l.ctx, claim, in.Data.TemplateId, in.Data.Id); err != nil {
		return nil, err
	}
	return &proto.EmptyResponse{}, nil
}
