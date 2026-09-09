package userlogic

import (
	"context"

	"zeebee-backend/app/lender/rpc/internal/svc"
	tenantpb "zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/domains/zeebee/proto"
)

type SwitchOrganizationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSwitchOrganizationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SwitchOrganizationLogic {
	return &SwitchOrganizationLogic{ctx: ctx, svcCtx: svcCtx}
}

// SwitchOrganization proxies to tenant. user_id and userType are resolved by
// tenant from the caller's JWT (forwarded by the lender rpc server interceptor).
func (l *SwitchOrganizationLogic) SwitchOrganization(in *proto.SwitchOrganizationRequest) (resp *proto.UserLoginResponse, err error) {
	return l.svcCtx.TenantClient.User.SwitchOrganization(l.ctx, &tenantpb.SwitchOrganizationRequest{
		Data: &tenantpb.SwitchOrganizationRequestData{
			OrganizationId: in.Data.OrganizationId,
		},
	})
}
