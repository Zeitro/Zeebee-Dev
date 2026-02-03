// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package organization_template

import (
	"context"

	"zeebee-backend/app/lender/api/internal/types/internal/svc"
	"zeebee-backend/app/lender/api/internal/types/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrganizationTemplatePageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// page
func NewOrganizationTemplatePageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationTemplatePageLogic {
	return &OrganizationTemplatePageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationTemplatePageLogic) OrganizationTemplatePage(req *types.OrganizationTemplatePageRequest) (resp *types.OrganizationTemplatePageResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
