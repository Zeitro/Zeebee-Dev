package borrower_property

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
)

type PropertyCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// create
func NewPropertyCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PropertyCreateLogic {
	return &PropertyCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PropertyCreateLogic) PropertyCreate(req *types.PropertyCreateRequest) (resp *types.PropertyCreateResponse, err error) {
	rpcReq := &borrower.PropertyCreateRequest{
		Data: &borrower.PropertyCreateRequestData{},
	}
	rpcReq.Data, err = copierx.Trans[types.PropertyCreateRequest, borrower.PropertyCreateRequestData](req)
	if err != nil {
		return nil, err
	}

	rpcResp, err := l.svcCtx.BorrowerPropertyRpcClient.Create(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	resp = &types.PropertyCreateResponse{}
	resp, err = copierx.Trans[borrower.BorrowerPropertyEntity, types.PropertyCreateResponse](rpcResp.Data)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
