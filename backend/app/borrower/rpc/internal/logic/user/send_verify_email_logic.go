package userlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/proto"
)

type SendVerifyEmailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendVerifyEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendVerifyEmailLogic {
	return &SendVerifyEmailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// SendVerifyEmail proxies to tenant; tenant rate-limits and mails the verify link.
func (l *SendVerifyEmailLogic) SendVerifyEmail(in *proto.SendVerifyEmailRequest) (resp *proto.SendVerifyEmailResponse, err error) {
	return l.svcCtx.TenantClient.User.SendVerifyEmail(l.ctx, &tenant.SendVerifyEmailRequest{
		Data: &tenant.SendVerifyEmailRequestData{
			UserType: string(constants.UserTypeBorrower),
			Email:    in.Data.Email,
		},
	})
}
