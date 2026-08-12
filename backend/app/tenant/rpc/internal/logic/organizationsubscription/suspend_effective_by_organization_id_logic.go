package organizationsubscriptionlogic

import (
	"context"

	"zeebee-backend/app/tenant/rpc/internal/logic"
	"zeebee-backend/app/tenant/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/proto"
)

type SuspendEffectiveByOrganizationIdLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewSuspendEffectiveByOrganizationIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SuspendEffectiveByOrganizationIdLogic {
	return &SuspendEffectiveByOrganizationIdLogic{ctx: ctx, svcCtx: svcCtx, baseLogic: logic.GetBaseLogic(svcCtx)}
}

// SuspendEffectiveByOrganizationId suspends a subscription when its organization is
// disabled; also exposed for admin trigger.
func (l *SuspendEffectiveByOrganizationIdLogic) SuspendEffectiveByOrganizationId(in *tenant.OrganizationSubscriptionSuspendEffectiveByOrganizationIdRequest) (resp *proto.EmptyResponse, err error) {
	if in.Data == nil {
		return nil, errorx.ErrorEmptyRequestData
	}
	userType, err := logic.ParseUserType(in.Data.UserType)
	if err != nil {
		return
	}

	if err = l.baseLogic.OrganizationSubscriptionLogic.SuspendEffectiveByOrganizationId(l.ctx, in.Data.OrganizationId, userType); err != nil {
		return
	}
	// Temporarily disabled; re-enable when tenant audit logging is needed.
	// l.baseLogic.RecordAuditFromClaim(l.ctx, constants.OperationEntityTypeOrganizationSubscription, in.Data.OrganizationId, constants.OperationActionSubscriptionSuspend)
	return &proto.EmptyResponse{}, nil
}
