package userlogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/copierx"
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

// Login proxies to tenant, the sole identity authority. user_type pins admin.
func (l *LoginLogic) Login(in *admin.UserLoginRequest) (resp *admin.UserLoginResponse, err error) {
	loginResp, err := l.svcCtx.TenantClient.User.Login(l.ctx, &tenant.LoginRequest{
		Data: &tenant.LoginRequestData{
			UserType: string(constants.UserTypeAdmin),
			Email:    in.Data.Email,
			Password: in.Data.Password,
		},
	})
	if err != nil {
		return nil, err
	}
	data, err := copierx.Trans[proto.UserLoginResponseData, admin.UserLoginResponseData](loginResp.Data)
	if err != nil {
		return nil, err
	}
	return &admin.UserLoginResponse{Data: data}, nil
}
