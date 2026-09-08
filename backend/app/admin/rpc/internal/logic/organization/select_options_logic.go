package organizationlogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/app/tenant/rpc/pb/tenant"

	"zeebee-backend/common/errorx"
)

type SelectOptionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewSelectOptionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SelectOptionsLogic {
	return &SelectOptionsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *SelectOptionsLogic) SelectOptions(in *admin.OrganizationSelectOptionsRequest) (resp *admin.OrganizationSelectOptionsResponse, err error) {

	if in.Data == nil {
		return nil, errorx.New("Data is required")
	}

	organization, err := l.svcCtx.TenantClient.Organization.SelectOptions(l.ctx, &tenant.OrganizationSelectOptionsRequest{Data: &tenant.OrganizationSelectOptionsRequestData{
		UserType: in.Data.Type,
	}})
	if err != nil {
		return nil, err
	}

	return &admin.OrganizationSelectOptionsResponse{
		Data: organization.Data,
	}, nil
}
