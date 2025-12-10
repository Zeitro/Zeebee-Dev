package resource

import (
	"context"

	"github.com/jinzhu/copier"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type ResourceTreeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// tree
func NewResourceTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResourceTreeLogic {
	return &ResourceTreeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ResourceTreeLogic) ResourceTree(req *types.ResourceTreeRequest) (resp *types.ResourceTreeResponse, err error) {
	resp = &types.ResourceTreeResponse{}
	rpcReq := &admin.ResourceTreeRequest{
		Data: &admin.ResourceTreeRequestData{},
	}
	err = copier.CopyWithOption(rpcReq.Data, req, copier.Option{DeepCopy: true})
	if err != nil {
		return nil, errorx.GetFailedTransReqErr(err)
	}

	rpcResp, err := l.svcCtx.ResourceRpcClient.ResourceTree(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// trans response
	err = copier.CopyWithOption(&resp.List, &rpcResp.Data, copier.Option{DeepCopy: true})
	if err != nil {
		return nil, errorx.GetFailedTransReqErr(err)
	}
	return resp, err
}
