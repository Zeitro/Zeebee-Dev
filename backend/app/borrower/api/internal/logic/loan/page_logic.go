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

	rpcReq := &borrower.PageLoanRequest{
		Data: &borrower.PageLoanRequestData{
			Filter: &borrower.PageLoanRequestFilter{},
		},
	}
	resp = &types.LoanPageResponse{}
	rpcResp := &borrower.PageLoanResponse{
		Data: &borrower.PageLoanResponseData{},
	}

	err = copier.Copy(&rpcReq.Data, req)
	if err != nil {
		return nil, errorx.GetFailedTransReqErr(err)
	}
	// get borrower id
	rpcReq.Data.Filter.BorrowerId, err = utils.GetAccountIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	rpcResp, err = l.svcCtx.LoanRpcClient.GetBorrowerLoanPage(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}
	err = copier.Copy(resp, &rpcResp.Data)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	return resp, err
}
