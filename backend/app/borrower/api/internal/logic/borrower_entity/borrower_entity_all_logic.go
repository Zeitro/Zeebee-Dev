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

type BorrowerEntityAllLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// all
func NewBorrowerEntityAllLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BorrowerEntityAllLogic {
	return &BorrowerEntityAllLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BorrowerEntityAllLogic) BorrowerEntityAll(req *types.EmptyRequest) (resp *types.BorrowerEntityAllResponse, err error) {
	rpcResp, err := l.svcCtx.BorrowerEntityRpcClient.All(l.ctx, &borrower.BorrowerEntityAllRequest{
		Data: &borrower.BorrowerEntityAllRequestData{},
	})
	if err != nil {
		return nil, err
	}

	resp = &types.BorrowerEntityAllResponse{}
	resp.Lists = make([]*types.BorrowerEntityPageEntry, 0, len(rpcResp.Data))

	err = copier.CopyWithOption(&resp.Lists, rpcResp.Data, copier.Option{
		DeepCopy:      true,
		CaseSensitive: true,
	})
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	return resp, nil
}
