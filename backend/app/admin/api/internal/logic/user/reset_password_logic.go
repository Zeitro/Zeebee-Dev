package user

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
)

type ResetPasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// reset password
func NewResetPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetPasswordLogic {
	return &ResetPasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ResetPasswordLogic) ResetPassword(req *types.UserResetPasswordRequest) (resp *types.UserLoginResponse, err error) {
	resp = &types.UserLoginResponse{}
	rpcReq := &admin.ResetPasswordRequest{
		Data: &admin.ResetPasswordRequestData{
			Password: req.Password,
			Code:     req.Code,
			Email:    req.Email,
		},
	}
	// call rpc
	rpcResp, err := l.svcCtx.UserRpcClient.ResetPassword(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// trans response
	resp, err = copierx.Trans[admin.UserLoginResponseData, types.UserLoginResponse](rpcResp.Data)
	return resp, err
}
