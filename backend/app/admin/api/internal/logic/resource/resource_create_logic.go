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

type ResourceCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// create
func NewResourceCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResourceCreateLogic {
	return &ResourceCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ResourceCreateLogic) ResourceCreate(req *types.ResourceCreateRequest) (resp *types.ResourceEntity, err error) {
	resp = &types.ResourceEntity{}
	rpcReq := &admin.ResourceCreateRequest{
		Data: &admin.ResourceCreateRequestData{},
	}

	rpcReq.Data, err = copierx.Trans[types.ResourceCreateRequest, admin.ResourceCreateRequestData](req)
	if err != nil {
		return nil, err
	}

	rpcResp, err := l.svcCtx.ResourceRpcClient.Create(l.ctx, rpcReq)
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
