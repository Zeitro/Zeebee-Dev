package borrower_property

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
)

type PropertyDetailSummaryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// detail
func NewPropertyDetailSummaryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PropertyDetailSummaryLogic {
	return &PropertyDetailSummaryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PropertyDetailSummaryLogic) PropertyDetailSummary(req *types.PropertyDetailRequest) (resp *types.PropertyDetailSummaryResponse, err error) {
	rpcReq := &borrower.PropertyDetailRequest{
		Data: &borrower.PropertyDetailRequestData{
			Id: req.Id,
		},
	}

	rpcResp, err := l.svcCtx.BorrowerPropertyRpcClient.DetailSummary(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	resp = &types.PropertyDetailSummaryResponse{}
	resp, err = copierx.Trans[borrower.PropertyDetailSummaryResponseData, types.PropertyDetailSummaryResponse](rpcResp.Data)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
