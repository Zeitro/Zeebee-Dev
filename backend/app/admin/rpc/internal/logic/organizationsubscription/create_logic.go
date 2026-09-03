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

type CreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *CreateLogic) Create(in *admin.OrganizationSubscriptionCreateRequest) (resp *admin.OrganizationSubscriptionDetailResponse, err error) {

	// organization_id comes from API layer context, passed via request metadata
	data, err := copierx.Trans[admin.OrganizationSubscriptionCreateRequestData, tenant.OrganizationSubscriptionCreateRequestData](in.Data)
	if err != nil {
		return nil, err
	}

	data.UserType = in.Data.Type
	if data.UserType == "" {
		data.UserType = string(constants.UserTypeAdmin)
	}
	result, err := l.svcCtx.TenantClient.OrganizationSubscription.Create(l.ctx, &tenant.OrganizationSubscriptionCreateRequest{Data: data})
	if err != nil {
		return nil, err
	}

	return &admin.OrganizationSubscriptionDetailResponse{
		Data: result.Data,
	}, nil
}
