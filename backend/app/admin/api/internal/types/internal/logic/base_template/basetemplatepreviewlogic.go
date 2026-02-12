// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package base_template

import (
	"context"

	"zeebee-backend/app/admin/api/internal/types/internal/svc"
	"zeebee-backend/app/admin/api/internal/types/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BaseTemplatePreviewLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Preview
func NewBaseTemplatePreviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BaseTemplatePreviewLogic {
	return &BaseTemplatePreviewLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BaseTemplatePreviewLogic) BaseTemplatePreview(req *types.BaseTemplatePreviewRequest) (resp *types.BaseTemplatePreviewResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
