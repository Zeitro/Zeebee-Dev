package pricingplanlogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
)

type ResourceListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewResourceListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResourceListLogic {
	return &ResourceListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *ResourceListLogic) ResourceList(in *admin.PricingPlanResourceListRequest) (resp *admin.PricingPlanResourceListResponse, err error) {

	resources, err := l.svcCtx.TenantClient.PricingPlan.ResourceList(l.ctx, &tenant.PricingPlanResourceListRequest{Data: &tenant.PricingPlanResourceListRequestData{
		UserType: string(constants.UserTypeAdmin),
		Id:       in.Data.Id,
	}})
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	return &admin.PricingPlanResourceListResponse{Data: resources.Data}, nil
}
