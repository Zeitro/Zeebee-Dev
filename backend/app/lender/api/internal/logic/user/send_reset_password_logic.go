package user

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"
	"zeebee-backend/domains/zeebee/proto"
)

type SendResetPasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// send reset password email
func NewSendResetPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendResetPasswordLogic {
	return &SendResetPasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendResetPasswordLogic) SendResetPassword(req *types.UserSendResetPasswordEmailRequest) (resp *types.EmptyResponse, err error) {
	resp = &types.EmptyResponse{}
	_, err = l.svcCtx.UserRpcClient.SendResetPasswordEmail(l.ctx, &proto.SendResetPasswordEmailRequest{
		Data: &proto.SendResetPasswordEmailRequestData{
			Email: req.Email,
		},
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}
