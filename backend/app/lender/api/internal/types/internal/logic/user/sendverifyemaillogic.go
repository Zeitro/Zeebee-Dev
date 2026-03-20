// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"context"

	"zeebee-backend/app/lender/api/internal/types/internal/svc"
	"zeebee-backend/app/lender/api/internal/types/internal/types"

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
	// todo: add your logic here and delete this line

	return
}
