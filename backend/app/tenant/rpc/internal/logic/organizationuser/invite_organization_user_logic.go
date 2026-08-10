package organizationuserlogic

import (
	"context"

	"zeebee-backend/app/tenant/rpc/internal/logic"
	"zeebee-backend/app/tenant/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/proto"
)

type InviteOrganizationUserLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewInviteOrganizationUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InviteOrganizationUserLogic {
	return &InviteOrganizationUserLogic{ctx: ctx, svcCtx: svcCtx, baseLogic: logic.GetBaseLogic(svcCtx)}
}

// InviteOrganizationUser invites a new member by email and assigns roles.
func (l *InviteOrganizationUserLogic) InviteOrganizationUser(in *tenant.OrganizationUserInviteRequest) (resp *tenant.OrganizationUserDetailResponse, err error) {
	if in.Data == nil {
		return nil, errorx.ErrorEmptyRequestData
	}
	userType, err := logic.ParseUserType(in.Data.UserType)
	if err != nil {
		return
	}

	input := dto.OrganizationUserCreateRequestDto{
		Email:          in.Data.Email,
		Phone:          in.Data.Phone,
		Remark:         in.Data.Remark,
		RoleIds:        in.Data.RoleIds,
		UserType:       userType,
		OrganizationId: in.Data.OrganizationId,
		EmailConfig:    toInviteEmailConfig(in.Data.EmailConfig),
		Operator:       in.Data.Operator,
	}
	result, err := l.baseLogic.OrganizationUserLogic.InviteOrganizationUser(l.ctx, input)
	if err != nil {
		return
	}
	l.baseLogic.RecordAuditFromClaim(l.ctx, constants.OperationEntityTypeOrganizationUser, result.Id, constants.OperationActionMemberInvite)

	data, err := copierx.Trans[dto.OrganizationUserEntityDto, proto.OrganizationUserEntity](result, copierx.TimeToRFC3339Str())
	if err != nil {
		return
	}
	return &tenant.OrganizationUserDetailResponse{Data: data}, nil
}
