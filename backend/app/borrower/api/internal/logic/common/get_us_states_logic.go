package common

import (
	"context"

	"github.com/jinzhu/copier"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/errorx"

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
	var rpcResp *borrower.UsStateResponse
	resp = &types.GetUsStatesResponse{
		List: make([]types.UsStateEntry, 0),
	}
	rpcResp, err = l.svcCtx.CommonRpcClient.GetUsState(l.ctx, &borrower.EmptyRequest{})
	if err != nil {
		return nil, err
	}
	err = copier.Copy(&resp.List, rpcResp.Data)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	return
}
