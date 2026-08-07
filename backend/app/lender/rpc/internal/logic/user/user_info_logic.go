package userlogic

import (
	"context"

	"zeebee-backend/app/lender/rpc/internal/svc"
	tenantpb "zeebee-backend/app/tenant/rpc/pb/tenant"
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
// the caller's JWT (forwarded by the lender rpc server interceptor).
func (l *UserInfoLogic) UserInfo(in *proto.EmptyRequest) (resp *proto.UserinfoResponse, err error) {
	return l.svcCtx.TenantClient.User.UserInfo(l.ctx, &tenantpb.UserInfoRequest{
		Data: &tenantpb.UserInfoRequestData{},
	})
}
