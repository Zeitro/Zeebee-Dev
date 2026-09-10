package organizationlogic

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

func (l *PageLogic) Page(in *admin.OrganizationPageRequest) (resp *admin.OrganizationPageResponse, err error) {

	pageReq, err := copierx.Trans[admin.OrganizationPageRequestData, tenant.OrganizationPageRequestData](in.Data)
	if err != nil {
		return nil, err
	}
	if pageReq.UserType == "" {
		pageReq.UserType = string(constants.UserTypeAdmin)
		if pageReq.Filter != nil && pageReq.Filter.OrganizationType != "" {
			pageReq.UserType = pageReq.Filter.OrganizationType
		}
	}

	pageResult, err := l.svcCtx.TenantClient.Organization.Page(l.ctx, &tenant.OrganizationPageRequest{Data: pageReq})
	if err != nil {
		return nil, err
	}

	data, err := copierx.Trans[tenant.OrganizationPageResponseData, admin.OrganizationPageResponseData](pageResult.Data, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}

	return &admin.OrganizationPageResponse{
		Data: data,
	}, nil
}
