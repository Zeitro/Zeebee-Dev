package user

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"
	"zeebee-backend/app/lender/rpc/pb/lender"
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

func (l *ResetPasswordLogic) ResetPassword(req *types.UserResetPasswordRequest) (resp *types.UserJwtResponse, err error) {
	rpcResp, err := l.svcCtx.UserRpcClient.ResetPassword(l.ctx, &lender.ResetPasswordRequest{
		Data: &lender.ResetPasswordRequestData{
			Email:    req.Email,
			Code:     req.Code,
			Password: req.Password,
		},
	})
	if err != nil {
		return nil, err
	}
	// trans response data
	resp, err = copierx.Trans[lender.UserJwtResponseData, types.UserJwtResponse](rpcResp.Data)
	return resp, err
}
