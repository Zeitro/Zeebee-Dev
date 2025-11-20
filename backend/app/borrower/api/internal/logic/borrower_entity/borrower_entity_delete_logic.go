package borrower_entity

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type BorrowerEntityDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// delete
func NewBorrowerEntityDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BorrowerEntityDeleteLogic {
	return &BorrowerEntityDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BorrowerEntityDeleteLogic) BorrowerEntityDelete(req *types.PathIdRequest) (resp *types.EmptyResponse, err error) {
	rpcReq := &borrower.BorrowerEntityDeleteRequest{
		Data: &borrower.BorrowerEntityDeleteRequestData{
			Id: req.Id,
		},
	}

	// get borrower id
	borrowerId, err := utils.GetAccountIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	rpcReq.Data.BorrowerId = borrowerId

	_, err = l.svcCtx.BorrowerEntityRpcClient.Delete(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	resp = &types.EmptyResponse{}
	return resp, nil
}
