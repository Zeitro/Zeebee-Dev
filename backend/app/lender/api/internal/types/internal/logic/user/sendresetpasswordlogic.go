// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"context"

	"zeebee-backend/app/lender/api/internal/types/internal/svc"
	"zeebee-backend/app/lender/api/internal/types/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
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
	// todo: add your logic here and delete this line

	return
}
