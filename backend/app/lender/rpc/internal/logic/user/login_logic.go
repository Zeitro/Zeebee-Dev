package userlogic

import (
	"context"

	"zeebee-backend/app/lender/rpc/internal/svc"
	tenantpb "zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/proto"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{ctx: ctx, svcCtx: svcCtx}
}

// Login proxies to tenant, the sole identity authority. user_type pins lender.
func (l *LoginLogic) Login(in *proto.UserLoginRequest) (resp *proto.UserLoginResponse, err error) {
	return l.svcCtx.TenantClient.User.Login(l.ctx, &tenantpb.LoginRequest{
		Data: &tenantpb.LoginRequestData{
			UserType: string(constants.UserTypeLender),
			Email:    in.Data.Email,
			Password: in.Data.Password,
		},
	})
}
