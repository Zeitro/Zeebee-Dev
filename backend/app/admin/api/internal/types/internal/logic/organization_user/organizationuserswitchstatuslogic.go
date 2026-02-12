// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package organization_user

import (
	"context"

	"zeebee-backend/app/admin/api/internal/types/internal/svc"
	"zeebee-backend/app/admin/api/internal/types/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrganizationUserSwitchStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// switch status
func NewOrganizationUserSwitchStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationUserSwitchStatusLogic {
	return &OrganizationUserSwitchStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationUserSwitchStatusLogic) OrganizationUserSwitchStatus(req *types.OrganizationUserSwitchStatusRequest) (resp *types.EmptyResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
