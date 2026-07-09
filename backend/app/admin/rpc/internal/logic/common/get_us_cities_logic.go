package commonlogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/domains/zeebee/proto"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
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

func (l *GetUsCitiesLogic) GetUsCities(in *admin.UsCitiesRequest) (resp *proto.StringSliceResponse, err error) {
	defer func() {
		if err != nil {
			l.svcCtx.Logger.WithContext(l.ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
			).Error(err)
		}
	}()
	cities, err := l.svcCtx.UsGeographyRepo.GetCitiesByStateUsps(l.ctx, in.Data.StateUsps)
	if err != nil {
		return nil, err
	}
	return &proto.StringSliceResponse{
		Data: cities,
	}, nil
}
