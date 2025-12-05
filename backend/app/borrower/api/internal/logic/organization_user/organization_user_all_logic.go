package organization_user

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrganizationUserAllLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// all
func NewOrganizationUserAllLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationUserAllLogic {
	return &OrganizationUserAllLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationUserAllLogic) OrganizationUserAll(req *types.EmptyRequest) (resp *types.OrganizationUserAllResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
