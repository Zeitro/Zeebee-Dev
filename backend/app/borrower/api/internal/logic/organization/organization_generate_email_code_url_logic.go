package organization

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrganizationGenerateEmailCodeUrlLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// generate email code url
func NewOrganizationGenerateEmailCodeUrlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationGenerateEmailCodeUrlLogic {
	return &OrganizationGenerateEmailCodeUrlLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationGenerateEmailCodeUrlLogic) OrganizationGenerateEmailCodeUrl(req *types.OrganizationGenerateEmailCodeUrlRequest) (resp *types.OrganizationGenerateEmailCodeUrlResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
