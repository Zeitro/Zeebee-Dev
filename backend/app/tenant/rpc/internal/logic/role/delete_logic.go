package rolelogic

import (
	"context"

	"zeebee-backend/app/tenant/rpc/internal/logic"
	"zeebee-backend/app/tenant/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/proto"
)

type DeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogic {
	return &DeleteLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

// Delete removes a role within the organization.
func (l *DeleteLogic) Delete(in *tenant.RoleDeleteRequest) (resp *proto.EmptyResponse, err error) {
	if in.Data == nil {
		return nil, errorx.ErrorEmptyRequestData
	}
	userType, err := logic.ParseUserType(in.Data.UserType)
	if err != nil {
		return nil, err
	}

	// before, _ := l.baseLogic.RoleLogic.Detail(l.ctx, in.Data.Id, in.Data.OrganizationId, userType)
	if err = l.baseLogic.RoleLogic.Delete(l.ctx, in.Data.Id, in.Data.OrganizationId, userType); err != nil {
		return nil, err
	}
	// Temporarily disabled; re-enable when tenant audit logging is needed.
	// l.baseLogic.RecordAuditFromClaim(l.ctx, constants.OperationEntityTypeRole, in.Data.Id, constants.OperationActionRoleDelete, logic.AuditChange{Before: before})
	return &proto.EmptyResponse{}, nil
}
