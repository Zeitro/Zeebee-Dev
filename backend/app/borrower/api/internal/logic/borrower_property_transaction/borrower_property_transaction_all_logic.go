package borrower_property_transaction

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/proto"
)

type BorrowerPropertyTransactionAllLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// all transactions for select
func NewBorrowerPropertyTransactionAllLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BorrowerPropertyTransactionAllLogic {
	return &BorrowerPropertyTransactionAllLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BorrowerPropertyTransactionAllLogic) BorrowerPropertyTransactionAll(req *types.EmptyRequest) (resp *types.BorrowerPropertyTransactionAllResponse, err error) {
	rpcResp, err := l.svcCtx.BorrowerPropertyTransactionRpcClient.All(l.ctx, &proto.EmptyRequest{})
	if err != nil {
		return nil, err
	}

	resp = &types.BorrowerPropertyTransactionAllResponse{}

	listPtr, err := copierx.Trans[[]*proto.BorrowerPropertyTransactionOptionEntity, []types.BorrowerPropertyTransactionOption](&rpcResp.Data)
	if err != nil {
		return nil, err
	}
	if listPtr == nil {
		return resp, nil
	}

	resp.Lists = *listPtr

	return resp, nil
}
