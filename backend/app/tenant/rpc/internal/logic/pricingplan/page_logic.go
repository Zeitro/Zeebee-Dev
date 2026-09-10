package pricingplanlogic

import (
	"context"

	"zeebee-backend/app/tenant/rpc/internal/logic"
	"zeebee-backend/app/tenant/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/data/dto"
)

type PageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PageLogic {
	return &PageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *PageLogic) Page(in *tenant.PricingPlanPageRequest) (resp *tenant.PricingPlanPageResponse, err error) {
	if in.Data == nil {
		return nil, errorx.ErrorEmptyRequestData
	}
	userType, err := logic.ParseUserType(in.Data.UserType)
	if err != nil {
		return nil, err
	}

	pageReq, err := copierx.Trans[tenant.PricingPlanPageRequestData, dto.PricingPlanPageRequestDto](in.Data)
	if err != nil {
		return nil, err
	}
	pageReq.Filter.Type = string(userType)

	page, err := l.baseLogic.PricingPlanLogic.Page(l.ctx, *pageReq)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	data, err := copierx.Trans[gormx.PageResponse[dto.PricingPlanPageEntityDto], tenant.PricingPlanPageResponseData](&page, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}
	return &tenant.PricingPlanPageResponse{Data: data}, nil
}
