package borrower_property_transaction

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
)

type BorrowerPropertyTransactionUpdateChecklistLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// update transaction checklist
func NewBorrowerPropertyTransactionUpdateChecklistLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BorrowerPropertyTransactionUpdateChecklistLogic {
	return &BorrowerPropertyTransactionUpdateChecklistLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BorrowerPropertyTransactionUpdateChecklistLogic) BorrowerPropertyTransactionUpdateChecklist(req *types.BorrowerPropertyTransactionChecklistRequest) (resp *types.EmptyResponse, err error) {
	rpcReq := &borrower.BorrowerPropertyTransactionChecklistRequest{
		Data: &borrower.BorrowerPropertyTransactionChecklistRequestData{},
	}

	rpcReq.Data, err = copierx.Trans[types.BorrowerPropertyTransactionChecklistRequest, borrower.BorrowerPropertyTransactionChecklistRequestData](req)
	if err != nil {
		return nil, err
	}

	_, err = l.svcCtx.BorrowerPropertyTransactionRpcClient.UpdateChecklist(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	return &types.EmptyResponse{}, nil
}
