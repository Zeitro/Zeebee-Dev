package user

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/utils"
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

	claim, err := utils.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	rpcReq.Data.UserId = claim.UserId
	rpcReq.Data.OrganizationId = claim.OrganizationId

	_, err = l.svcCtx.UserRpcClient.UserDelete(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	return resp, err
}
