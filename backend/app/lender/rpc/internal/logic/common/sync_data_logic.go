package commonlogic

import (
	"context"

	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/domains/zeebee/proto"
)

type SyncDataLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSyncDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SyncDataLogic {
	return &SyncDataLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SyncDataLogic) SyncData(*proto.EmptyRequest) (*proto.EmptyResponse, error) {
	_, err := l.svcCtx.TenantClient.OrganizationSetting.SyncLenderDefaults(l.ctx, &tenant.OrganizationSettingSyncLenderDefaultsRequest{Data: &tenant.OrganizationSettingSyncLenderDefaultsRequestData{}})
	if err != nil {
		return nil, err
	}
	return &proto.EmptyResponse{}, nil
}
