package rolelogic

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

type AssignOrganizationUsersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewAssignOrganizationUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AssignOrganizationUsersLogic {
	return &AssignOrganizationUsersLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

// AssignOrganizationUsers binds organization users to a role.
func (l *AssignOrganizationUsersLogic) AssignOrganizationUsers(in *tenant.RoleAssignOrganizationUsersRequest) (resp *tenant.RoleDetailResponse, err error) {
	if in.Data == nil {
		return nil, errorx.ErrorEmptyRequestData
	}
	userType, err := logic.ParseUserType(in.Data.UserType)
	if err != nil {
		return nil, err
	}

	input, err := copierx.Trans[tenant.RoleAssignOrganizationUsersRequestData, dto.RoleAssignOrganizationUsersRequestDto](in.Data)
	if err != nil {
		return nil, err
	}
	input.Type = userType
	input.OrganizationId = in.Data.OrganizationId

	role, err := l.baseLogic.RoleLogic.AssignOrganizationUsers(l.ctx, *input)
	if err != nil {
		return nil, err
	}

	data, err := copierx.Trans[dto.RoleEntityDto, proto.RoleEntity](role, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}
	return &tenant.RoleDetailResponse{Data: data}, nil
}
