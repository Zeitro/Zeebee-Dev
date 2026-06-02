package borrower_property_transaction

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
)

type OperationListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// transaction operation list
func NewOperationListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OperationListLogic {
	return &OperationListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OperationListLogic) OperationList(req *types.TransactionOperationListRequest) (resp *types.TransactionOperationListResponse, err error) {
	rpcReq := &borrower.TransactionOperationListRequest{}
	rpcReq.Data, err = copierx.Trans[types.TransactionOperationListRequest, borrower.TransactionOperationListRequestData](req)
	if err != nil {
		return nil, err
	}

	rpcResp, err := l.svcCtx.BorrowerPropertyTransactionRpcClient.OperationList(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	resp, err = copierx.Trans[borrower.TransactionOperationListResponseData, types.TransactionOperationListResponse](rpcResp.Data, copierx.StrToStrMap())
	if err != nil {
		return nil, err
	}

	return resp, nil
}
