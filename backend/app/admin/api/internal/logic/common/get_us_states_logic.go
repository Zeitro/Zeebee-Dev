package common

import (
	"context"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/proto"

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
		Logger: svcCtx.Logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUsStatesLogic) GetUsStates(req *types.EmptyRequest) (resp *types.GetUsStatesResponse, err error) {
	var rpcResp *proto.StringSliceResponse
	rpcResp, err = l.svcCtx.CommonRpcClient.GetUsState(l.ctx, &proto.EmptyRequest{})
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	return &types.GetUsStatesResponse{
		List: rpcResp.Data,
	}, nil
}
