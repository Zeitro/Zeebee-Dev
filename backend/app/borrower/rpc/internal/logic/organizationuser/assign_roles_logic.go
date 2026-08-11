package organizationuserlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/proto"
)

type AssignRolesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewAssignRolesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AssignRolesLogic {
	return &AssignRolesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *AssignRolesLogic) AssignRoles(in *borrower.OrganizationUserAssignRolesRequest) (resp *proto.EmptyResponse, err error) {

	claim, err := l.baseLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	return l.svcCtx.TenantClient.OrganizationUser.AssignRoles(l.ctx, &tenant.OrganizationUserAssignRolesRequest{Data: &tenant.OrganizationUserAssignRolesRequestData{
		UserType:       string(constants.UserTypeBorrower),
		OrganizationId: claim.OrganizationId,
		Id:             in.Data.Id,
		RoleIds:        in.Data.RoleIds,
	}})
}
