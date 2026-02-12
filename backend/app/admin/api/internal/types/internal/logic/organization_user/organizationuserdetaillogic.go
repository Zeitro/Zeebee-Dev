// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package organization_user

import (
	"context"

	"zeebee-backend/app/admin/api/internal/types/internal/svc"
	"zeebee-backend/app/admin/api/internal/types/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrganizationUserDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// detail
func NewOrganizationUserDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationUserDetailLogic {
	return &OrganizationUserDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationUserDetailLogic) OrganizationUserDetail(req *types.PathIdRequest) (resp *types.OrganizationUserDetailResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
