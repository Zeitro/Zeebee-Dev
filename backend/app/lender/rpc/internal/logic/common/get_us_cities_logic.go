package commonlogic

import (
	"context"

	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/domains/zeebee/proto"
)

type GetUsCitiesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUsCitiesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUsCitiesLogic {
	return &GetUsCitiesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUsCitiesLogic) GetUsCities(in *lender.UsCitiesRequest) (resp *proto.StringSliceResponse, err error) {

	cities, err := l.svcCtx.UsGeographyRepo.GetCitiesByStateUsps(l.ctx, in.Data.StateUsps)
	if err != nil {
		return nil, err
	}
	return &proto.StringSliceResponse{
		Data: cities,
	}, nil
}
