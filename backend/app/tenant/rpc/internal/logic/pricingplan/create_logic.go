package pricingplanlogic

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

func (l *CreateLogic) Create(in *tenant.PricingPlanCreateRequest) (resp *tenant.PricingPlanDetailResponse, err error) {
	if in.Data == nil {
		return nil, errorx.ErrorEmptyRequestData
	}
	if _, err = logic.ParseUserType(in.Data.UserType); err != nil {
		return nil, err
	}

	input, err := copierx.Trans[tenant.PricingPlanCreateRequestData, dto.PricingPlanCreateRequestDto](in.Data)
	if err != nil {
		return nil, err
	}
	result, err := l.baseLogic.PricingPlanLogic.Create(l.ctx, *input)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	l.baseLogic.RecordAuditFromClaim(l.ctx, constants.OperationEntityTypePricingPlan, result.Id, constants.OperationActionPlanCreate)
	data, err := copierx.Trans[dto.PricingPlanEntityDto, proto.PricingPlanEntity](result, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}
	return &tenant.PricingPlanDetailResponse{Data: data}, nil
}
