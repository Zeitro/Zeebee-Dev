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

type ResourceTreeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// tree
func NewResourceTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResourceTreeLogic {
	return &ResourceTreeLogic{
		Logger: svcCtx.Logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ResourceTreeLogic) ResourceTree(req *types.ResourceTreeRequest) (resp *types.ResourceTreeResponse, err error) {
	resp = &types.ResourceTreeResponse{}
	rpcReq := &admin.ResourceTreeRequest{
		Data: &admin.ResourceTreeRequestData{},
	}

	rpcReq.Data, err = copierx.Trans[types.ResourceTreeRequest, admin.ResourceTreeRequestData](req)
	if err != nil {
		return nil, err
	}

	rpcResp, err := l.svcCtx.ResourceRpcClient.ResourceTree(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// trans response
	list, err := copierx.Trans[[]*proto.ResourceTreeEntity, []*types.ResourceTreeEntity](&rpcResp.Data)
	if err != nil {
		return nil, err
	}
	resp.List = *list

	return resp, err
}
