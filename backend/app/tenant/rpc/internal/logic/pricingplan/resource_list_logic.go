package pricingplanlogic

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

type ResourceListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewResourceListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResourceListLogic {
	return &ResourceListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *ResourceListLogic) ResourceList(in *tenant.PricingPlanResourceListRequest) (resp *tenant.PricingPlanResourceListResponse, err error) {
	if in.Data == nil {
		return nil, errorx.ErrorEmptyRequestData
	}
	if _, err = logic.ParseUserType(in.Data.UserType); err != nil {
		return nil, err
	}

	resources, err := l.baseLogic.PricingPlanLogic.ResourceList(l.ctx, in.Data.Id)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	var list []*proto.ResourceTreeEntity
	for _, r := range resources {
		item, innerErr := copierx.Trans[dto.PricingPlanResourceEntityDto, proto.ResourceTreeEntity](r)
		if innerErr != nil {
			return nil, innerErr
		}
		list = append(list, item)
	}
	return &tenant.PricingPlanResourceListResponse{Data: list}, nil
}
