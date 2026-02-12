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

type LoanDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// detail
func NewLoanDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoanDetailLogic {
	return &LoanDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoanDetailLogic) LoanDetail(req *types.PathIdRequest) (resp *types.LoanDetailResponse, err error) {
	rpcReq := &admin.IdRequest{Data: req.Id}

	rpcResp, err := l.svcCtx.LoanRpcClient.LoanDetail(l.ctx, rpcReq)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	resp, err = copierx.Trans[admin.LoanEntity, types.LoanDetailResponse](rpcResp.Data)

	return resp, err
}
