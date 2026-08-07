package userlogic

import (
	"context"

	"zeebee-backend/app/lender/rpc/internal/svc"
	tenantpb "zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/domains/zeebee/proto"
)

type RefreshTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRefreshTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshTokenLogic {
	return &RefreshTokenLogic{ctx: ctx, svcCtx: svcCtx}
}

// RefreshToken proxies to tenant; tenant resolves userType from the refresh token.
func (l *RefreshTokenLogic) RefreshToken(in *proto.RefreshTokenRequest) (resp *proto.UserLoginResponse, err error) {
	return l.svcCtx.TenantClient.User.RefreshToken(l.ctx, &tenantpb.RefreshTokenRequest{
		Data: &tenantpb.RefreshTokenRequestData{
			RefreshToken: in.Data.RefreshToken,
		},
	})
}
