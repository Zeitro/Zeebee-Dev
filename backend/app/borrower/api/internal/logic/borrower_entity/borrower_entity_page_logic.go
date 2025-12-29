package borrower_entity

import (
	"context"

	"github.com/jinzhu/copier"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
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
	rpcReq := &borrower.BorrowerEntityPageRequest{
		Data: &borrower.BorrowerEntityPageRequestData{
			Filter: &borrower.BorrowerEntityPageRequestFilter{},
			Sorter: &borrower.BorrowerEntityPageRequestSorter{},
		},
	}
	resp = &types.BorrowerEntityPageResponse{}
	rpcResp := &borrower.BorrowerEntityPageResponse{
		Data: &borrower.BorrowerEntityPageResponseData{},
	}

	err = copier.CopyWithOption(&rpcReq.Data, req, copier.Option{DeepCopy: true})
	if err != nil {
		return nil, errorx.GetFailedTransReqErr(err)
	}
	rpcResp, err = l.svcCtx.BorrowerEntityRpcClient.Page(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}
	err = copier.Copy(resp, &rpcResp.Data)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	return resp, err
}
