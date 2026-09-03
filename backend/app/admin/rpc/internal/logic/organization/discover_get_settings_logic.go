package organizationlogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/proto"
)

type DiscoverGetSettingsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDiscoverGetSettingsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DiscoverGetSettingsLogic {
	return &DiscoverGetSettingsLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *DiscoverGetSettingsLogic) DiscoverGetSettings(in *proto.EmptyRequest) (*admin.LenderDiscoverSettingsResponse, error) {
	resp, err := l.svcCtx.TenantClient.Organization.DiscoverGetSettings(l.ctx, in)
	if err != nil {
		return nil, err
	}
	return copierx.Trans[tenant.LenderDiscoverSettingsResponse, admin.LenderDiscoverSettingsResponse](resp)
}
