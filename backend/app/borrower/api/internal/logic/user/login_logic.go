package user

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
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
	resp = &types.UserLoginResponse{}
	// call rpc
	rpcResp, err := l.svcCtx.UserRpcClient.Login(l.ctx, &borrower.UserLoginRequest{
		Data: &borrower.UserLoginRequestData{
			Email:    req.Email,
			Password: req.Password,
		},
	})
	if err != nil {
		return nil, err
	}

	// trans response
	resp, err = copierx.Trans[borrower.UserLoginResponseData, types.UserLoginResponse](rpcResp.Data)
	if err != nil {
		return nil, err
	}
	return resp, err
}
