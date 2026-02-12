// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package organization

import (
	"context"

	"zeebee-backend/app/admin/api/internal/types/internal/svc"
	"zeebee-backend/app/admin/api/internal/types/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrganizationDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// delete
func NewOrganizationDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationDeleteLogic {
	return &OrganizationDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationDeleteLogic) OrganizationDelete(req *types.PathIdRequest) (resp *types.EmptyResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
