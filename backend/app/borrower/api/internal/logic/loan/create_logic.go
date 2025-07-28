package loan

import (
	"context"

	"github.com/jinzhu/copier"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// create loan
func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLogic) Create(req *types.LoanCreateRequest) (resp *types.LoanCreateResponse, err error) {
	rpcReq := &borrower.CreateLoanRequest{}
	resp = &types.LoanCreateResponse{}
	rpcResp := &borrower.CreateLoanResponse{}

	err = copier.Copy(rpcReq, req)
	if err != nil {
		return nil, errorx.GetFailedTransReqErr(err)
	}

	rpcResp, err = l.svcCtx.LoanRpcClient.CreateBorrowerLoan(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}
	err = copier.Copy(resp, rpcResp)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	return resp, err
}
