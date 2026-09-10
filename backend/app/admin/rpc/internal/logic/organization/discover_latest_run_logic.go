package organizationlogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/proto"
)

type DiscoverLatestRunLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDiscoverLatestRunLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DiscoverLatestRunLogic {
	return &DiscoverLatestRunLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *DiscoverLatestRunLogic) DiscoverLatestRun(in *proto.EmptyRequest) (*admin.LenderDiscoverRunResponse, error) {
	resp, err := l.svcCtx.TenantClient.Organization.DiscoverLatestRun(l.ctx, in)
	if err != nil {
		return nil, err
	}
	return copierx.Trans[tenant.LenderDiscoverRunResponse, admin.LenderDiscoverRunResponse](resp)
}
