package pricingplanpolicylogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/constants"
)

type CreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{ctx: ctx, svcCtx: svcCtx}
}

// Create proxies to tenant; tenant is the RBAC/权益 single authority.
func (l *CreateLogic) Create(in *admin.PricingPlanPolicyCreateRequest) (resp *admin.PricingPlanPolicyDetailResponse, err error) {
	data, err := copierx.Trans[admin.PricingPlanPolicyCreateRequestData, tenant.PricingPlanPolicyCreateRequestData](in.Data)
	if err != nil {
		return nil, err
	}
	data.UserType = string(constants.UserTypeAdmin)
	tenantResp, err := l.svcCtx.TenantClient.PricingPlanPolicy.Create(l.ctx, &tenant.PricingPlanPolicyCreateRequest{Data: data})
	if err != nil {
		return nil, err
	}
	return &admin.PricingPlanPolicyDetailResponse{Data: tenantResp.Data}, nil
}
