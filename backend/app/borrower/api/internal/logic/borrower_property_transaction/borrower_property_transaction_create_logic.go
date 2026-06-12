package borrower_property_transaction

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
)

type BorrowerPropertyTransactionCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// create transaction
func NewBorrowerPropertyTransactionCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BorrowerPropertyTransactionCreateLogic {
	return &BorrowerPropertyTransactionCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BorrowerPropertyTransactionCreateLogic) BorrowerPropertyTransactionCreate(req *types.BorrowerPropertyTransactionCreateRequest) (resp *types.BorrowerPropertyTransactionCreateResponse, err error) {
	rpcReq := &borrower.BorrowerPropertyTransactionCreateRequest{
		Data: &borrower.BorrowerPropertyTransactionCreateRequestData{},
	}

	rpcReq.Data, err = copierx.Trans[types.BorrowerPropertyTransactionCreateRequest, borrower.BorrowerPropertyTransactionCreateRequestData](req)
	if err != nil {
		return nil, err
	}

	rpcResp, err := l.svcCtx.BorrowerPropertyTransactionRpcClient.Create(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	resp = &types.BorrowerPropertyTransactionCreateResponse{}
	resp, err = copierx.Trans[borrower.BorrowerPropertyTransactionEntity, types.BorrowerPropertyTransactionCreateResponse](rpcResp.Data)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
