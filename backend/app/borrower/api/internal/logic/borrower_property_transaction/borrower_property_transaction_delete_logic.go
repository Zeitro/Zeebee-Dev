package borrower_property_transaction

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
)

type BorrowerPropertyTransactionDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// delete transaction
func NewBorrowerPropertyTransactionDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BorrowerPropertyTransactionDeleteLogic {
	return &BorrowerPropertyTransactionDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BorrowerPropertyTransactionDeleteLogic) BorrowerPropertyTransactionDelete(req *types.BorrowerPropertyTransactionDeleteRequest) (resp *types.EmptyResponse, err error) {
	rpcReq := &borrower.BorrowerPropertyTransactionDeleteRequest{
		Data: &borrower.BorrowerPropertyTransactionDeleteRequestData{},
	}

	rpcReq.Data, err = copierx.Trans[types.BorrowerPropertyTransactionDeleteRequest, borrower.BorrowerPropertyTransactionDeleteRequestData](req)
	if err != nil {
		return nil, err
	}

	_, err = l.svcCtx.BorrowerPropertyTransactionRpcClient.Delete(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	return &types.EmptyResponse{}, nil
}
