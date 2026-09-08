package userlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/domains/zeebee/proto"
)

type ChangePasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChangePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangePasswordLogic {
	return &ChangePasswordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ChangePassword proxies to tenant. user_id and userType are resolved by tenant
// from the caller's JWT (forwarded by the borrower rpc server interceptor).
func (l *ChangePasswordLogic) ChangePassword(in *proto.ChangePasswordRequest) (resp *proto.EmptyResponse, err error) {
	return l.svcCtx.TenantClient.User.ChangePassword(l.ctx, &tenant.ChangePasswordRequest{
		Data: &tenant.ChangePasswordRequestData{
			OldPassword: in.Data.OldPassword,
			NewPassword: in.Data.NewPassword,
		},
	})
}
