package organizationuserlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/domains/zeebee/proto"
)

type SendInvitationEmailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewSendInvitationEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendInvitationEmailLogic {
	return &SendInvitationEmailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *SendInvitationEmailLogic) SendInvitationEmail(in *proto.PathIdRequest) (resp *proto.EmptyResponse, err error) {

	claim, err := l.baseLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	return l.svcCtx.TenantClient.OrganizationUser.SendInviteEmail(l.ctx, &tenant.OrganizationUserSendInviteEmailRequest{Data: &tenant.OrganizationUserSendInviteEmailRequestData{
		UserType:       string(claim.UserType),
		OrganizationId: claim.OrganizationId,
		Id:             in.Data.Id,
		Operator:       claim.Username,
	}})
}
