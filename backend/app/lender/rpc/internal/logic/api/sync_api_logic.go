package apilogic

import (
	"context"

	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/proto"
)

type SyncApiLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSyncApiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SyncApiLogic {
	return &SyncApiLogic{ctx: ctx, svcCtx: svcCtx}
}

// SyncApi proxies to tenant; tenant upserts the manifest into its registry.
func (l *SyncApiLogic) SyncApi(in *lender.SyncApiRequest) (resp *proto.EmptyResponse, err error) {
	return l.svcCtx.TenantClient.Api.SyncApi(l.ctx, &tenant.SyncApiRequest{
		Data: &tenant.SyncApiRequestData{
			UserType: string(constants.UserTypeLender),
			ApiJson:  in.Data.ApiJson,
		},
	})
}
