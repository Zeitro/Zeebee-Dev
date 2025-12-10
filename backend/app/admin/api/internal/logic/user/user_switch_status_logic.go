package user

import (
	"context"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/errorx"
)

type UserSwitchStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// update status
func NewUserSwitchStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserSwitchStatusLogic {
	return &UserSwitchStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserSwitchStatusLogic) UserSwitchStatus(req *types.UserSwitchStatusRequest) (resp *types.EmptyResponse, err error) {
	resp = &types.EmptyResponse{}
	rpcReq := &admin.UserSwitchStatusRequest{
		Data: &admin.UserSwitchStatusRequestData{},
	}
	err = copier.CopyWithOption(rpcReq.Data, req, copier.Option{DeepCopy: true})
	if err != nil {
		return nil, errorx.GetFailedTransReqErr(err)
	}

	_, err = l.svcCtx.UserRpcClient.UserSwitchStatus(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	return resp, err
}
