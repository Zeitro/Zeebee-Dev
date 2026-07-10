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
	defer func() {
		if err != nil {
			l.svcCtx.Logger.WithContext(l.ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
			).Error(err)
		}
	}()
	req := dto.PricingPlanPageRequestDto{
		BasePageRequestDto: dto.BasePageRequestDto{Page: in.Data.Page, Limit: in.Data.Limit},
	}
	if in.Data.Filter != nil {
		req.Filter = dto.PricingPlanPageRequestFilter{
			Keyword:  in.Data.Filter.Keyword,
			Status:   in.Data.Filter.Status,
			Type:     in.Data.Filter.Type,
			PlanType: in.Data.Filter.PlanType,
		}
	}
	if in.Data.Sorter != nil {
		req.Sorter = dto.BasePageSorter{Column: in.Data.Sorter.Column, Sort: in.Data.Sorter.Sort}
	}

	pageResult, err := l.baseLogic.PricingPlanLogic.Page(l.ctx, req)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	pageDto := dto.PageResultDto[dto.PricingPlanPageEntityDto](pageResult)
	protoData, err := copierx.Trans[dto.PageResultDto[dto.PricingPlanPageEntityDto], admin.PricingPlanPageResponseData](
		&pageDto, copierx.TimeToRFC3339Str(),
	)
	if err != nil {
		return nil, err
	}

	return &admin.PricingPlanPageResponse{Data: protoData}, nil
}
