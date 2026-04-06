package lender_discover

import (
	"context"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type LenderDiscoverLatestRunLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// lender discover latest run
func NewLenderDiscoverLatestRunLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LenderDiscoverLatestRunLogic {
	return &LenderDiscoverLatestRunLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LenderDiscoverLatestRunLogic) LenderDiscoverLatestRun(_ *types.EmptyRequest) (resp *types.LenderDiscoverRunResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	rpcResp, err := l.svcCtx.LenderDiscoverRpcClient.LatestRun(l.ctx, &admin.EmptyRequest{})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	return &types.LenderDiscoverRunResponse{
		Data: protoRunToTypes(rpcResp.GetData()),
	}, nil
}
