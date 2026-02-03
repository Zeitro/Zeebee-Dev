// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package common

import (
	"context"

	"zeebee-backend/app/admin/api/internal/types/internal/svc"
	"zeebee-backend/app/admin/api/internal/types/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SelectOptionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// select options
func NewSelectOptionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SelectOptionsLogic {
	return &SelectOptionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SelectOptionsLogic) SelectOptions(req *types.GetSelectOptionRequest) (resp *types.GetSelectOptionResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
