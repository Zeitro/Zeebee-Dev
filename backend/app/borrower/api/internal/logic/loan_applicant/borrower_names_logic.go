package loan_applicant

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
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
	rpcResp, err := l.svcCtx.LoanApplicantRpcClient.BorrowerNames(l.ctx, &borrower.BorrowerNamesRequest{
		Data: &borrower.BorrowerNamesRequestData{},
	})
	if err != nil {
		return nil, err
	}
	resp = &types.BorrowerNamesResponse{
		BorrowerNames: rpcResp.Data,
	}

	return resp, nil
}
