package user

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChangePasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// change password
func NewChangePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangePasswordLogic {
	return &ChangePasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChangePasswordLogic) ChangePassword(req *types.UserChangePasswordRequest) (resp *types.EmptyResponse, err error) {
	rpcReq := &borrower.ChangePasswordRequest{
		Data: &borrower.ChangePasswordRequestData{},
	}
	rpcReq.Data, err = copierx.Trans[types.UserChangePasswordRequest, borrower.ChangePasswordRequestData](req)
	if err != nil {
		return nil, err
	}

	_, err = l.svcCtx.UserRpcClient.ChangePassword(l.ctx, rpcReq)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	resp = &types.EmptyResponse{}
	return resp, nil
}
