package pricingplanlogic

import (
	"context"

	"zeebee-backend/app/tenant/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"

	"zeebee-backend/common/errorx"
)

type ResourceListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewResourceListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResourceListLogic {
	return &ResourceListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ResourceListLogic) ResourceList(in *tenant.PricingPlanResourceListRequest) (resp *tenant.PricingPlanResourceListResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	// todo: add your logic here and delete this line

	return &tenant.PricingPlanResourceListResponse{}, nil
}
