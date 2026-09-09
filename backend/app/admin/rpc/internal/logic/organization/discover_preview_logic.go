package organizationlogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/proto"
)

type DiscoverPreviewLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDiscoverPreviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DiscoverPreviewLogic {
	return &DiscoverPreviewLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *DiscoverPreviewLogic) DiscoverPreview(in *proto.EmptyRequest) (*admin.LenderDiscoverPreviewResponse, error) {
	resp, err := l.svcCtx.TenantClient.Organization.DiscoverPreview(l.ctx, in)
	if err != nil {
		return nil, err
	}
	return copierx.Trans[tenant.LenderDiscoverPreviewResponse, admin.LenderDiscoverPreviewResponse](resp)
}
