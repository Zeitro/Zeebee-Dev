package pricingplanlogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
)

type UpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *UpdateLogic) Update(in *admin.PricingPlanUpdateRequest) (resp *admin.PricingPlanDetailResponse, err error) {

	data, err := copierx.Trans[admin.PricingPlanUpdateRequestData, tenant.PricingPlanUpdateRequestData](in.Data)
	if err != nil {
		return nil, err
	}

	data.UserType = string(constants.UserTypeAdmin)
	result, err := l.svcCtx.TenantClient.PricingPlan.Update(l.ctx, &tenant.PricingPlanUpdateRequest{Data: data})
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	return &admin.PricingPlanDetailResponse{Data: result.Data}, nil
}
