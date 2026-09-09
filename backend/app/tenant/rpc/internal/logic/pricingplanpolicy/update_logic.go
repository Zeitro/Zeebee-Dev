package pricingplanpolicylogic

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

type UpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *UpdateLogic) Update(in *tenant.PricingPlanPolicyUpdateRequest) (resp *tenant.PricingPlanPolicyDetailResponse, err error) {
	if in.Data == nil {
		return nil, errorx.ErrorEmptyRequestData
	}
	if _, err = logic.ParseUserType(in.Data.UserType); err != nil {
		return nil, err
	}

	input, err := copierx.Trans[tenant.PricingPlanPolicyUpdateRequestData, dto.PricingPlanPolicyUpdateRequestDto](in.Data)
	if err != nil {
		return nil, err
	}
	// before, _ := l.baseLogic.PricingPlanPolicyLogic.Detail(l.ctx, in.Data.Id)
	result, err := l.baseLogic.PricingPlanPolicyLogic.Update(l.ctx, *input)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	// Temporarily disabled; re-enable when tenant audit logging is needed.
	// l.baseLogic.RecordAuditFromClaim(l.ctx, constants.OperationEntityTypePricingPlanPolicy, result.Id, constants.OperationActionPlanPolicyUpdate, logic.AuditChange{Before: before, After: result})
	data, err := copierx.Trans[dto.PricingPlanPolicyEntityDto, proto.PricingPlanPolicyEntity](result, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}
	return &tenant.PricingPlanPolicyDetailResponse{Data: data}, nil
}
