package loan_applicant

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type BorrowerNamesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// borrower names
func NewBorrowerNamesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BorrowerNamesLogic {
	return &BorrowerNamesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BorrowerNamesLogic) BorrowerNames(req *types.EmptyRequest) (resp *types.BorrowerNamesResponse, err error) {
	// get borrower id
	borrowerId, err := utils.GetAccountIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	rpcResp, err := l.svcCtx.LoanApplicantRpcClient.BorrowerNames(l.ctx, &borrower.BorrowerNamesRequest{
		Data: &borrower.BorrowerNamesRequestData{
			BorrowerId: borrowerId,
		},
	})
	if err != nil {
		return nil, err
	}
	resp = &types.BorrowerNamesResponse{
		BorrowerNames: rpcResp.Data,
	}

	return resp, nil
}
