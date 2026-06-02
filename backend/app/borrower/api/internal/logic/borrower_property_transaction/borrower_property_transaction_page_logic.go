package borrower_property_transaction

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
)

type BorrowerPropertyTransactionPageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// page transactions
func NewBorrowerPropertyTransactionPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BorrowerPropertyTransactionPageLogic {
	return &BorrowerPropertyTransactionPageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BorrowerPropertyTransactionPageLogic) BorrowerPropertyTransactionPage(req *types.BorrowerPropertyTransactionPageRequest) (resp *types.BorrowerPropertyTransactionPageResponse, err error) {
	rpcReq := &borrower.BorrowerPropertyTransactionPageRequest{}

	rpcReq.Data, err = copierx.Trans[types.BorrowerPropertyTransactionPageRequest, borrower.BorrowerPropertyTransactionPageRequestData](req)
	if err != nil {
		return nil, err
	}

	rpcResp, err := l.svcCtx.BorrowerPropertyTransactionRpcClient.Page(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	resp, err = copierx.Trans[borrower.BorrowerPropertyTransactionPageResponseData, types.BorrowerPropertyTransactionPageResponse](rpcResp.Data)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
