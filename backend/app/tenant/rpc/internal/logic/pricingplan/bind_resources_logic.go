package pricingplanlogic

import (
	"context"

	"zeebee-backend/app/tenant/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/domains/zeebee/proto"

	"zeebee-backend/common/errorx"
)

type BindResourcesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBindResourcesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BindResourcesLogic {
	return &BindResourcesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BindResourcesLogic) BindResources(in *tenant.PricingPlanBindResourcesRequest) (resp *proto.EmptyResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	// todo: add your logic here and delete this line

	return &proto.EmptyResponse{}, nil
}
