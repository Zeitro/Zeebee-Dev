package borrower_property

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
)

type PropertyUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// update
func NewPropertyUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PropertyUpdateLogic {
	return &PropertyUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PropertyUpdateLogic) PropertyUpdate(req *types.PropertyUpdateRequest) (resp *types.EmptyResponse, err error) {
	rpcReq := &borrower.PropertyUpdateRequest{
		Data: &borrower.PropertyUpdateRequestData{},
	}

	rpcReq.Data, err = copierx.Trans[types.PropertyUpdateRequest, borrower.PropertyUpdateRequestData](req)
	if err != nil {
		return nil, err
	}

	_, err = l.svcCtx.BorrowerPropertyRpcClient.Update(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	resp = &types.EmptyResponse{}
	return resp, nil
}
