package organizationuserlogic

import (
	"context"

	"zeebee-backend/app/tenant/rpc/internal/logic"
	"zeebee-backend/app/tenant/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/proto"
)

type AcceptInvitationLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewAcceptInvitationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AcceptInvitationLogic {
	return &AcceptInvitationLogic{ctx: ctx, svcCtx: svcCtx, baseLogic: logic.GetBaseLogic(svcCtx)}
}

// AcceptInvitation completes an invitation and returns the activated user.
func (l *AcceptInvitationLogic) AcceptInvitation(in *tenant.OrganizationUserAcceptInvitationRequest) (resp *tenant.OrganizationUserAcceptInvitationResponse, err error) {
	if in.Data == nil {
		return nil, errorx.ErrorEmptyRequestData
	}
	userType, err := logic.ParseUserType(in.Data.UserType)
	if err != nil {
		return
	}

	input := dto.OrganizationUserAcceptInvitationRequestDto{
		Code:           in.Data.Code,
		UserType:       userType,
		Email:          in.Data.Email,
		OrganizationId: in.Data.OrganizationId,
		FirstName:      in.Data.FirstName,
		LastName:       in.Data.LastName,
	}
	user, err := l.baseLogic.OrganizationUserLogic.AcceptInvitation(l.ctx, input)
	if err != nil {
		return
	}

	data, err := copierx.Trans[dto.UserEntityDto, proto.UserEntity](&user, copierx.TimeToRFC3339Str())
	if err != nil {
		return
	}
	return &tenant.OrganizationUserAcceptInvitationResponse{Data: data}, nil
}
