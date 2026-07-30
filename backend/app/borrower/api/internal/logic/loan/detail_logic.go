package loan

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/proto"
)

type DetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// get loan
func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DetailLogic) Detail(req *types.LoanDetailRequest) (resp *types.LoanDetailResponse, err error) {
	resp = &types.LoanDetailResponse{}

	rpcResp, err := l.svcCtx.LoanRpcClient.BorrowerLoanDetail(l.ctx, &borrower.LoanDetailRequest{
		Data: &borrower.LoanDetailRequestData{
			Id: req.Id,
		},
	})
	if err != nil {
		return nil, err
	}
	resp, err = copierx.Trans[proto.LoanEntity, types.LoanDetailResponse](rpcResp.Data)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
