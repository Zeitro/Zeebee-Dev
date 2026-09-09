package userlogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/proto"
)

type EditProfileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEditProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EditProfileLogic {
	return &EditProfileLogic{ctx: ctx, svcCtx: svcCtx}
}

// EditProfile proxies to tenant. user_id and userType are resolved by tenant
// from the caller's JWT (forwarded by the admin rpc server interceptor).
func (l *EditProfileLogic) EditProfile(in *admin.UserProfileEditRequest) (resp *admin.UserinfoResponse, err error) {
	userInfoResp, err := l.svcCtx.TenantClient.User.EditProfile(l.ctx, &tenant.EditProfileRequest{
		Data: &tenant.EditProfileRequestData{
			FirstName: in.Data.FirstName,
			LastName:  in.Data.LastName,
			Phone:     in.Data.Phone,
		},
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
