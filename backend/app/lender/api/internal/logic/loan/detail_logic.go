package loan

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/proto"
)

type DetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// detail
func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DetailLogic) Detail(req *types.LoanDetailRequest) (resp *types.LoanDetailResponse, err error) {
	resp = &types.LoanDetailResponse{}

	rpcResp, err := l.svcCtx.LoanRpcClient.LoanDetail(l.ctx, &proto.PathIdRequest{
		Data: &proto.PathIdRequestData{
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
