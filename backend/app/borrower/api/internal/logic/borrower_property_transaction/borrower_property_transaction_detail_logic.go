package borrower_property_transaction

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
)

type BorrowerPropertyTransactionDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// transaction detail
func NewBorrowerPropertyTransactionDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BorrowerPropertyTransactionDetailLogic {
	return &BorrowerPropertyTransactionDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BorrowerPropertyTransactionDetailLogic) BorrowerPropertyTransactionDetail(req *types.BorrowerPropertyTransactionDetailRequest) (resp *types.BorrowerPropertyTransactionDetailResponse, err error) {
	rpcResp, err := l.svcCtx.BorrowerPropertyTransactionRpcClient.Detail(l.ctx, &borrower.BorrowerPropertyTransactionDetailRequest{
		Data: &borrower.BorrowerPropertyTransactionDetailRequestData{
			Id: req.Id,
		},
	})
	if err != nil {
		return nil, err
	}

	resp = &types.BorrowerPropertyTransactionDetailResponse{}
	resp, err = copierx.Trans[borrower.BorrowerPropertyTransactionEntity, types.BorrowerPropertyTransactionDetailResponse](rpcResp.Data)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
