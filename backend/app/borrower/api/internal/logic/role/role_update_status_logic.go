package role

import (
	"context"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/utils"
)

type RoleUpdateStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// update status
func NewRoleUpdateStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleUpdateStatusLogic {
	return &RoleUpdateStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RoleUpdateStatusLogic) RoleUpdateStatus(req *types.RoleUpdateStatusRequest) (resp *types.EmptyResponse, err error) {
	resp = &types.EmptyResponse{}
	rpcReq := &borrower.RoleUpdateStatusRequest{
		Data: &borrower.RoleUpdateStatusRequestData{},
	}
	err = copier.CopyWithOption(rpcReq.Data, req, copier.Option{DeepCopy: true})
	if err != nil {
		return nil, errorx.GetFailedTransReqErr(err)
	}

	claim, err := utils.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	rpcReq.Data.UserId = claim.UserId
	rpcReq.Data.OrganizationId = claim.OrganizationId

	_, err = l.svcCtx.RoleRpcClient.UpdateStatus(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	return resp, err
}
