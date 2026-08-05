package rolelogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/constants"
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
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *AssignOrganizationUsersLogic) AssignOrganizationUsers(in *admin.RoleAssignOrganizationUsersRequest) (resp *admin.RoleDetailResponse, err error) {

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx, constants.UserTypeAdmin)
	if err != nil {
		return nil, err
	}
	input, err := copierx.Trans[admin.RoleAssignOrganizationUsersRequestData, dto.RoleAssignOrganizationUsersRequestDto](in.Data)
	if err != nil {
		return nil, err
	}
	input.Type = constants.UserTypeAdmin
	input.OrganizationId = claim.OrganizationId

	role, err := l.baseLogic.RoleLogic.AssignOrganizationUsers(l.ctx, *input)
	if err != nil {
		return nil, err
	}

	data, err := copierx.Trans[dto.RoleEntityDto, proto.RoleEntity](role, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}
	return &admin.RoleDetailResponse{
		Data: data,
	}, nil
}
