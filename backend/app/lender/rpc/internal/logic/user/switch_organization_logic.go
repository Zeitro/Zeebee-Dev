package userlogic

import (
	"context"

	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/app/lender/rpc/internal/svc"
	tenantpb "zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/domains/zeebee/proto"
)

type SwitchOrganizationLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewSwitchOrganizationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SwitchOrganizationLogic {
	return &SwitchOrganizationLogic{ctx: ctx, svcCtx: svcCtx, baseLogic: logic.GetBaseLogic(svcCtx)}
}

// SwitchOrganization proxies to tenant. user_id is resolved from the caller's JWT
// (lender proto carries only organization_id); the lender rpc server interceptor
// forwards the token so tenant can read userType from the claim.
func (l *SwitchOrganizationLogic) SwitchOrganization(in *proto.SwitchOrganizationRequest) (resp *proto.UserLoginResponse, err error) {
	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	return l.svcCtx.TenantClient.User.SwitchOrganization(l.ctx, &tenantpb.SwitchOrganizationRequest{
		Data: &tenantpb.SwitchOrganizationRequestData{
			UserId:         claim.UserId,
			OrganizationId: in.Data.OrganizationId,
		},
	})
}
