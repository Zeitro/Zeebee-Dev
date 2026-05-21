package borrower_property_transaction

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
)

type BorrowerPropertyTransactionSwitchStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// switch transaction status
func NewBorrowerPropertyTransactionSwitchStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BorrowerPropertyTransactionSwitchStatusLogic {
	return &BorrowerPropertyTransactionSwitchStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BorrowerPropertyTransactionSwitchStatusLogic) BorrowerPropertyTransactionSwitchStatus(req *types.BorrowerPropertyTransactionSwitchStatusRequest) (resp *types.EmptyResponse, err error) {
	rpcReq := &borrower.BorrowerPropertyTransactionSwitchStatusRequest{
		Data: &borrower.BorrowerPropertyTransactionSwitchStatusRequestData{},
	}

	rpcReq.Data, err = copierx.Trans[types.BorrowerPropertyTransactionSwitchStatusRequest, borrower.BorrowerPropertyTransactionSwitchStatusRequestData](req)
	if err != nil {
		return nil, err
	}

	_, err = l.svcCtx.BorrowerPropertyTransactionRpcClient.SwitchStatus(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	return &types.EmptyResponse{}, nil
}
