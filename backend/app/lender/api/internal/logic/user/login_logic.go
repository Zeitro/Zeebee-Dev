package user

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/copierx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// login
func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.UserLoginRequest) (resp *types.UserJwtResponse, err error) {
	// call rpc
	rpcResp, err := l.svcCtx.UserRpcClient.Login(l.ctx, &lender.UserLoginRequest{
		Data: &lender.UserLoginRequestData{
			Email:    req.Email,
			Password: req.Password,
		},
	})
	if err != nil {
		return nil, err
	}
	// trans response
	resp, err = copierx.Trans[lender.UserJwtResponseData, types.UserJwtResponse](rpcResp.Data)
	if err != nil {
		return nil, err
	}

	return resp, err
}
