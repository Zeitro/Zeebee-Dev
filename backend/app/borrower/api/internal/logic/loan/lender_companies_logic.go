package loan

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"

	"github.com/zeromicro/go-zero/core/logx"
)

type LenderCompaniesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// lender companies
func NewLenderCompaniesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LenderCompaniesLogic {
	return &LenderCompaniesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LenderCompaniesLogic) LenderCompanies(req *types.EmptyRequest) (resp *types.LoanLenderCompanyResponse, err error) {
	rpcReq := &borrower.LoanLenderCompaniesRequest{
		Data: &borrower.LoanLenderCompaniesRequestData{},
	}

	rpcResp, err := l.svcCtx.LoanRpcClient.BorrowerLoanLenderCompanies(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}
	resp = &types.LoanLenderCompanyResponse{
		LenderCompanies: rpcResp.Data,
	}

	return resp, nil
}
