package borrower_entity

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/domains/zeebee/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type BorrowerEntityDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// delete
func NewBorrowerEntityDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BorrowerEntityDeleteLogic {
	return &BorrowerEntityDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BorrowerEntityDeleteLogic) BorrowerEntityDelete(req *types.PathIdRequest) (resp *types.EmptyResponse, err error) {
	rpcReq := &proto.PathIdRequest{
		Data: &proto.PathIdRequestData{
			Id: req.Id,
		},
	}

	_, err = l.svcCtx.BorrowerEntityRpcClient.Delete(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	resp = &types.EmptyResponse{}
	return resp, nil
}
