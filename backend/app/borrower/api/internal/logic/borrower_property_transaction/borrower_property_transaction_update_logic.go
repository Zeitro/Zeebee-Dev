package borrower_property_transaction

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
)

type BorrowerPropertyTransactionUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// update transaction
func NewBorrowerPropertyTransactionUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BorrowerPropertyTransactionUpdateLogic {
	return &BorrowerPropertyTransactionUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BorrowerPropertyTransactionUpdateLogic) BorrowerPropertyTransactionUpdate(req *types.BorrowerPropertyTransactionUpdateRequest) (resp *types.EmptyResponse, err error) {
	rpcReq := &borrower.BorrowerPropertyTransactionUpdateRequest{
		Data: &borrower.BorrowerPropertyTransactionUpdateRequestData{},
	}

	rpcReq.Data, err = copierx.Trans[types.BorrowerPropertyTransactionUpdateRequest, borrower.BorrowerPropertyTransactionUpdateRequestData](req)
	if err != nil {
		return nil, err
	}

	_, err = l.svcCtx.BorrowerPropertyTransactionRpcClient.Update(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	return &types.EmptyResponse{}, nil
}
