package organizationsubscriptionlogic

import (
	"context"

	"zeebee-backend/app/tenant/rpc/internal/logic"
	"zeebee-backend/app/tenant/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/proto"
)

type SubscribeDefaultPlanLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewSubscribeDefaultPlanLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SubscribeDefaultPlanLogic {
	return &SubscribeDefaultPlanLogic{ctx: ctx, svcCtx: svcCtx, baseLogic: logic.GetBaseLogic(svcCtx)}
}

// SubscribeDefaultPlan onboards an organization onto its default pricing plan.
func (l *SubscribeDefaultPlanLogic) SubscribeDefaultPlan(in *tenant.OrganizationSubscriptionSubscribeDefaultPlanRequest) (resp *proto.EmptyResponse, err error) {
	if in.Data == nil {
		return nil, errorx.ErrorEmptyRequestData
	}
	userType, err := logic.ParseUserType(in.Data.UserType)
	if err != nil {
		return
	}

	if err = l.baseLogic.OrganizationSubscriptionLogic.SubscribeDefaultPlan(l.ctx, in.Data.OrganizationId, userType); err != nil {
		return
	}
	l.baseLogic.RecordAuditFromClaim(l.ctx, constants.OperationEntityTypeOrganizationSubscription, in.Data.OrganizationId, constants.OperationActionSubscriptionDefaultPlan)
	return &proto.EmptyResponse{}, nil
}
