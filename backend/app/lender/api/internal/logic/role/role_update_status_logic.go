package role

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
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
	rpcReq := &lender.RoleUpdateStatusRequest{
		Data: &lender.RoleUpdateStatusRequestData{},
	}
	data, err := copierx.Trans[types.RoleUpdateStatusRequest, lender.RoleUpdateStatusRequestData](req)
	if err != nil {
		return nil, errorx.GetFailedTransReqErr(err)
	}
	rpcReq.Data = data

	_, err = l.svcCtx.RoleRpcClient.UpdateStatus(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	return &types.EmptyResponse{}, nil
}
