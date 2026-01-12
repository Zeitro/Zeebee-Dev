package organization

import (
	"context"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrganizationSwitchStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// switch status
func NewOrganizationSwitchStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationSwitchStatusLogic {
	return &OrganizationSwitchStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationSwitchStatusLogic) OrganizationSwitchStatus(req *types.OrganizationSwitchStatusRequest) (resp *types.EmptyResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
