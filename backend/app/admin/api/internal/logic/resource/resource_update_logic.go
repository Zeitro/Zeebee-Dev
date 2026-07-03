package resource

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/proto"
)

type ResourceUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// update
func NewResourceUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResourceUpdateLogic {
	return &ResourceUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ResourceUpdateLogic) ResourceUpdate(req *types.ResourceUpdateRequest) (resp *types.ResourceEntity, err error) {
	resp = &types.ResourceEntity{}
	rpcReq := &admin.ResourceUpdateRequest{
		Data: &admin.ResourceUpdateRequestData{},
	}

	rpcReq.Data, err = copierx.Trans[types.ResourceUpdateRequest, admin.ResourceUpdateRequestData](req)
	if err != nil {
		return nil, err
	}

	rpcResp, err := l.svcCtx.ResourceRpcClient.Update(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// trans response
	resp, err = copierx.Trans[proto.ResourceEntity, types.ResourceEntity](rpcResp.Data)
	if err != nil {
		return nil, err
	}
	return resp, err
}
