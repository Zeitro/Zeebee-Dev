package common

import (
	"context"

	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"
	"zeebee-backend/app/lender/rpc/pb/lender"

	"github.com/zeromicro/go-zero/core/logx"
)

type SyncDataLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// sync data
func NewSyncDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SyncDataLogic {
	return &SyncDataLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SyncDataLogic) SyncData(req *types.EmptyRequest) (resp *types.EmptyResponse, err error) {
	_, err = l.svcCtx.CommonRpcClient.SyncData(l.ctx, &lender.EmptyRequest{})
	if err != nil {
		return nil, err
	}
	return
}
