package organizationsubscriptionlogic

import (
	"context"

	"zeebee-backend/app/tenant/rpc/internal/logic"
	"zeebee-backend/app/tenant/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/proto"
)

type RenewLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewRenewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RenewLogic {
	return &RenewLogic{ctx: ctx, svcCtx: svcCtx, baseLogic: logic.GetBaseLogic(svcCtx)}
}

// Renew extends a subscription by one cycle.
func (l *RenewLogic) Renew(in *tenant.OrganizationSubscriptionRenewRequest) (resp *tenant.OrganizationSubscriptionDetailResponse, err error) {
	if in.Data == nil {
		return nil, errorx.ErrorEmptyRequestData
	}
	if _, err = logic.ParseUserType(in.Data.UserType); err != nil {
		return
	}

	result, err := l.baseLogic.OrganizationSubscriptionLogic.Renew(l.ctx, in.Data.Id)
	if err != nil {
		return
	}
	l.baseLogic.RecordAuditFromClaim(l.ctx, constants.OperationEntityTypeOrganizationSubscription, result.Id, constants.OperationActionSubscriptionRenew)

	data, err := copierx.Trans[dto.OrganizationSubscriptionEntity, proto.OrganizationSubscriptionEntity](result, copierx.TimeToRFC3339Str())
	if err != nil {
		return
	}
	return &tenant.OrganizationSubscriptionDetailResponse{Data: data}, nil
}
