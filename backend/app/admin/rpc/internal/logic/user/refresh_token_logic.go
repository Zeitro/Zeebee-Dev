package userlogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/proto"
)

type RefreshTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRefreshTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshTokenLogic {
	return &RefreshTokenLogic{ctx: ctx, svcCtx: svcCtx}
}

// RefreshToken proxies to tenant; tenant resolves userType from the refresh
// token. All backends share the same full-TokenSet refresh contract (ODD-6).
func (l *RefreshTokenLogic) RefreshToken(in *admin.RefreshTokenRequest) (resp *admin.RefreshTokenResponse, err error) {
	loginResp, err := l.svcCtx.TenantClient.User.RefreshToken(l.ctx, &tenant.RefreshTokenRequest{
		Data: &tenant.RefreshTokenRequestData{
			RefreshToken: in.Data.RefreshToken,
		},
	})
	if err != nil {
		return nil, err
	}
	data, err := copierx.Trans[proto.UserLoginResponseData, admin.RefreshTokenResponseData](loginResp.Data)
	if err != nil {
		return nil, err
	}
	return &admin.RefreshTokenResponse{Data: data}, nil
}
