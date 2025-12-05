package organization_user

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrganizationUserPageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// page
func NewOrganizationUserPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationUserPageLogic {
	return &OrganizationUserPageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationUserPageLogic) OrganizationUserPage(req *types.OrganizationUserPageRequest) (resp *types.OrganizationUserPageResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
