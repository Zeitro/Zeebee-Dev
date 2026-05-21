package property_transaction

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
)

type TransactionUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTransactionUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TransactionUpdateLogic {
	return &TransactionUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TransactionUpdateLogic) TransactionUpdate(req *types.TransactionUpdateRequest) (resp *types.EmptyResponse, err error) {
	rpcReq := &borrower.BorrowerPropertyTransactionUpdateRequest{
		Data: &borrower.BorrowerPropertyTransactionUpdateRequestData{},
	}

	rpcReq.Data, err = copierx.Trans[types.TransactionUpdateRequest, borrower.BorrowerPropertyTransactionUpdateRequestData](req)
	if err != nil {
		return nil, err
	}

	_, err = l.svcCtx.BorrowerPropertyTransactionRpcClient.Update(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	return &types.EmptyResponse{}, nil
}
