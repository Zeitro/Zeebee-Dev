package borrower_property

import (
	"context"

	"github.com/jinzhu/copier"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/utils"

	"github.com/zeromicro/go-zero/core/logx"
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

	err = copier.CopyWithOption(&rpcReq.Data, req, copier.Option{
		DeepCopy:      true,
		IgnoreEmpty:   true,
		CaseSensitive: true,
	})
	if err != nil {
		return nil, errorx.GetFailedTransReqErr(err)
	}

	// get borrower id
	borrowerId, err := utils.GetAccountIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	rpcReq.Data.BorrowerId = borrowerId

	_, err = l.svcCtx.PropertyRpcClient.UpdateProperty(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	resp = &types.EmptyResponse{}
	return resp, nil
}
