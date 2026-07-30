package borrower_property_transaction

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/proto"
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
	rpcReq := &proto.PathIdRequest{
		Data: &proto.PathIdRequestData{},
	}

	rpcReq.Data, err = copierx.Trans[types.BorrowerPropertyTransactionDeleteRequest, proto.PathIdRequestData](req)
	if err != nil {
		return nil, err
	}

	_, err = l.svcCtx.BorrowerPropertyTransactionRpcClient.Delete(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	return &types.EmptyResponse{}, nil
}
