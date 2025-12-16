package user

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type EmailAuthLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// email auth
func NewEmailAuthLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EmailAuthLogic {
	return &EmailAuthLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EmailAuthLogic) EmailAuth(req *types.UserEmailAuthRequest) (resp *types.EmptyResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
