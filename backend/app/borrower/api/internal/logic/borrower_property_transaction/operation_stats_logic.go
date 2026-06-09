package borrower_property_transaction

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
)

type OperationStatsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// transaction operation stats
func NewOperationStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OperationStatsLogic {
	return &OperationStatsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OperationStatsLogic) OperationStats(req *types.TransactionOperationStatsRequest) (resp *types.TransactionOperationStatsResponse, err error) {
	rpcReq := &borrower.TransactionOperationStatsRequest{}
	rpcReq.Data, err = copierx.Trans[types.TransactionOperationStatsRequest, borrower.TransactionOperationStatsRequestData](req)
	if err != nil {
		return nil, err
	}

	rpcResp, err := l.svcCtx.BorrowerPropertyTransactionRpcClient.OperationStats(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	resp, err = copierx.Trans[borrower.TransactionOperationStatsResponseData, types.TransactionOperationStatsResponse](rpcResp.Data)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
