// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package resource

import (
	"context"

	"zeebee-backend/app/admin/api/internal/types/internal/svc"
	"zeebee-backend/app/admin/api/internal/types/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ResourceTreeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// tree
func NewResourceTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResourceTreeLogic {
	return &ResourceTreeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ResourceTreeLogic) ResourceTree(req *types.ResourceTreeRequest) (resp *types.ResourceTreeResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
