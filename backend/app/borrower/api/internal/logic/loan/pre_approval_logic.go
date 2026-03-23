package loan

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
)

type PreApprovalLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// pre_approve
func NewPreApprovalLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PreApprovalLogic {
	return &PreApprovalLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PreApprovalLogic) PreApproval(req *types.LoanPreApprovalRequest) (resp *types.LoanPreApprovalResponse, err error) {
	rpcReq := &borrower.LoanPreApprovalRequest{
		Data: &borrower.LoanPreApprovalRequestData{},
	}
	resp = &types.LoanPreApprovalResponse{}
	rpcResp := &borrower.LoanPreApprovalResponse{}

	rpcReq.Data, err = copierx.Trans[types.LoanPreApprovalRequest, borrower.LoanPreApprovalRequestData](req)
	if err != nil {
		return nil, err
	}

	rpcResp, err = l.svcCtx.LoanRpcClient.PreApprovalLoan(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	resp, err = copierx.Trans[borrower.LoanPreApprovalResponse, types.LoanPreApprovalResponse](rpcResp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
