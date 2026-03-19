package api

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
)

type DeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// delete
func NewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogic {
	return &DeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteLogic) Delete(req *types.ApiDeleteRequest) (resp *types.EmptyResponse, err error) {
	resp = &types.EmptyResponse{}
	rpcReq := &admin.ApiDeleteRequest{
		Data: &admin.ApiDeleteRequestData{},
	}
	rpcReq.Data, err = copierx.Trans[types.ApiDeleteRequest, admin.ApiDeleteRequestData](req)
	if err != nil {
		return nil, err
	}
	_, err = l.svcCtx.ApiRpcClient.Delete(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	return resp, err
}
