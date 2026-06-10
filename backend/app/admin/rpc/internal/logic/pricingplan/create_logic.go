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

type CreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *CreateLogic) Create(in *admin.PricingPlanCreateRequest) (resp *admin.PricingPlanCreateResponse, err error) {

	dtoData, err := copierx.Trans[admin.PricingPlanCreateRequestData, dto.PricingPlanCreateDto](in.Data)
	if err != nil {
		return nil, err
	}

	result, err := l.baseLogic.PricingPlanLogic.Create(l.ctx, *dtoData)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	protoData, err := copierx.Trans[dto.PricingPlanDetailEntityDto, admin.PricingPlanEntity](result, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}

	return &admin.PricingPlanCreateResponse{Data: protoData}, nil
}
