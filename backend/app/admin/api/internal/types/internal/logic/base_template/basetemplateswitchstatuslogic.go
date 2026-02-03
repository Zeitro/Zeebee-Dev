// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package base_template

import (
	"context"

	"zeebee-backend/app/admin/api/internal/types/internal/svc"
	"zeebee-backend/app/admin/api/internal/types/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BaseTemplateSwitchStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Switch status
func NewBaseTemplateSwitchStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BaseTemplateSwitchStatusLogic {
	return &BaseTemplateSwitchStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BaseTemplateSwitchStatusLogic) BaseTemplateSwitchStatus(req *types.BaseTemplateSwitchStatusRequest) (resp *types.EmptyResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
