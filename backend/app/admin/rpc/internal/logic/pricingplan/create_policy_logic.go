package pricingplanlogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/data/dto"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreatePolicyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewCreatePolicyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePolicyLogic {
	return &CreatePolicyLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *CreatePolicyLogic) CreatePolicy(in *admin.PricingPlanPolicyCreateRequest) (resp *admin.PricingPlanPolicyCreateResponse, err error) {

	dtoData, err := copierx.Trans[admin.PricingPlanPolicyCreateRequestData, dto.PricingPlanPolicyCreateDto](in.Data)
	if err != nil {
		return nil, err
	}

	result, err := l.baseLogic.PricingPlanPolicyLogic.Create(l.ctx, *dtoData)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	protoData, err := copierx.Trans[dto.PricingPlanPolicyEntityDto, admin.PricingPlanPolicyEntity](result, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}

	return &admin.PricingPlanPolicyCreateResponse{Data: protoData}, nil
}
