package organizationlogic

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

// Page returns a paginated list of organizations.
func (l *PageLogic) Page(in *tenant.OrganizationPageRequest) (resp *tenant.OrganizationPageResponse, err error) {
	if in.Data == nil {
		return nil, errorx.ErrorEmptyRequestData
	}
	userType, err := logic.ParseUserType(in.Data.UserType)
	if err != nil {
		return
	}

	pageReq := dto.OrganizationPageRequestDto{
		BasePageRequestDto: dto.BasePageRequestDto{Page: in.Data.Page, Limit: in.Data.Limit},
		Sorter:             dto.BasePageSorter{Column: in.Data.Sorter.Column, Sort: in.Data.Sorter.Sort},
		Filter: dto.OrganizationPageFilter{
			Keyword:          in.Data.Filter.Keyword,
			OrganizationType: userType,
			Status:           constants.OrganizationStatus(in.Data.Filter.Status),
		},
	}

	page, err := l.baseLogic.OrganizationLogic.Page(l.ctx, pageReq)
	if err != nil {
		return
	}

	data, err := copierx.Trans[dto.PageResultDto[dto.OrganizationPageEntityDto], tenant.OrganizationPageResponseData](&page, copierx.TimeToRFC3339Str())
	if err != nil {
		return
	}
	return &tenant.OrganizationPageResponse{Data: data}, nil
}
