// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package organization_user

import (
	"context"

	"zeebee-backend/app/lender/api/internal/types/internal/svc"
	"zeebee-backend/app/lender/api/internal/types/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrganizationUserAssignRolesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// assign roles
func NewOrganizationUserAssignRolesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationUserAssignRolesLogic {
	return &OrganizationUserAssignRolesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationUserAssignRolesLogic) OrganizationUserAssignRoles(req *types.OrganizationUserAssignRolesRequest) (resp *types.EmptyResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
