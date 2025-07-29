package loan

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
)

type DetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// get loan
func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DetailLogic) Detail(req *types.LoanDetailRequest) (resp *types.LoanDetailResponse, err error) {
	// get borrower id
	// borrowerId, err := utils.GetAccountIdFromCtx(l.ctx)
	// if err != nil {
	// 	return nil, err
	// }
	//
	// rpcResp, err := l.svcCtx.LoanRpcClient.GetBorrowerLoan(l.ctx, &borrower.GetLoanRequest{
	// 	Data: &borrower.GetLoanRequestData{
	// 		Id:         req.Id,
	// 		BorrowerId: borrowerId,
	// 	},
	// })
	// if err != nil {
	// 	return nil, err
	// }

	return resp, nil
}
