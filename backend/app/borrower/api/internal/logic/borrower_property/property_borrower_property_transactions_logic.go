package borrower_property

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"

	"github.com/zeromicro/go-zero/core/logx"
)

type PropertyBorrowerPropertyTransactionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// borrower property transactions
func NewPropertyBorrowerPropertyTransactionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PropertyBorrowerPropertyTransactionsLogic {
	return &PropertyBorrowerPropertyTransactionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PropertyBorrowerPropertyTransactionsLogic) PropertyBorrowerPropertyTransactions(req *types.BorrowerPropertyTransactionListByPropertyRequest) (resp *types.BorrowerPropertyTransactionListByPropertyResponse, err error) {
	rpcResp, err := l.svcCtx.BorrowerPropertyTransactionRpcClient.ListByBorrowerProperty(l.ctx, &borrower.BorrowerPropertyTransactionListByPropertyRequest{
		Data: &borrower.BorrowerPropertyTransactionListByPropertyRequestData{
			BorrowerPropertyId: req.Id,
		},
	})
	if err != nil {
		return nil, err
	}

	resp = &types.BorrowerPropertyTransactionListByPropertyResponse{}
	list, err := copierx.Trans[[]*borrower.BorrowerPropertyTransactionEntity, []types.BorrowerPropertyTransactionEntity](&rpcResp.Data)
	if err != nil {
		return nil, err
	}
	if list != nil {
		resp.List = *list
	}

	return resp, nil
}
