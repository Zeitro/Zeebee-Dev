package organization

import (
	"context"

	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrganizationProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// profile
func NewOrganizationProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationProfileLogic {
	return &OrganizationProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationProfileLogic) OrganizationProfile(req *types.EmptyRequest) (resp *types.OrganizationEntity, err error) {
	// todo: add your logic here and delete this line

	return
}
