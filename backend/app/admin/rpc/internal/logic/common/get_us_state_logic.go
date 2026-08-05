package commonlogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/domains/zeebee/proto"
)

type GetUsStateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUsStateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUsStateLogic {
	return &GetUsStateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUsStateLogic) GetUsState(in *proto.EmptyRequest) (resp *proto.StringSliceResponse, err error) {

	states, err := l.svcCtx.UsGeographyRepo.GetStateUsps(l.ctx)
	if err != nil {
		return nil, err
	}

	return &proto.StringSliceResponse{
		Data: states,
	}, nil
}
