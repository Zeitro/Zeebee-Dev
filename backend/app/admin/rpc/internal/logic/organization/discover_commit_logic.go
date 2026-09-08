package organizationlogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/copierx"
)

type DiscoverCommitLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDiscoverCommitLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DiscoverCommitLogic {
	return &DiscoverCommitLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *DiscoverCommitLogic) DiscoverCommit(in *admin.LenderDiscoverCommitRequest) (*admin.LenderDiscoverRunResponse, error) {
	req, err := copierx.Trans[admin.LenderDiscoverCommitRequest, tenant.LenderDiscoverCommitRequest](in)
	if err != nil {
		return nil, err
	}
	resp, err := l.svcCtx.TenantClient.Organization.DiscoverCommit(l.ctx, req)
	if err != nil {
		return nil, err
	}
	return copierx.Trans[tenant.LenderDiscoverRunResponse, admin.LenderDiscoverRunResponse](resp)
}
