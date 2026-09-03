package pricingplanpolicylogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/proto"
)

type DeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogic {
	return &DeleteLogic{ctx: ctx, svcCtx: svcCtx}
}

// Delete proxies to tenant; tenant is the RBAC/权益 single authority.
func (l *DeleteLogic) Delete(in *proto.PathIdRequest) (resp *proto.EmptyResponse, err error) {
	return l.svcCtx.TenantClient.PricingPlanPolicy.Delete(l.ctx, &tenant.PricingPlanPolicyDeleteRequest{
		Data: &tenant.PricingPlanPolicyDeleteRequestData{
			UserType: string(constants.UserTypeAdmin),
			Id:       in.Data.Id,
		},
	})
}
