package organization_user

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrganizationUserCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// create
func NewOrganizationUserCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationUserCreateLogic {
	return &OrganizationUserCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationUserCreateLogic) OrganizationUserCreate(req *types.OrganizationUserCreateRequest) (resp *types.OrganizationUserEntity, err error) {
	// todo: add your logic here and delete this line

	return
}
