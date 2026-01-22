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

	err = copier.CopyWithOption(&rpcReq.Data, req, copier.Option{
		DeepCopy:      true,
		CaseSensitive: true,
	})
	if err != nil {
		return nil, errorx.GetFailedTransReqErr(err)
	}
	rpcResp, err = l.svcCtx.LoanRpcClient.PreApprovalLoan(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	err = copier.CopyWithOption(resp, rpcResp.Data, copier.Option{
		DeepCopy:      true,
		CaseSensitive: true,
	})
	if err != nil {
		return nil, errorx.GetFailedTransReqErr(err)
	}

	return resp, nil
}
