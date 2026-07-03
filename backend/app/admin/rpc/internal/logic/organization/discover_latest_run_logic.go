package organizationlogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/domains/zeebee/proto"
)

type DiscoverLatestRunLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDiscoverLatestRunLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DiscoverLatestRunLogic {
	return &DiscoverLatestRunLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

// DiscoverLatestRun returns an empty stub (runs table has been removed; endpoint kept for backward compatibility).
func (l *DiscoverLatestRunLogic) DiscoverLatestRun(_ *proto.EmptyRequest) (*admin.LenderDiscoverRunResponse, error) {
	return &admin.LenderDiscoverRunResponse{Data: &admin.LenderDiscoverRunEntity{}}, nil
}
