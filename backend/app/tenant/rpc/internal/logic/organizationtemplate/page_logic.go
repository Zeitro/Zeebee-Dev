package organizationtemplatelogic

import (
	"context"

	"zeebee-backend/app/tenant/rpc/internal/logic"
	"zeebee-backend/app/tenant/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
)

type PageLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PageLogic {
	return &PageLogic{ctx: ctx, svcCtx: svcCtx, baseLogic: logic.GetBaseLogic(svcCtx)}
}

// Page returns a paginated list of organization templates.
func (l *PageLogic) Page(in *tenant.OrganizationTemplatePageRequest) (resp *tenant.OrganizationTemplatePageResponse, err error) {
	if in.Data == nil {
		return nil, errorx.ErrorEmptyRequestData
	}
	if _, err = logic.ParseUserType(in.Data.UserType); err != nil {
		return
	}

	pageReq := dto.OrganizationTemplatePageRequestDto{
		OrganizationPageRequestDto: dto.OrganizationPageRequestDto{
			BasePageRequestDto: dto.BasePageRequestDto{Page: in.Data.Page, Limit: in.Data.Limit},
		},
		Sorter: dto.BasePageSorter{Column: in.Data.Sorter.Column, Sort: in.Data.Sorter.Sort},
		Filter: dto.OrganizationTemplatePageFilter{
			Keyword:        in.Data.Filter.Keyword,
			Type:           constants.UserType(in.Data.Filter.Type),
			Category:       constants.TemplateCategory(in.Data.Filter.Category),
			Status:         constants.CommonStatus(in.Data.Filter.Status),
			OrganizationId: in.Data.Filter.OrganizationId,
		},
	}

	page, err := l.baseLogic.OrganizationTemplateLogic.Page(l.ctx, pageReq)
	if err != nil {
		return
	}

	// copierx converts the whole page result (incl. the element slice) in one pass.
	data, err := copierx.Trans[dto.PageResultDto[dto.OrganizationTemplatePageEntityDto], tenant.OrganizationTemplatePageResponseData](&page, copierx.TimeToRFC3339Str())
	if err != nil {
		return
	}
	return &tenant.OrganizationTemplatePageResponse{Data: data}, nil
}
