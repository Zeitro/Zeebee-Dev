package organization

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrganizationEmailAuthorizeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// email authorize
func NewOrganizationEmailAuthorizeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationEmailAuthorizeLogic {
	return &OrganizationEmailAuthorizeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationEmailAuthorizeLogic) OrganizationEmailAuthorize(req *types.OrganizationEmailAuthorizeRequest) (resp *types.EmptyResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
