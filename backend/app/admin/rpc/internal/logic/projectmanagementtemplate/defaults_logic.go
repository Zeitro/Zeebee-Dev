package projectmanagementtemplatelogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/proto"
)

type DefaultsLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewDefaultsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DefaultsLogic {
	return &DefaultsLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *DefaultsLogic) Defaults(in *proto.PathIdRequest) (resp *proto.EmptyResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	base := l.baseLogic
	claim, err := base.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	if err = base.ProjectManagementTemplateLogic.SetDefault(l.ctx, claim, in.Data.Id); err != nil {
		return nil, err
	}
	return &proto.EmptyResponse{}, nil
}
