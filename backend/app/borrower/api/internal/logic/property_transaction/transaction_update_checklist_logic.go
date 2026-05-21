package property_transaction

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
)

type TransactionUpdateChecklistLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTransactionUpdateChecklistLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TransactionUpdateChecklistLogic {
	return &TransactionUpdateChecklistLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TransactionUpdateChecklistLogic) TransactionUpdateChecklist(req *types.TransactionChecklistRequest) (resp *types.EmptyResponse, err error) {
	rpcReq := &borrower.BorrowerPropertyTransactionChecklistRequest{
		Data: &borrower.BorrowerPropertyTransactionChecklistRequestData{},
	}

	rpcReq.Data, err = copierx.Trans[types.TransactionChecklistRequest, borrower.BorrowerPropertyTransactionChecklistRequestData](req)
	if err != nil {
		return nil, err
	}

	_, err = l.svcCtx.BorrowerPropertyTransactionRpcClient.UpdateChecklist(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	return &types.EmptyResponse{}, nil
}
