// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package role

import (
	"context"

	"zeebee-backend/app/lender/api/internal/types/internal/svc"
	"zeebee-backend/app/lender/api/internal/types/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleUpdateStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// update status
func NewRoleUpdateStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleUpdateStatusLogic {
	return &RoleUpdateStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RoleUpdateStatusLogic) RoleUpdateStatus(req *types.RoleUpdateStatusRequest) (resp *types.EmptyResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
