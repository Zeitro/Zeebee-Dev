package organization

import (
	"context"

	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrganizationUpdateProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// update profile
func NewOrganizationUpdateProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationUpdateProfileLogic {
	return &OrganizationUpdateProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationUpdateProfileLogic) OrganizationUpdateProfile(req *types.OrganizationUpdateProfileRequest) (resp *types.OrganizationEntity, err error) {
	// todo: add your logic here and delete this line

	return
}
