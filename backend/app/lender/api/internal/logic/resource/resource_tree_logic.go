package resource

import (
	"context"

	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/proto"

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

func (l *ResourceTreeLogic) ResourceTree(req *types.EmptyRequest) (resp *types.ResourceTreeResponse, err error) {
	resp = &types.ResourceTreeResponse{}
	rpcReq := &lender.ResourceTreeRequest{
		Data: &lender.ResourceTreeRequestData{
			Type: string(constants.UserTypeLender),
		},
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
