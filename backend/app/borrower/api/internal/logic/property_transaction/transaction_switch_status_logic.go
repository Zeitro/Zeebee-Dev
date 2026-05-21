package property_transaction

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
)

type TransactionSwitchStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTransactionSwitchStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TransactionSwitchStatusLogic {
	return &TransactionSwitchStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TransactionSwitchStatusLogic) TransactionSwitchStatus(req *types.TransactionSwitchStatusRequest) (resp *types.EmptyResponse, err error) {
	rpcReq := &borrower.BorrowerPropertyTransactionSwitchStatusRequest{
		Data: &borrower.BorrowerPropertyTransactionSwitchStatusRequestData{},
	}

	rpcReq.Data, err = copierx.Trans[types.TransactionSwitchStatusRequest, borrower.BorrowerPropertyTransactionSwitchStatusRequestData](req)
	if err != nil {
		return nil, err
	}

	_, err = l.svcCtx.BorrowerPropertyTransactionRpcClient.SwitchStatus(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	return &types.EmptyResponse{}, nil
}
