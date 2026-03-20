package loan

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/copierx"
)

type PageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// page
func NewPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PageLogic {
	return &PageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PageLogic) Page(req *types.LoanPageRequest) (resp *types.LoanPageResponse, err error) {
	rpcReq := &lender.LoanPageRequest{
		Data: &lender.LoanPageRequestData{
			Filter: &lender.LoanPageRequestFilter{},
		},
	}
	resp = &types.LoanPageResponse{}
	rpcReq.Data, err = copierx.Trans[types.LoanPageRequest, lender.LoanPageRequestData](req)
	if err != nil {
		return nil, err
	}

	rpcResp, err := l.svcCtx.LoanRpcClient.LoanPage(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	resp, err = copierx.Trans[lender.LoanPageResponseData, types.LoanPageResponse](rpcResp.Data)
	if err != nil {
		return nil, err
	}
	return resp, err
}
