package userlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/proto"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Login proxies to tenant, the sole identity authority. user_type pins borrower.
func (l *LoginLogic) Login(in *proto.UserLoginRequest) (resp *proto.UserLoginResponse, err error) {
	return l.svcCtx.TenantClient.User.Login(l.ctx, &tenant.LoginRequest{
		Data: &tenant.LoginRequestData{
			UserType: string(constants.UserTypeBorrower),
			Email:    in.Data.Email,
			Password: in.Data.Password,
		},
	})
}
