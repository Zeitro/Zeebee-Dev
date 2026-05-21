package property_transaction

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
)

type TransactionAllLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTransactionAllLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TransactionAllLogic {
	return &TransactionAllLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TransactionAllLogic) TransactionAll(req *types.EmptyRequest) (resp *types.TransactionAllResponse, err error) {
	rpcResp, err := l.svcCtx.BorrowerPropertyTransactionRpcClient.All(l.ctx, &borrower.BorrowerPropertyTransactionAllRequest{})
	if err != nil {
		return nil, err
	}

	resp = &types.TransactionAllResponse{}
	for _, opt := range rpcResp.GetData() {
		item, innerErr := copierx.Trans[borrower.BorrowerPropertyTransactionOptionEntity, types.TransactionOption](opt)
		if innerErr != nil {
			continue
		}
		resp.Lists = append(resp.Lists, *item)
	}

	return resp, nil
}
