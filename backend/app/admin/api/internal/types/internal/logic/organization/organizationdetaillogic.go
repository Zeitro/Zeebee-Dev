// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package organization

import (
	"context"

	"zeebee-backend/app/admin/api/internal/types/internal/svc"
	"zeebee-backend/app/admin/api/internal/types/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrganizationDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// detail
func NewOrganizationDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationDetailLogic {
	return &OrganizationDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationDetailLogic) OrganizationDetail(req *types.PathIdRequest) (resp *types.OrganizationEntity, err error) {
	// todo: add your logic here and delete this line

	return
}
