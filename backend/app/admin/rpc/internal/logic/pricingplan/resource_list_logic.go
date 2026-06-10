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

type ResourceListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewResourceListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResourceListLogic {
	return &ResourceListLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *ResourceListLogic) ResourceList(in *admin.PricingPlanResourceListRequest) (resp *admin.PricingPlanResourceListResponse, err error) {

	resources, err := l.baseLogic.PricingPlanLogic.ResourceList(l.ctx, in.Data.Id)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	var list []*admin.ResourceTreeEntity
	for _, r := range resources {
		item, innerErr := copierx.Trans[dto.PricingPlanResourceEntityDto, admin.ResourceTreeEntity](r)
		if innerErr != nil {
			return nil, innerErr
		}
		list = append(list, item)
	}

	return &admin.PricingPlanResourceListResponse{Data: list}, nil
}
