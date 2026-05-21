package borrower_property_transaction

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
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
	rpcResp, err := l.svcCtx.BorrowerPropertyTransactionRpcClient.All(l.ctx, &borrower.BorrowerPropertyTransactionAllRequest{})
	if err != nil {
		return nil, err
	}

	resp = &types.BorrowerPropertyTransactionAllResponse{}
	for _, opt := range rpcResp.GetData() {
		item, innerErr := copierx.Trans[borrower.BorrowerPropertyTransactionOptionEntity, types.BorrowerPropertyTransactionOption](opt)
		if innerErr != nil {
			continue
		}
		resp.Lists = append(resp.Lists, *item)
	}

	return resp, nil
}
