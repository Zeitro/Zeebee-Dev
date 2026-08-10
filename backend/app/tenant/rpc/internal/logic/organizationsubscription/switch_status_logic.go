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

type SwitchStatusLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewSwitchStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SwitchStatusLogic {
	return &SwitchStatusLogic{ctx: ctx, svcCtx: svcCtx, baseLogic: logic.GetBaseLogic(svcCtx)}
}

// SwitchStatus updates a subscription's status.
func (l *SwitchStatusLogic) SwitchStatus(in *tenant.OrganizationSubscriptionSwitchStatusRequest) (resp *proto.EmptyResponse, err error) {
	if in.Data == nil {
		return nil, errorx.ErrorEmptyRequestData
	}
	if _, err = logic.ParseUserType(in.Data.UserType); err != nil {
		return
	}

	if err = l.baseLogic.OrganizationSubscriptionLogic.SwitchStatus(l.ctx, in.Data.Id, constants.OrganizationSubscriptionStatus(in.Data.Status)); err != nil {
		return
	}
	l.baseLogic.RecordAuditFromClaim(l.ctx, constants.OperationEntityTypeOrganizationSubscription, in.Data.Id, constants.OperationActionSubscriptionStatusSwitch)
	return &proto.EmptyResponse{}, nil
}
