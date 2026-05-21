package property_transaction

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
)

type TransactionCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTransactionCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TransactionCreateLogic {
	return &TransactionCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TransactionCreateLogic) TransactionCreate(req *types.TransactionCreateRequest) (resp *types.TransactionCreateResponse, err error) {
	rpcReq := &borrower.BorrowerPropertyTransactionCreateRequest{
		Data: &borrower.BorrowerPropertyTransactionCreateRequestData{},
	}

	rpcReq.Data, err = copierx.Trans[types.TransactionCreateRequest, borrower.BorrowerPropertyTransactionCreateRequestData](req)
	if err != nil {
		return nil, err
	}

	rpcResp, err := l.svcCtx.BorrowerPropertyTransactionRpcClient.Create(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	resp = &types.TransactionCreateResponse{}
	resp, err = copierx.Trans[borrower.BorrowerPropertyTransactionEntity, types.TransactionCreateResponse](rpcResp.Data)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
