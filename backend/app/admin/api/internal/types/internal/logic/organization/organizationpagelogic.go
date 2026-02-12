// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package organization

import (
	"context"

	"zeebee-backend/app/admin/api/internal/types/internal/svc"
	"zeebee-backend/app/admin/api/internal/types/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrganizationPageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// page
func NewOrganizationPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationPageLogic {
	return &OrganizationPageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationPageLogic) OrganizationPage(req *types.OrganizationPageRequest) (resp *types.OrganizationPageResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
