package borrower_property

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
)

type PropertyAllLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// all
func NewPropertyAllLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PropertyAllLogic {
	return &PropertyAllLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PropertyAllLogic) PropertyAll(req *types.EmptyRequest) (resp *types.PropertyAllResponse, err error) {
	rpcResp, err := l.svcCtx.BorrowerPropertyRpcClient.All(l.ctx, &borrower.PropertyAllRequest{
		Data: &borrower.PropertyAllRequestData{},
	})
	if err != nil {
		return nil, err
	}

	resp = &types.PropertyAllResponse{}
	resp.Lists = make([]*types.BorrowerPropertyEntity, 0, len(rpcResp.Data))
	listPtr, err := copierx.Trans[[]*borrower.BorrowerPropertyEntity, []*types.BorrowerPropertyEntity](&rpcResp.Data)
	if err != nil {
		return nil, err
	}
	resp.Lists = *listPtr

	return resp, nil
}
