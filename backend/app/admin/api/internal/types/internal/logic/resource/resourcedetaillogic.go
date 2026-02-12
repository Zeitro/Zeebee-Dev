// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package resource

import (
	"context"

	"zeebee-backend/app/admin/api/internal/types/internal/svc"
	"zeebee-backend/app/admin/api/internal/types/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ResourceDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// detail
func NewResourceDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResourceDetailLogic {
	return &ResourceDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ResourceDetailLogic) ResourceDetail(req *types.ResourceDetailRequest) (resp *types.ResourceEntity, err error) {
	// todo: add your logic here and delete this line

	return
}
