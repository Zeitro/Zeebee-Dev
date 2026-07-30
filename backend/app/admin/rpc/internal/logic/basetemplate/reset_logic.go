package basetemplatelogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/domains/zeebee/proto"
)

type ResetLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewResetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetLogic {
	return &ResetLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *ResetLogic) Reset(in *proto.IdRequest) (resp *proto.EmptyResponse, err error) {

	id := in.Data
	err = l.baseLogic.OrganizationTemplateLogic.ResetByBaseTemplateId(l.ctx, id)
	return &proto.EmptyResponse{}, err
}
