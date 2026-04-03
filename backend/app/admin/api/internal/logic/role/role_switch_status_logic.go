package role

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
)

type RoleSwitchStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// update status
func NewRoleSwitchStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleSwitchStatusLogic {
	return &RoleSwitchStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RoleSwitchStatusLogic) RoleSwitchStatus(req *types.RoleSwitchStatusRequest) (resp *types.EmptyResponse, err error) {
	resp = &types.EmptyResponse{}
	rpcReq := &admin.RoleSwitchStatusRequest{
		Data: &admin.RoleSwitchStatusRequestData{},
	}

	rpcReq.Data, err = copierx.Trans[types.RoleSwitchStatusRequest, admin.RoleSwitchStatusRequestData](req)
	if err != nil {
		return nil, err
	}

	_, err = l.svcCtx.RoleRpcClient.SwitchStatus(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	return resp, err
}
