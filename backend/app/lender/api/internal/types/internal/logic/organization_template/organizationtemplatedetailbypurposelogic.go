// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package organization_template

import (
	"context"

	"zeebee-backend/app/lender/api/internal/types/internal/svc"
	"zeebee-backend/app/lender/api/internal/types/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrganizationTemplateDetailByPurposeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Get Detail by purpose
func NewOrganizationTemplateDetailByPurposeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationTemplateDetailByPurposeLogic {
	return &OrganizationTemplateDetailByPurposeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationTemplateDetailByPurposeLogic) OrganizationTemplateDetailByPurpose(req *types.OrganizationTemplateDetailByPurposeRequest) (resp *types.OrganizationTemplateDetailResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
