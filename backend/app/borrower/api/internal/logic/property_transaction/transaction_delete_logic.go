package property_transaction

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
)

type TransactionDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTransactionDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TransactionDeleteLogic {
	return &TransactionDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TransactionDeleteLogic) TransactionDelete(req *types.TransactionDeleteRequest) (resp *types.EmptyResponse, err error) {
	rpcReq := &borrower.BorrowerPropertyTransactionDeleteRequest{
		Data: &borrower.BorrowerPropertyTransactionDeleteRequestData{},
	}

	rpcReq.Data, err = copierx.Trans[types.TransactionDeleteRequest, borrower.BorrowerPropertyTransactionDeleteRequestData](req)
	if err != nil {
		return nil, err
	}

	_, err = l.svcCtx.BorrowerPropertyTransactionRpcClient.Delete(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	return &types.EmptyResponse{}, nil
}
