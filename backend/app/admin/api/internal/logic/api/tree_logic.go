package api

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/proto"
)

type TreeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// tree
func NewTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TreeLogic {
	return &TreeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TreeLogic) Tree(req *types.ApiTreeRequest) (resp *types.ApiTreeResponse, err error) {
	resp = &types.ApiTreeResponse{}
	rpcReq := &admin.ApiTreeRequest{
		Data: &admin.ApiTreeRequestData{},
	}

	rpcReq.Data, err = copierx.Trans[types.ApiTreeRequest, admin.ApiTreeRequestData](req)
	if err != nil {
		return nil, err
	}

	rpcResp, err := l.svcCtx.ApiRpcClient.ApiTree(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// trans response
	list, err := copierx.Trans[[]*proto.ApiTreeEntity, []*types.ApiTreeEntity](&rpcResp.Data)
	if err != nil {
		return nil, err
	}
	resp.List = *list

	return resp, err
}
