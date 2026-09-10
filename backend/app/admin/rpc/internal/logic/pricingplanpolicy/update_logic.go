package pricingplanpolicylogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/constants"
)

type UpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{ctx: ctx, svcCtx: svcCtx}
}

// Update proxies to tenant; tenant is the RBAC/权益 single authority.
func (l *UpdateLogic) Update(in *admin.PricingPlanPolicyUpdateRequest) (resp *admin.PricingPlanPolicyDetailResponse, err error) {
	data, err := copierx.Trans[admin.PricingPlanPolicyUpdateRequestData, tenant.PricingPlanPolicyUpdateRequestData](in.Data)
	if err != nil {
		return nil, err
	}
	data.UserType = string(constants.UserTypeAdmin)
	tenantResp, err := l.svcCtx.TenantClient.PricingPlanPolicy.Update(l.ctx, &tenant.PricingPlanPolicyUpdateRequest{Data: data})
	if err != nil {
		return nil, err
	}
	return &admin.PricingPlanPolicyDetailResponse{Data: tenantResp.Data}, nil
}
