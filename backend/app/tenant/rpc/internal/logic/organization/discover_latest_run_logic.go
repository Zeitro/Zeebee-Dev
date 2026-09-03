package organizationlogic

import (
	"context"

	"zeebee-backend/app/tenant/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/domains/zeebee/proto"
)

type DiscoverLatestRunLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDiscoverLatestRunLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DiscoverLatestRunLogic {
	return &DiscoverLatestRunLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DiscoverLatestRunLogic) DiscoverLatestRun(in *proto.EmptyRequest) (resp *tenant.LenderDiscoverRunResponse, err error) {
	return &tenant.LenderDiscoverRunResponse{Data: &tenant.LenderDiscoverRunEntity{}}, nil
}
