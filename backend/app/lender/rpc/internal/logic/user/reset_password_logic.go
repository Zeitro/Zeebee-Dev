package userlogic

import (
	"context"

	"zeebee-backend/app/lender/rpc/internal/svc"
	tenantpb "zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/proto"
)

type ResetPasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewResetPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetPasswordLogic {
	return &ResetPasswordLogic{ctx: ctx, svcCtx: svcCtx}
}

// ResetPassword proxies to tenant (ResetPasswordByCode); tenant verifies the
// one-time code and re-issues the token set.
func (l *ResetPasswordLogic) ResetPassword(in *proto.ResetPasswordRequest) (resp *proto.UserLoginResponse, err error) {
	return l.svcCtx.TenantClient.User.ResetPasswordByCode(l.ctx, &tenantpb.ResetPasswordByCodeRequest{
		Data: &tenantpb.ResetPasswordByCodeRequestData{
			UserType:    string(constants.UserTypeLender),
			Email:       in.Data.Email,
			Code:        in.Data.Code,
			NewPassword: in.Data.Password,
		},
	})
}
