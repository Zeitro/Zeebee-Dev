package organizationuserlogic

import (
	"context"

	"zeebee-backend/app/tenant/rpc/internal/logic"
	"zeebee-backend/app/tenant/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/data/dto"
)

type CheckInvitationCodeLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewCheckInvitationCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckInvitationCodeLogic {
	return &CheckInvitationCodeLogic{ctx: ctx, svcCtx: svcCtx, baseLogic: logic.GetBaseLogic(svcCtx)}
}

// CheckInvitationCode reports whether an invitation code is expired.
func (l *CheckInvitationCodeLogic) CheckInvitationCode(in *tenant.OrganizationUserCheckInvitationCodeRequest) (resp *tenant.OrganizationUserCheckInvitationCodeResponse, err error) {
	if in.Data == nil {
		return nil, errorx.ErrorEmptyRequestData
	}
	userType, err := logic.ParseUserType(in.Data.UserType)
	if err != nil {
		return
	}

	input := dto.OrganizationUserCheckInvitationCodeRequestDto{
		Code:           in.Data.Code,
		UserType:       userType,
		Email:          in.Data.Email,
		OrganizationId: in.Data.OrganizationId,
	}
	isExpired, err := l.baseLogic.OrganizationUserLogic.CheckInvitationCode(l.ctx, input)
	if err != nil {
		return
	}
	return &tenant.OrganizationUserCheckInvitationCodeResponse{IsExpired: isExpired}, nil
}
