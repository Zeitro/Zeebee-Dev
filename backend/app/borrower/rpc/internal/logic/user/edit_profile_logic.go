package userlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/domains/zeebee/proto"
)

type EditProfileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEditProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EditProfileLogic {
	return &EditProfileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// EditProfile proxies to tenant. user_id and userType are resolved by tenant
// from the caller's JWT (forwarded by the borrower rpc server interceptor).
func (l *EditProfileLogic) EditProfile(in *borrower.UserProfileEditRequest) (resp *proto.UserinfoResponse, err error) {
	return l.svcCtx.TenantClient.User.EditProfile(l.ctx, &tenant.EditProfileRequest{
		Data: &tenant.EditProfileRequestData{
			FirstName: in.Data.FirstName,
			LastName:  in.Data.LastName,
			Phone:     in.Data.Phone,
		},
	})
}
