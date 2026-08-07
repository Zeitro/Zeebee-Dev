package borrower_entity

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/proto"
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
	rpcReq := &proto.PathIdRequest{
		Data: &proto.PathIdRequestData{
			Id: req.Id,
		},
	}

	rpcResp, err := l.svcCtx.BorrowerEntityRpcClient.Detail(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	resp = &types.BorrowerEntityDetailResponse{}
	entity, err := copierx.Trans[proto.BorrowerEntityEntity, types.BorrowerEntityEntity](rpcResp.Data)
	if err != nil {
		return nil, err
	}
	resp.BorrowerEntityEntity = *entity

	return resp, nil
}
