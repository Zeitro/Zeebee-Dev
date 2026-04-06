package lender_discover

import (
	"context"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type LenderDiscoverCommitLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// lender discover commit selected new lenders
func NewLenderDiscoverCommitLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LenderDiscoverCommitLogic {
	return &LenderDiscoverCommitLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LenderDiscoverCommitLogic) LenderDiscoverCommit(req *types.LenderDiscoverCommitRequest) (resp *types.LenderDiscoverRunResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	rpcResp, err := l.svcCtx.LenderDiscoverRpcClient.Commit(context.WithoutCancel(l.ctx), typesCommitRequestToProto(req))
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	return &types.LenderDiscoverRunResponse{
		Data: protoRunToTypes(rpcResp.GetData()),
	}, nil
}
