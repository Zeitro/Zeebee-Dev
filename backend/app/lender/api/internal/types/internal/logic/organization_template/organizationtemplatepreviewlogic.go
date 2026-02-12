// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package organization_template

import (
	"context"

	"zeebee-backend/app/lender/api/internal/types/internal/svc"
	"zeebee-backend/app/lender/api/internal/types/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrganizationTemplatePreviewLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Preview
func NewOrganizationTemplatePreviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationTemplatePreviewLogic {
	return &OrganizationTemplatePreviewLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationTemplatePreviewLogic) OrganizationTemplatePreview(req *types.OrganizationTemplatePreviewRequest) (resp *types.OrganizationTemplatePreviewResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
