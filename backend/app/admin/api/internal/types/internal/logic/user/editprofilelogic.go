// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"context"

	"zeebee-backend/app/admin/api/internal/types/internal/svc"
	"zeebee-backend/app/admin/api/internal/types/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type EditProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// edit profile
func NewEditProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EditProfileLogic {
	return &EditProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EditProfileLogic) EditProfile(req *types.UserEditProfileRequest) (resp *types.UserInfoResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
