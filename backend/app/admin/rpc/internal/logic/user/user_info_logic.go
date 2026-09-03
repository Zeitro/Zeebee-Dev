package userlogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/proto"
)

type UserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{ctx: ctx, svcCtx: svcCtx}
}

// UserInfo proxies to tenant. user_id and userType are resolved by tenant from
// the caller's JWT (forwarded by the admin rpc server interceptor).
func (l *UserInfoLogic) UserInfo(in *proto.EmptyRequest) (resp *admin.UserinfoResponse, err error) {
	userInfoResp, err := l.svcCtx.TenantClient.User.UserInfo(l.ctx, &tenant.UserInfoRequest{
		Data: &tenant.UserInfoRequestData{},
	})
	if err != nil {
		return nil, err
	}
	data, err := copierx.Trans[proto.UserinfoResponseData, admin.UserinfoResponseData](userInfoResp.Data)
	if err != nil {
		return nil, err
	}
	return &admin.UserinfoResponse{Data: data}, nil
}
