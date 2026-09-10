package userlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/proto"
)

type ResetPasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewResetPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetPasswordLogic {
	return &ResetPasswordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ResetPassword proxies to tenant (ResetPasswordByCode); tenant verifies the
// one-time code and re-issues the token set.
func (l *ResetPasswordLogic) ResetPassword(in *proto.ResetPasswordRequest) (resp *proto.ResetPasswordResponse, err error) {
	loginResp, err := l.svcCtx.TenantClient.User.ResetPasswordByCode(l.ctx, &tenant.ResetPasswordByCodeRequest{
		Data: &tenant.ResetPasswordByCodeRequestData{
			UserType:    string(constants.UserTypeBorrower),
			Email:       in.Data.Email,
			Code:        in.Data.Code,
			NewPassword: in.Data.Password,
		},
	})
	if err != nil {
		return nil, err
	}
	return &proto.ResetPasswordResponse{
		Data: loginResp.Data,
	}, nil
}
