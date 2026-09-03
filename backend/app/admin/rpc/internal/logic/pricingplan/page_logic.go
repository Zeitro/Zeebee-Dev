package pricingplanlogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/constants"
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

func (l *PageLogic) Page(in *admin.PricingPlanPageRequest) (resp *admin.PricingPlanPageResponse, err error) {

	req, err := copierx.Trans[admin.PricingPlanPageRequestData, tenant.PricingPlanPageRequestData](in.Data)
	if err != nil {
		return nil, err
	}
	req.UserType = string(constants.UserTypeAdmin)
	if req.Filter != nil && req.Filter.Type != "" {
		req.UserType = req.Filter.Type
	}
	pageResult, err := l.svcCtx.TenantClient.PricingPlan.Page(l.ctx, &tenant.PricingPlanPageRequest{Data: req})
	if err != nil {
		return nil, err
	}
	protoData, err := copierx.Trans[tenant.PricingPlanPageResponseData, admin.PricingPlanPageResponseData](pageResult.Data, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}

	return &admin.PricingPlanPageResponse{Data: protoData}, nil
}
