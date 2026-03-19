// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package organization_user

import (
	"context"

	"zeebee-backend/app/lender/api/internal/types/internal/svc"
	"zeebee-backend/app/lender/api/internal/types/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrganizationUserUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// update
func NewOrganizationUserUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationUserUpdateLogic {
	return &OrganizationUserUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationUserUpdateLogic) OrganizationUserUpdate(req *types.OrganizationUserUpdateRequest) (resp *types.OrganizationUserEntity, err error) {
	// todo: add your logic here and delete this line

	return
}
