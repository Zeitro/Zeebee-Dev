package organizationlogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/copierx"
)

type DiscoverUpdateSettingsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDiscoverUpdateSettingsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DiscoverUpdateSettingsLogic {
	return &DiscoverUpdateSettingsLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *DiscoverUpdateSettingsLogic) DiscoverUpdateSettings(in *admin.LenderDiscoverUpdateSettingsRequest) (*admin.LenderDiscoverSettingsResponse, error) {
	req, err := copierx.Trans[admin.LenderDiscoverUpdateSettingsRequest, tenant.LenderDiscoverUpdateSettingsRequest](in)
	if err != nil {
		return nil, err
	}
	resp, err := l.svcCtx.TenantClient.Organization.DiscoverUpdateSettings(l.ctx, req)
	if err != nil {
		return nil, err
	}
	return copierx.Trans[tenant.LenderDiscoverSettingsResponse, admin.LenderDiscoverSettingsResponse](resp)
}
