package organizationuserlogic

import (
	"context"

	"zeebee-backend/app/tenant/rpc/internal/logic"
	"zeebee-backend/app/tenant/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/proto"
)

type AssignRolesLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewAssignRolesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AssignRolesLogic {
	return &AssignRolesLogic{ctx: ctx, svcCtx: svcCtx, baseLogic: logic.GetBaseLogic(svcCtx)}
}

// AssignRoles binds roles to an organization member.
func (l *AssignRolesLogic) AssignRoles(in *tenant.OrganizationUserAssignRolesRequest) (resp *proto.EmptyResponse, err error) {
	if in.Data == nil {
		return nil, errorx.ErrorEmptyRequestData
	}
	userType, err := logic.ParseUserType(in.Data.UserType)
	if err != nil {
		return
	}

	if err = l.baseLogic.OrganizationUserLogic.AssignRoles(l.ctx, in.Data.Id, in.Data.OrganizationId, in.Data.RoleIds, userType); err != nil {
		return
	}
	// Temporarily disabled; re-enable when tenant audit logging is needed.
	// l.baseLogic.RecordAuditFromClaim(l.ctx, constants.OperationEntityTypeOrganizationUser, in.Data.Id, constants.OperationActionMemberAssignRoles)
	return &proto.EmptyResponse{}, nil
}
