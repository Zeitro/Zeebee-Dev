package common

import (
	"context"

	"github.com/pkg/errors"

	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"
	"zeebee-backend/app/lender/rpc/pb/lender"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUsStatesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// get us states
func NewGetUsStatesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUsStatesLogic {
	return &GetUsStatesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUsStatesLogic) GetUsStates(req *types.EmptyRequest) (resp *types.GetUsStatesResponse, err error) {
	var rpcResp *lender.StringSliceResponse
	rpcResp, err = l.svcCtx.CommonRpcClient.GetUsState(l.ctx, &lender.EmptyRequest{})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &types.GetUsStatesResponse{
		List: rpcResp.Data,
	}, nil
}
