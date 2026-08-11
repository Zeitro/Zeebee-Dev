package organizationuserlogic

import (
	"context"

	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/proto"
)

type AcceptInvitationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAcceptInvitationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AcceptInvitationLogic {
	return &AcceptInvitationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AcceptInvitationLogic) AcceptInvitation(in *lender.OrganizationUserAcceptInvitationRequest) (resp *proto.UserLoginResponse, err error) {

	input, err := copierx.Trans[lender.OrganizationUserAcceptInvitationRequestData, tenant.OrganizationUserAcceptInvitationRequestData](in.Data)
	if err != nil {
		return nil, err
	}
	input.UserType = string(constants.UserTypeLender)
	user, err := l.svcCtx.TenantClient.OrganizationUser.AcceptInvitation(l.ctx, &tenant.OrganizationUserAcceptInvitationRequest{Data: input})
	if err != nil {
		return nil, err
	}
	return user.Token, nil
}
