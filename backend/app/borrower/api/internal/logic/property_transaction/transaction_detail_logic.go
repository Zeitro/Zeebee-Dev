package property_transaction

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
)

type TransactionDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTransactionDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TransactionDetailLogic {
	return &TransactionDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TransactionDetailLogic) TransactionDetail(req *types.TransactionDetailRequest) (resp *types.TransactionDetailResponse, err error) {
	rpcResp, err := l.svcCtx.BorrowerPropertyTransactionRpcClient.Detail(l.ctx, &borrower.BorrowerPropertyTransactionDetailRequest{
		Data: &borrower.BorrowerPropertyTransactionDetailRequestData{
			Id: req.Id,
		},
	})
	if err != nil {
		return nil, err
	}

	resp = &types.TransactionDetailResponse{}
	resp, err = copierx.Trans[borrower.BorrowerPropertyTransactionEntity, types.TransactionDetailResponse](rpcResp.Data)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
