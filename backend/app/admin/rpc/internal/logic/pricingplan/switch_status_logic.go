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

type SwitchStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewSwitchStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SwitchStatusLogic {
	return &SwitchStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *SwitchStatusLogic) SwitchStatus(in *admin.PricingPlanSwitchStatusRequest) (resp *proto.EmptyResponse, err error) {

	return l.svcCtx.TenantClient.PricingPlan.SwitchStatus(l.ctx, &tenant.PricingPlanSwitchStatusRequest{Data: &tenant.PricingPlanSwitchStatusRequestData{
		UserType: string(constants.UserTypeAdmin),
		Id:       in.Data.Id,
		Status:   in.Data.Status,
	}})
}
