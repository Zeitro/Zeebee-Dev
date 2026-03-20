// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package base_template

import (
	"context"

	"zeebee-backend/app/admin/api/internal/types/internal/svc"
	"zeebee-backend/app/admin/api/internal/types/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BaseTemplateResetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Reset
func NewBaseTemplateResetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BaseTemplateResetLogic {
	return &BaseTemplateResetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BaseTemplateResetLogic) BaseTemplateReset(req *types.PathIdRequest) (resp *types.BaseTemplateDetailResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
