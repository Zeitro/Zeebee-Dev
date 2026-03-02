// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package resource

import (
	"context"

	"zeebee-backend/app/admin/api/internal/types/internal/svc"
	"zeebee-backend/app/admin/api/internal/types/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ResourceUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// update
func NewResourceUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResourceUpdateLogic {
	return &ResourceUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ResourceUpdateLogic) ResourceUpdate(req *types.ResourceUpdateRequest) (resp *types.ResourceEntity, err error) {
	// todo: add your logic here and delete this line

	return
}
