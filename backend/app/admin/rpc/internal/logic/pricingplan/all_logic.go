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

type AllLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewAllLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AllLogic {
	return &AllLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *AllLogic) All(in *admin.PricingPlanAllRequest) (resp *admin.PricingPlanAllResponse, err error) {

	if in.Data == nil {
		return nil, errorx.New("invalid request: data is required")
	}

	plans, err := l.svcCtx.TenantClient.PricingPlan.All(l.ctx, &tenant.PricingPlanAllRequest{Data: &tenant.PricingPlanAllRequestData{
		UserType: string(constants.UserTypeAdmin),
		Type:     in.Data.Type,
	}})
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	return &admin.PricingPlanAllResponse{Data: plans.Data}, nil
}
