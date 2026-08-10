package organizationuserlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/proto"
)

type AcceptInvitationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewAcceptInvitationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AcceptInvitationLogic {
	return &AcceptInvitationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *AcceptInvitationLogic) AcceptInvitation(in *borrower.OrganizationUserAcceptInvitationRequest) (resp *proto.UserLoginResponse, err error) {

	input, err := copierx.Trans[borrower.OrganizationUserAcceptInvitationRequestData, tenant.OrganizationUserAcceptInvitationRequestData](in.Data)
	if err != nil {
		return nil, err
	}

	input.UserType = string(constants.UserTypeBorrower)
	user, err := l.svcCtx.TenantClient.OrganizationUser.AcceptInvitation(l.ctx, &tenant.OrganizationUserAcceptInvitationRequest{Data: input})
	if err != nil {
		return nil, err
	}
	token, err := l.baseLogic.UserLogic.GenerateJwtTokenById(l.ctx, user.Data.Id, user.Data.DefaultOrganizationId, constants.UserTypeBorrower)
	if err != nil {
		return nil, err
	}
	data, err := copierx.Trans[dto.UserJwtTokenDto, proto.UserLoginResponseData](token)
	if err != nil {
		return nil, err
	}

	return &proto.UserLoginResponse{Data: data}, nil
}
