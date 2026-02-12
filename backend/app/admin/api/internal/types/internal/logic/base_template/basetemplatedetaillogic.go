// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package base_template

import (
	"context"

	"zeebee-backend/app/admin/api/internal/types/internal/svc"
	"zeebee-backend/app/admin/api/internal/types/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BaseTemplateDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Detail
func NewBaseTemplateDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BaseTemplateDetailLogic {
	return &BaseTemplateDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BaseTemplateDetailLogic) BaseTemplateDetail(req *types.PathIdRequest) (resp *types.BaseTemplateDetailResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
