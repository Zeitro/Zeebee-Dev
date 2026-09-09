package organizationlogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/copierx"
)

type DiscoverRunLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDiscoverRunLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DiscoverRunLogic {
	return &DiscoverRunLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *DiscoverRunLogic) DiscoverRun(in *admin.LenderDiscoverRunRequest) (*admin.LenderDiscoverRunResponse, error) {
	req, err := copierx.Trans[admin.LenderDiscoverRunRequest, tenant.LenderDiscoverRunRequest](in)
	if err != nil {
		return nil, err
	}
	resp, err := l.svcCtx.TenantClient.Organization.DiscoverRun(l.ctx, req)
	if err != nil {
		return nil, err
	}
	return copierx.Trans[tenant.LenderDiscoverRunResponse, admin.LenderDiscoverRunResponse](resp)
}
