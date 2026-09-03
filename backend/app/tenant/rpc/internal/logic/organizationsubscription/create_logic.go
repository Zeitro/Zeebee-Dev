package organizationsubscriptionlogic

import (
	"context"

	"zeebee-backend/app/tenant/rpc/internal/logic"
	"zeebee-backend/app/tenant/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/proto"
)

type CreateLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{ctx: ctx, svcCtx: svcCtx, baseLogic: logic.GetBaseLogic(svcCtx)}
}

// Create creates an organization subscription.
func (l *CreateLogic) Create(in *tenant.OrganizationSubscriptionCreateRequest) (resp *tenant.OrganizationSubscriptionDetailResponse, err error) {
	if in.Data == nil {
		return nil, errorx.ErrorEmptyRequestData
	}
	userType, err := logic.ParseUserType(in.Data.UserType)
	if err != nil {
		return
	}

	input := dto.OrganizationSubscriptionCreateDto{
		OrganizationId:      in.Data.OrganizationId,
		PricingPlanId:       in.Data.PricingPlanId,
		PricingPlanPolicyId: in.Data.PricingPlanPolicyId,
		Type:                userType,
	}
	result, err := l.baseLogic.OrganizationSubscriptionLogic.Create(l.ctx, input)
	if err != nil {
		return
	}
	// Temporarily disabled; re-enable when tenant audit logging is needed.
	// l.baseLogic.RecordAuditFromClaim(l.ctx, constants.OperationEntityTypeOrganizationSubscription, result.Id, constants.OperationActionSubscriptionCreate, logic.AuditChange{After: result})

	data, err := copierx.Trans[dto.OrganizationSubscriptionEntity, proto.OrganizationSubscriptionEntity](result, copierx.TimeToRFC3339Str())
	if err != nil {
		return
	}
	return &tenant.OrganizationSubscriptionDetailResponse{Data: data}, nil
}
