package user

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/proto"
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

	rpcResp, err := l.svcCtx.UserRpcClient.ResetPassword(l.ctx, &proto.ResetPasswordRequest{
		Data: &proto.ResetPasswordRequestData{
			Email:    req.Email,
			Code:     req.Code,
			Password: req.Password,
		},
	})
	if err != nil {
		return nil, err
	}
	// trans response data
	resp, err = copierx.Trans[proto.UserLoginResponseData, types.UserLoginResponse](rpcResp.Data)

	if err != nil {
		return nil, err
	}
	return resp, nil
}
