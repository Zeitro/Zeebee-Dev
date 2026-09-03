package userlogic

import (
	"context"

	"zeebee-backend/app/lender/rpc/internal/svc"
	tenantpb "zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/proto"
)

type VerifyEmailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVerifyEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerifyEmailLogic {
	return &VerifyEmailLogic{ctx: ctx, svcCtx: svcCtx}
}

// VerifyEmail proxies to tenant; tenant activates the user and issues the token set.
func (l *VerifyEmailLogic) VerifyEmail(in *proto.VerifyEmailRequest) (resp *proto.UserLoginResponse, err error) {
	return l.svcCtx.TenantClient.User.VerifyEmail(l.ctx, &tenantpb.VerifyEmailRequest{
		Data: &tenantpb.VerifyEmailRequestData{
			UserType: string(constants.UserTypeLender),
			Email:    in.Data.Email,
			Code:     in.Data.Code,
		},
	})
}
