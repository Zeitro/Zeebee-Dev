package organizationsubscriptionlogic

import (
	"context"

	"zeebee-backend/app/tenant/rpc/internal/logic"
	"zeebee-backend/app/tenant/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
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

// Page returns a paginated list of organization subscriptions.
func (l *PageLogic) Page(in *tenant.OrganizationSubscriptionPageRequest) (resp *tenant.OrganizationSubscriptionPageResponse, err error) {
	if in.Data == nil {
		return nil, errorx.ErrorEmptyRequestData
	}
	userType, err := logic.ParseUserType(in.Data.UserType)
	if err != nil {
		return
	}

	statuses := make([]constants.OrganizationSubscriptionStatus, 0, len(in.Data.Filter.Status))
	for _, s := range in.Data.Filter.Status {
		statuses = append(statuses, constants.OrganizationSubscriptionStatus(s))
	}
	pageReq := dto.OrganizationSubscriptionPageRequestDto{
		BasePageRequestDto: dto.BasePageRequestDto{Page: in.Data.Page, Limit: in.Data.Limit},
		Sorter:             dto.BasePageSorter{Column: in.Data.Sorter.Column, Sort: in.Data.Sorter.Sort},
		Filter: dto.OrganizationSubscriptionPageRequestFilter{
			Keyword: in.Data.Filter.Keyword,
			Type:    userType,
			Status:  statuses,
		},
	}

	page, err := l.baseLogic.OrganizationSubscriptionLogic.Page(l.ctx, pageReq)
	if err != nil {
		return
	}

	data, err := copierx.Trans[gormx.PageResponse[dto.OrganizationSubscriptionPageEntity], tenant.OrganizationSubscriptionPageResponseData](&page, copierx.TimeToRFC3339Str())
	if err != nil {
		return
	}
	return &tenant.OrganizationSubscriptionPageResponse{Data: data}, nil
}
