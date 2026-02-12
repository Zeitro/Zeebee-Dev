package loan

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/app/admin/rpc/pb/admin"
)

type LenderCompaniesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// lender companies list
func NewLenderCompaniesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LenderCompaniesLogic {
	return &LenderCompaniesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LenderCompaniesLogic) LenderCompanies(req *types.EmptyRequest) (resp *types.LoanLenderCompanyResponse, err error) {

	rpcResp, err := l.svcCtx.LoanRpcClient.LoanLenderCompanies(l.ctx, &admin.EmptyRequest{})
	if err != nil {
		return nil, err
	}
	resp = &types.LoanLenderCompanyResponse{
		LenderCompanies: rpcResp.Data,
	}

	return resp, nil
}
