// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package organization_user

import (
	"context"

	"zeebee-backend/app/lender/api/internal/types/internal/svc"
	"zeebee-backend/app/lender/api/internal/types/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrganizationUserDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// delete
func NewOrganizationUserDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationUserDeleteLogic {
	return &OrganizationUserDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationUserDeleteLogic) OrganizationUserDelete(req *types.PathIdRequest) (resp *types.EmptyResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
