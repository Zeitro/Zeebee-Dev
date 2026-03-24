package borrower_entity

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
)

type BorrowerEntityUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// update
func NewBorrowerEntityUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BorrowerEntityUpdateLogic {
	return &BorrowerEntityUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BorrowerEntityUpdateLogic) BorrowerEntityUpdate(req *types.BorrowerEntityUpdateRequest) (resp *types.EmptyResponse, err error) {
	rpcReq := &borrower.BorrowerEntityUpdateRequest{
		Data: &borrower.BorrowerEntityEditEntity{},
	}
	rpcReq.Data, err = copierx.Trans[types.BorrowerEntityEditRequest, borrower.BorrowerEntityEditEntity](&req.BorrowerEntityEditRequest)
	if err != nil {
		return nil, err
	}

	_, err = l.svcCtx.BorrowerEntityRpcClient.Update(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	return &types.EmptyResponse{}, nil
}
