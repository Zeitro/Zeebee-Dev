package commonlogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/static"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/proto"
)

type SyncDataLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSyncDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SyncDataLogic {
	return &SyncDataLogic{ctx: ctx, svcCtx: svcCtx}
}

// SyncData keeps the admin endpoint as an orchestration shim; tenant owns both registries.
func (l *SyncDataLogic) SyncData(_ *proto.EmptyRequest) (*proto.EmptyResponse, error) {
	if _, err := l.svcCtx.TenantClient.Api.SyncApi(l.ctx, &tenant.SyncApiRequest{Data: &tenant.SyncApiRequestData{
		UserType: string(constants.UserTypeAdmin),
		ApiJson:  string(static.ApiStatic),
	}}); err != nil {
		return nil, err
	}
	return l.svcCtx.TenantClient.Resource.SyncResource(l.ctx, &tenant.SyncResourceRequest{Data: &tenant.SyncResourceRequestData{
		UserType:     string(constants.UserTypeAdmin),
		ResourceJson: string(static.ResourceStatic),
	}})
}
