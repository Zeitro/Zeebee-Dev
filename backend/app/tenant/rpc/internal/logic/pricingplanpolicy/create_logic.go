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

type CreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *CreateLogic) Create(in *tenant.PricingPlanPolicyCreateRequest) (resp *tenant.PricingPlanPolicyDetailResponse, err error) {
	if in.Data == nil {
		return nil, errorx.ErrorEmptyRequestData
	}
	if _, err = logic.ParseUserType(in.Data.UserType); err != nil {
		return nil, err
	}

	input, err := copierx.Trans[tenant.PricingPlanPolicyCreateRequestData, dto.PricingPlanPolicyCreateRequestDto](in.Data)
	if err != nil {
		return nil, err
	}
	result, err := l.baseLogic.PricingPlanPolicyLogic.Create(l.ctx, *input)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	// Temporarily disabled; re-enable when tenant audit logging is needed.
	// l.baseLogic.RecordAuditFromClaim(l.ctx, constants.OperationEntityTypePricingPlanPolicy, result.Id, constants.OperationActionPlanPolicyCreate, logic.AuditChange{After: result})
	data, err := copierx.Trans[dto.PricingPlanPolicyEntityDto, proto.PricingPlanPolicyEntity](result, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}
	return &tenant.PricingPlanPolicyDetailResponse{Data: data}, nil
}
