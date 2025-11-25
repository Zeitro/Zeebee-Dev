package borrower_property

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type PropertyDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// delete
func NewPropertyDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PropertyDeleteLogic {
	return &PropertyDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PropertyDeleteLogic) PropertyDelete(req *types.PathIdRequest) (resp *types.EmptyResponse, err error) {
	rpcReq := &borrower.PropertyDeleteRequest{
		Data: &borrower.PropertyDeleteRequestData{
			Id: req.Id,
		},
	}

	// get borrower id
	claim, err := utils.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	rpcReq.Data.UserId = claim.UserId
	rpcReq.Data.OrganizationId = claim.OrganizationId

	_, err = l.svcCtx.BorrowerPropertyRpcClient.Delete(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	resp = &types.EmptyResponse{}
	return resp, nil
}
