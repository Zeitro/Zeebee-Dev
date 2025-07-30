package loan

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteLoanLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// delete loan
func NewDeleteLoanLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLoanLogic {
	return &DeleteLoanLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteLoanLogic) DeleteLoan(req *types.LoanDeleteRequest) (resp *types.EmptyResponse, err error) {
	// get borrower id
	borrowerId, err := utils.GetAccountIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	_, err = l.svcCtx.LoanRpcClient.DeleteBorrowerLoan(l.ctx, &borrower.DeleteLoanRequest{
		Data: &borrower.DeleteLoanRequestData{
			Id:         req.Id,
			BorrowerId: borrowerId,
		},
	})

	return
}
