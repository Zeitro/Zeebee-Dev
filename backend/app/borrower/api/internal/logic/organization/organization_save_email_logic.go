package organization

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrganizationSaveEmailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// save email
func NewOrganizationSaveEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationSaveEmailLogic {
	return &OrganizationSaveEmailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationSaveEmailLogic) OrganizationSaveEmail(req *types.OrganizationSaveEmailRequest) (resp *types.EmptyResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
