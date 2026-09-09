package resource

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
)

type ResourceDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// delete
func NewResourceDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResourceDeleteLogic {
	return &ResourceDeleteLogic{
		Logger: svcCtx.Logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ResourceDeleteLogic) ResourceDelete(req *types.ResourceDeleteRequest) (resp *types.EmptyResponse, err error) {
	resp = &types.EmptyResponse{}
	rpcReq := &admin.ResourceDeleteRequest{
		Data: &admin.ResourceDeleteRequestData{},
	}
	rpcReq.Data, err = copierx.Trans[types.ResourceDeleteRequest, admin.ResourceDeleteRequestData](req)
	if err != nil {
		return nil, err
	}
	_, err = l.svcCtx.ResourceRpcClient.Delete(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	return resp, err
}
