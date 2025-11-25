package borrower_entity

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

type BorrowerEntityUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// update
func NewBorrowerEntityUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BorrowerEntityUpdateLogic {
	return &BorrowerEntityUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BorrowerEntityUpdateLogic) BorrowerEntityUpdate(req *types.BorrowerEntityUpdateRequest) (resp *types.EmptyResponse, err error) {
	rpcReq := &borrower.BorrowerEntityUpdateRequest{
		Data: &borrower.BorrowerEntityEditEntry{},
	}

	err = copier.CopyWithOption(&rpcReq.Data, &req.BorrowerEntityEditRequest, copier.Option{
		DeepCopy:      true,
		IgnoreEmpty:   true,
		CaseSensitive: true,
	})
	if err != nil {
		return nil, errorx.GetFailedTransReqErr(err)
	}

	// get borrower id
	claim, err := utils.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	rpcReq.Data.UserId = claim.UserId
	rpcReq.Data.OrganizationId = claim.OrganizationId

	_, err = l.svcCtx.BorrowerEntityRpcClient.Update(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	return &types.EmptyResponse{}, nil
}
