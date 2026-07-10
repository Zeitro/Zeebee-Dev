package user

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/domains/zeebee/proto"

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

	_, err = l.svcCtx.UserRpcClient.SendVerifyEmail(l.ctx, &proto.SendVerifyEmailRequest{
		Data: &proto.SendVerifyEmailRequestData{
			Email: req.Email,
		},
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}
