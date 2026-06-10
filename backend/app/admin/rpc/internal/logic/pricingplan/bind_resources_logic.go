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

type BindResourcesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewBindResourcesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BindResourcesLogic {
	return &BindResourcesLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *BindResourcesLogic) BindResources(in *admin.PricingPlanBindResourcesRequest) (resp *admin.PricingPlanBindResourcesResponse, err error) {

	_, err = l.baseLogic.PricingPlanLogic.BindResources(l.ctx, in.Data.Id, in.Data.ResourceIds)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	// fetch plan resource list as tree
	resources, err := l.baseLogic.PricingPlanLogic.ResourceList(l.ctx, in.Data.Id)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	// convert resource bindings to proto
	var list []*admin.ResourceTreeEntity
	for _, r := range resources {
		item, innerErr := copierx.Trans[dto.PricingPlanResourceEntityDto, admin.ResourceTreeEntity](r)
		if innerErr != nil {
			return nil, innerErr
		}
		list = append(list, item)
	}

	return &admin.PricingPlanBindResourcesResponse{Data: list}, nil
}
