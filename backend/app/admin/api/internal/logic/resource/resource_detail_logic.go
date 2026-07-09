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

type ResourceDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// detail
func NewResourceDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResourceDetailLogic {
	return &ResourceDetailLogic{
		Logger: svcCtx.Logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ResourceDetailLogic) ResourceDetail(req *types.ResourceDetailRequest) (resp *types.ResourceEntity, err error) {
	resp = &types.ResourceEntity{}
	rpcReq := &admin.ResourceDetailRequest{
		Data: &admin.ResourceDetailRequestData{},
	}

	rpcReq.Data, err = copierx.Trans[types.ResourceDetailRequest, admin.ResourceDetailRequestData](req)
	if err != nil {
		return nil, err
	}

	rpcResp, err := l.svcCtx.ResourceRpcClient.Detail(l.ctx, rpcReq)
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
