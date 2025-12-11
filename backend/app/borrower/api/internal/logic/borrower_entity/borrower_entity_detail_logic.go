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

type BorrowerEntityDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// detail
func NewBorrowerEntityDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BorrowerEntityDetailLogic {
	return &BorrowerEntityDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BorrowerEntityDetailLogic) BorrowerEntityDetail(req *types.PathIdRequest) (resp *types.BorrowerEntityDetailResponse, err error) {
	rpcReq := &borrower.BorrowerEntityDetailRequest{
		Data: &borrower.BorrowerEntityDetailRequestData{
			Id: req.Id,
		},
	}

	rpcResp, err := l.svcCtx.BorrowerEntityRpcClient.Detail(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	resp = &types.BorrowerEntityDetailResponse{}
	err = copier.CopyWithOption(&resp.BorrowerEntityEntry, &rpcResp.Data, copier.Option{DeepCopy: true})
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	return resp, nil
}
