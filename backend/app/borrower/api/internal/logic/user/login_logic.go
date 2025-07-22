package user

import (
	"context"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// login user
func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.UserLoginRequest) (resp *types.UserLoginResponse, err error) {
	rpcReq := &borrower.UserLoginRequest{}
	resp = &types.UserLoginResponse{}
	// trans req
	err = copier.Copy(rpcReq, req)
	if err != nil {
		return nil, errors.Wrap(err, "trans request error")
	}
	// call rpc
	rpcResp, err := l.svcCtx.UserRpcClient.Login(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// trans response
	err = copier.Copy(resp, rpcResp)
	if err != nil {
		return nil, errors.Wrap(err, "trans response error")
	}

	return resp, err
}
