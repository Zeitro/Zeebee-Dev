package user

import (
	"context"

	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"
	"zeebee-backend/app/lender/rpc/pb/lender"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendVerifyEmailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// send verify email
func NewSendVerifyEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendVerifyEmailLogic {
	return &SendVerifyEmailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendVerifyEmailLogic) SendVerifyEmail(req *types.UserSendVerifyEmailRequest) (resp *types.EmptyResponse, err error) {
	resp = &types.EmptyResponse{}

	_, err = l.svcCtx.UserRpcClient.SendVerifyEmail(l.ctx, &lender.SendVerifyEmailRequest{
		Data: &lender.SendVerifyEmailRequestData{
			Email: req.Email,
		},
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}
