// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package base_template

import (
	"context"

	"zeebee-backend/app/admin/api/internal/types/internal/svc"
	"zeebee-backend/app/admin/api/internal/types/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BaseTemplateDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Delete
func NewBaseTemplateDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BaseTemplateDeleteLogic {
	return &BaseTemplateDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BaseTemplateDeleteLogic) BaseTemplateDelete(req *types.PathIdRequest) (resp *types.EmptyResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
