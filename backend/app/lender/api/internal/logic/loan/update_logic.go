package loan

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
)

type UpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// update
func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLogic) Update(req *types.LoanUpdateRequest) (resp *types.LoanUpdateResponse, err error) {
	rpcReq := &lender.LoanUpdateRequest{
		Data: &lender.LoanUpdateRequestData{},
	}
	resp = &types.LoanUpdateResponse{}
	rpcReq.Data, err = copierx.Trans[types.LoanUpdateRequest, lender.LoanUpdateRequestData](req)
	if err != nil {
		return nil, err
	}

	rpcResp, err := l.svcCtx.LoanRpcClient.UpdateLoan(l.ctx, rpcReq)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	resp, err = copierx.Trans[lender.LoanUpdateResponseData, types.LoanUpdateResponse](rpcResp.Data)
	if err != nil {
		return nil, err
	}

	return resp, err
}
