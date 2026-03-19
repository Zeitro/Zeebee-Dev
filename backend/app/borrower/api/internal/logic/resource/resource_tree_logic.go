package resource

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/constants"

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
	rpcReq := &borrower.ResourceTreeRequest{
		Data: &borrower.ResourceTreeRequestData{
			Type: string(constants.UserTypeBorrower),
		},
	}

	rpcResp, err := l.svcCtx.ResourceRpcClient.ResourceTree(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// trans response
	list, err := copierx.Trans[[]*borrower.ResourceTreeEntity, []*types.ResourceTreeEntity](&rpcResp.Data)
	if err != nil {
		return nil, err
	}
	resp.List = *list

	return resp, err
}
