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

type LoanPageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// page
func NewLoanPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoanPageLogic {
	return &LoanPageLogic{
		Logger: svcCtx.Logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoanPageLogic) LoanPage(req *types.LoanPageRequest) (resp *types.LoanPageResponse, err error) {
	rpcReq := &admin.LoanPageRequest{}
	rpcReq.Data, err = copierx.Trans[types.LoanPageRequest, admin.LoanPageRequestData](req)

	if err != nil {
		return nil, err
	}

	rpcResp, err := l.svcCtx.LoanRpcClient.LoanPage(l.ctx, rpcReq)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	resp, err = copierx.Trans[admin.LoanPageResponseData, types.LoanPageResponse](rpcResp.Data)

	return resp, err
}
