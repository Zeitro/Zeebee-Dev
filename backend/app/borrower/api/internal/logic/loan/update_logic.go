package loan

import (
	"context"

	"github.com/jinzhu/copier"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// update loan
func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLogic) Update(req *types.LoanUpdateRequest) (resp *types.EmptyResponse, err error) {
	rpcReq := &borrower.LoanUpdateRequest{
		Data: &borrower.LoanUpdateRequestData{},
	}
	resp = &types.EmptyResponse{}

	err = copier.Copy(&rpcReq.Data, req)
	if err != nil {
		return nil, errorx.GetFailedTransReqErr(err)
	}
	rpcReq.Data.Property.Address.State = req.Property.Address.State
	// get borrower id
	rpcReq.Data.BorrowerId, err = utils.GetAccountIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	_, err = l.svcCtx.LoanRpcClient.UpdateBorrowerLoan(l.ctx, rpcReq)

	return nil, err
}
