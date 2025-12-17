package organization

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrganizationInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// info
func NewOrganizationInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationInfoLogic {
	return &OrganizationInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationInfoLogic) OrganizationInfo(req *types.EmptyRequest) (resp *types.OrganizationEntity, err error) {
	// todo: add your logic here and delete this line

	return
}
