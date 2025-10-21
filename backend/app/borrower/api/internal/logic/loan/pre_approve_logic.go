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

type PreApproveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// pre_approve
func NewPreApproveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PreApproveLogic {
	return &PreApproveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PreApproveLogic) PreApprove(req *types.LoanPreApproveRequest) (resp *types.LoanPreApproveResponse, err error) {
	rpcReq := &borrower.LoanPreApproveRequest{
		Data: &borrower.LoanPreApproveRequestData{},
	}
	resp = &types.LoanPreApproveResponse{}
	rpcResp := &borrower.LoanPreApproveResponse{}

	err = copier.CopyWithOption(&rpcReq.Data, req, copier.Option{
		DeepCopy:      true,
		CaseSensitive: true,
	})
	if err != nil {
		return nil, errorx.GetFailedTransReqErr(err)
	}
	rpcResp, err = l.svcCtx.LoanRpcClient.PreApproveLoan(l.ctx, rpcReq)
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
