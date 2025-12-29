package common

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUuidLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// get uuid
func NewGetUuidLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUuidLogic {
	return &GetUuidLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUuidLogic) GetUuid(req *types.EmptyRequest) (resp *types.GetUuidResponse, err error) {
	rpcReq := &borrower.EmptyRequest{}

	rpcResp, err := l.svcCtx.CommonRpcClient.GetUuid(l.ctx, rpcReq)
	if err != nil {
		return
	}
	return &types.GetUuidResponse{
		Uuid: rpcResp.Data,
	}, nil
}
