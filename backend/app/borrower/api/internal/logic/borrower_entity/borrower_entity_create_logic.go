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

type BorrowerEntityCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// create
func NewBorrowerEntityCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BorrowerEntityCreateLogic {
	return &BorrowerEntityCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BorrowerEntityCreateLogic) BorrowerEntityCreate(req *types.BorrowerEntityCreateRequest) (resp *types.BorrowerEntityCreateResponse, err error) {
	rpcReq := &borrower.BorrowerEntityCreateRequest{
		Data: &borrower.BorrowerEntityEditEntity{},
	}
	rpcReq.Data, err = copierx.Trans[types.BorrowerEntityCreateRequest, borrower.BorrowerEntityEditEntity](req)
	if err != nil {
		return nil, err
	}

	rpcResp, err := l.svcCtx.BorrowerEntityRpcClient.Create(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	resp = &types.BorrowerEntityCreateResponse{}
	entity, err := copierx.Trans[proto.BorrowerEntityEntity, types.BorrowerEntityEntity](rpcResp.Data)
	if err != nil {
		return nil, err
	}
	resp.BorrowerEntityEntity = *entity

	return resp, nil
}
