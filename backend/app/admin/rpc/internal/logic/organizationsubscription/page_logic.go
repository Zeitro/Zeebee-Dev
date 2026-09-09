package organizationsubscriptionlogic

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

func (l *PageLogic) Page(in *admin.OrganizationSubscriptionPageRequest) (resp *admin.OrganizationSubscriptionPageResponse, err error) {

	req, err := copierx.Trans[admin.OrganizationSubscriptionPageRequestData, tenant.OrganizationSubscriptionPageRequestData](in.Data)
	if err != nil {
		return nil, err
	}

	req.UserType = string(constants.UserTypeAdmin)
	pageResult, err := l.svcCtx.TenantClient.OrganizationSubscription.Page(l.ctx, &tenant.OrganizationSubscriptionPageRequest{Data: req})
	if err != nil {
		return nil, err
	}
	protoData, err := copierx.Trans[tenant.OrganizationSubscriptionPageResponseData, admin.OrganizationSubscriptionPageResponseData](pageResult.Data, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}

	return &admin.OrganizationSubscriptionPageResponse{Data: protoData}, nil
}
