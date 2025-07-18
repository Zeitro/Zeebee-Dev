package user

import (
	"context"

	"github.com/jinzhu/copier"

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
	// 请求数据转换
	err = copier.Copy(rpcReq, req)
	if err != nil {
		return nil, err
	}
	// 调用rpc服务
	rpcResp, err := l.svcCtx.RpcClients.UserRpcClient.Login(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 响应数据转换
	err = copier.Copy(resp, rpcResp)

	return resp, err
}
