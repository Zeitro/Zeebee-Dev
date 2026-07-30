package borrower_entity

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/proto"
)

type BorrowerEntityPageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// page
func NewBorrowerEntityPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BorrowerEntityPageLogic {
	return &BorrowerEntityPageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BorrowerEntityPageLogic) BorrowerEntityPage(req *types.BorrowerEntityPageRequest) (resp *types.BorrowerEntityPageResponse, err error) {
	rpcReq := &proto.CommonPageRequest{
		Data: &proto.CommonPageRequestData{
			Filter: &proto.CommonPageFilter{},
			Sorter: &proto.CommonPageSorter{},
		},
	}
	resp = &types.BorrowerEntityPageResponse{}
	rpcResp := &borrower.BorrowerEntityPageResponse{
		Data: &borrower.BorrowerEntityPageResponseData{},
	}
	rpcReq.Data, err = copierx.Trans[types.BorrowerEntityPageRequest, proto.CommonPageRequestData](req)
	if err != nil {
		return nil, err
	}

	rpcResp, err = l.svcCtx.BorrowerEntityRpcClient.Page(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}
	resp, err = copierx.Trans[borrower.BorrowerEntityPageResponseData, types.BorrowerEntityPageResponse](rpcResp.Data)
	if err != nil {
		return nil, err
	}
	return resp, err
}
