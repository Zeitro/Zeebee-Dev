package rolelogic

import (
	"context"

	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/constants"
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

func (l *AssignOrganizationUsersLogic) AssignOrganizationUsers(in *lender.RoleAssignOrganizationUsersRequest) (resp *lender.RoleDetailResponse, err error) {

	claim, err := l.baseLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	input, err := copierx.Trans[lender.RoleAssignOrganizationUsersRequestData, tenant.RoleAssignOrganizationUsersRequestData](in.Data)
	if err != nil {
		return nil, err
	}
	input.UserType = string(constants.UserTypeLender)
	input.OrganizationId = claim.OrganizationId

	role, err := l.svcCtx.TenantClient.Role.AssignOrganizationUsers(l.ctx, &tenant.RoleAssignOrganizationUsersRequest{Data: input})
	if err != nil {
		return nil, err
	}

	return &lender.RoleDetailResponse{
		Data: role.Data,
	}, nil
}
