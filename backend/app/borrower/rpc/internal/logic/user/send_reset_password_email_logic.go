package userlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/proto"
)

type SendResetPasswordEmailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendResetPasswordEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendResetPasswordEmailLogic {
	return &SendResetPasswordEmailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// SendResetPasswordEmail proxies to tenant; tenant rate-limits and mails the link.
func (l *SendResetPasswordEmailLogic) SendResetPasswordEmail(in *proto.SendResetPasswordEmailRequest) (resp *proto.SendResetPasswordEmailResponse, err error) {
	return l.svcCtx.TenantClient.User.SendResetPasswordEmail(l.ctx, &tenant.SendResetPasswordEmailRequest{
		Data: &tenant.SendResetPasswordEmailRequestData{
			UserType: string(constants.UserTypeBorrower),
			Email:    in.Data.Email,
		},
	})
}
