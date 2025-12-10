package user

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/app/admin/rpc/pb/admin"
)

type UserDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// delete
func NewUserDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserDeleteLogic {
	return &UserDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserDeleteLogic) UserDelete(req *types.PathIdRequest) (resp *types.EmptyResponse, err error) {
	resp = &types.EmptyResponse{}
	rpcReq := &admin.UserDeleteRequest{
		Data: &admin.UserDeleteRequestData{},
	}
	rpcReq.Data.Id = req.Id

	_, err = l.svcCtx.UserRpcClient.UserDelete(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	return resp, err
}
