package borrower_property

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
)

type PropertyPageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// page
func NewPropertyPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PropertyPageLogic {
	return &PropertyPageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PropertyPageLogic) PropertyPage(req *types.PropertyPageRequest) (resp *types.PropertyPageResponse, err error) {
	rpcReq := &borrower.PropertyPageRequest{
		Data: &borrower.PropertyPageRequestData{
			Filter: &borrower.PropertyPageRequestFilter{},
			Sorter: &borrower.PropertyPageRequestSorter{},
		},
	}
	resp = &types.PropertyPageResponse{}
	rpcResp := &borrower.PropertyPageResponse{
		Data: &borrower.PropertyPageResponseData{},
	}
	rpcReq.Data, err = copierx.Trans[types.PropertyPageRequest, borrower.PropertyPageRequestData](req)
	if err != nil {
		return nil, err
	}

	rpcResp, err = l.svcCtx.BorrowerPropertyRpcClient.Page(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}
	resp, err = copierx.Trans[borrower.PropertyPageResponseData, types.PropertyPageResponse](rpcResp.Data)

	if err != nil {
		return nil, err
	}
	return resp, err
}
