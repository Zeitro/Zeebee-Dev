package resourcelogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/proto"
)

type SyncResourceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSyncResourceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SyncResourceLogic {
	return &SyncResourceLogic{ctx: ctx, svcCtx: svcCtx}
}

// SyncResource proxies to tenant; tenant upserts the manifest into its registry.
func (l *SyncResourceLogic) SyncResource(in *admin.SyncResourceRequest) (resp *proto.EmptyResponse, err error) {
	return l.svcCtx.TenantClient.Resource.SyncResource(l.ctx, &tenant.SyncResourceRequest{
		Data: &tenant.SyncResourceRequestData{
			UserType:     string(constants.UserTypeAdmin),
			ResourceJson: in.Data.ResourceJson,
		},
	})
}
