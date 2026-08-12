package rolelogic

import (
	"context"

	"zeebee-backend/app/tenant/rpc/internal/logic"
	"zeebee-backend/app/tenant/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/proto"
)

type SwitchStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewSwitchStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SwitchStatusLogic {
	return &SwitchStatusLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

// SwitchStatus enables / disables a role; direction B makes the change effective
// on the next AuthzCheck without any invalidation broadcast.
func (l *SwitchStatusLogic) SwitchStatus(in *tenant.RoleSwitchStatusRequest) (resp *proto.EmptyResponse, err error) {
	if in.Data == nil {
		return nil, errorx.ErrorEmptyRequestData
	}
	userType, err := logic.ParseUserType(in.Data.UserType)
	if err != nil {
		return nil, err
	}

	// before, _ := l.baseLogic.RoleLogic.Detail(l.ctx, in.Data.Id, in.Data.OrganizationId, userType)
	if err = l.baseLogic.RoleLogic.SwitchStatus(l.ctx, in.Data.Id, in.Data.OrganizationId, userType, constants.CommonStatus(in.Data.Status)); err != nil {
		return nil, err
	}
	// after, _ := l.baseLogic.RoleLogic.Detail(l.ctx, in.Data.Id, in.Data.OrganizationId, userType)
	// Temporarily disabled; re-enable when tenant audit logging is needed.
	// l.baseLogic.RecordAuditFromClaim(l.ctx, constants.OperationEntityTypeRole, in.Data.Id, constants.OperationActionRoleStatusSwitch, logic.AuditChange{Before: before, After: after})
	return &proto.EmptyResponse{}, nil
}
