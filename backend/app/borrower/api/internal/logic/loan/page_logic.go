package loan

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
)

type PageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// loan page
func NewPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PageLogic {
	return &PageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PageLogic) Page(req *types.LoanPageRequest) (resp *types.LoanPageResponse, err error) {

	rpcReq := &borrower.LoanPageRequest{
		Data: &borrower.LoanPageRequestData{
			Filter: &borrower.LoanPageRequestFilter{},
		},
	}
	resp = &types.LoanPageResponse{}
	rpcResp := &borrower.LoanPageResponse{
		Data: &borrower.LoanPageResponseData{},
	}
	rpcReq.Data, err = copierx.Trans[types.LoanPageRequest, borrower.LoanPageRequestData](req)
	if err != nil {
		return nil, err
	}
	rpcResp, err = l.svcCtx.LoanRpcClient.BorrowerLoanPage(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}
	resp, err = copierx.Trans[borrower.LoanPageResponseData, types.LoanPageResponse](rpcResp.Data)
	if err != nil {
		return nil, err
	}
	return resp, err
}
