package organizationuserlogic

import (
	"context"

	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/constants"
)

type CheckInvitationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewCheckInvitationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckInvitationLogic {
	return &CheckInvitationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *CheckInvitationLogic) CheckInvitation(in *lender.OrganizationUserCheckInvitationRequest) (resp *lender.OrganizationUserCheckInvitationResponse, err error) {

	input, err := copierx.Trans[lender.OrganizationUserCheckInvitationRequestData, tenant.OrganizationUserCheckInvitationCodeRequestData](in.Data)
	if err != nil {
		return
	}
	input.UserType = string(constants.UserTypeLender)
	result, err := l.svcCtx.TenantClient.OrganizationUser.CheckInvitationCode(l.ctx, &tenant.OrganizationUserCheckInvitationCodeRequest{Data: input})
	if err != nil {
		return nil, err
	}

	return &lender.OrganizationUserCheckInvitationResponse{
		Data: &lender.OrganizationUserCheckInvitationResponseData{
			IsExpired: result.IsExpired,
		},
	}, nil
}
