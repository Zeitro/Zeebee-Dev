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
	claim, err := utils.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	_, err = l.svcCtx.LoanRpcClient.DeleteBorrowerLoan(l.ctx, &borrower.LoanDeleteRequest{
		Data: &borrower.LoanDeleteRequestData{
			Id:             req.Id,
			UserId:         claim.UserId,
			OrganizationId: claim.OrganizationId,
		},
	})

	return
}
