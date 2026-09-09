package organizationuserlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
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

func (l *CheckInvitationLogic) CheckInvitation(in *borrower.OrganizationUserCheckInvitationRequest) (resp *borrower.OrganizationUserCheckInvitationResponse, err error) {

	input, err := copierx.Trans[borrower.OrganizationUserCheckInvitationRequestData, dto.OrganizationUserCheckInvitationCodeRequestDto](in.Data)
	if err != nil {
		return
	}
	input.UserType = constants.UserTypeBorrower
	isExpired, err := l.baseLogic.OrganizationUserLogic.CheckInvitationCode(l.ctx, *input)
	if err != nil {
		return nil, err
	}

	return &borrower.OrganizationUserCheckInvitationResponse{
		Data: &borrower.OrganizationUserCheckInvitationResponseData{
			IsExpired: isExpired,
		},
	}, nil
}
