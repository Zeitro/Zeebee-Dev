package organization_user

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrganizationUserUpdateStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// swaitch status
func NewOrganizationUserUpdateStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationUserUpdateStatusLogic {
	return &OrganizationUserUpdateStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationUserUpdateStatusLogic) OrganizationUserUpdateStatus(req *types.OrganizationUserSwitchStatusRequest) (resp *types.EmptyResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
