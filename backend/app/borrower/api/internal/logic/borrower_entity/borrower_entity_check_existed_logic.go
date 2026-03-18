package borrower_entity

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
)

type BorrowerEntityCheckExistedLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// existed
func NewBorrowerEntityCheckExistedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BorrowerEntityCheckExistedLogic {
	return &BorrowerEntityCheckExistedLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BorrowerEntityCheckExistedLogic) BorrowerEntityCheckExisted(req *types.BorrowerEntityCheckExistedRequest) (resp *types.BorrowerEntityCheckExistedResponse, err error) {
	rpcReq := &borrower.BorrowerEntityCheckExistedRequest{
		Data: &borrower.BorrowerEntityCheckExistedRequestData{
			Name: req.Name,
			Id:   req.Id,
		},
	}

	rpcResp, err := l.svcCtx.BorrowerEntityRpcClient.CheckExisted(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	resp = &types.BorrowerEntityCheckExistedResponse{}
	resp, err = copierx.Trans[borrower.BorrowerEntityCheckExistedResponseData, types.BorrowerEntityCheckExistedResponse](rpcResp.Data)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
