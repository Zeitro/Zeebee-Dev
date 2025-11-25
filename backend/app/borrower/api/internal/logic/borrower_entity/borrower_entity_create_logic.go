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

	rpcResp, err := l.svcCtx.BorrowerEntityRpcClient.Create(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	resp = &types.BorrowerEntityCreateResponse{}
	err = copier.CopyWithOption(&resp.BorrowerEntityEntry, &rpcResp.Data, copier.Option{DeepCopy: true})
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	return resp, nil
}
