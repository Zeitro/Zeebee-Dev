package loan

import (
	"context"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoanReassignLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// reassign
func NewLoanReassignLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoanReassignLogic {
	return &LoanReassignLogic{
		Logger: svcCtx.Logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoanReassignLogic) LoanReassign(req *types.LoanReassignRequest) (resp *types.EmptyResponse, err error) {
	rpcReq := &admin.ReassignLoanRequest{}
	rpcReq.Data, err = copierx.Trans[types.LoanReassignRequest, admin.ReassignLoanRequestData](req)

	if err != nil {
		return nil, err
	}

	_, err = l.svcCtx.LoanRpcClient.ReassignLoan(l.ctx, rpcReq)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	return &types.EmptyResponse{}, nil
}
