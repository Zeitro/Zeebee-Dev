package pricingplanlogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/proto"
)

type BindResourcesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewBindResourcesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BindResourcesLogic {
	return &BindResourcesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *BindResourcesLogic) BindResources(in *admin.PricingPlanBindResourcesRequest) (resp *proto.EmptyResponse, err error) {

	return l.svcCtx.TenantClient.PricingPlan.BindResources(l.ctx, &tenant.PricingPlanBindResourcesRequest{Data: &tenant.PricingPlanBindResourcesRequestData{
		UserType:    string(constants.UserTypeAdmin),
		Id:          in.Data.Id,
		ResourceIds: in.Data.ResourceIds,
	}})
}
