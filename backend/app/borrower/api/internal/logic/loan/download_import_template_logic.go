package loan

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DownloadImportTemplateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// download import template
func NewDownloadImportTemplateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DownloadImportTemplateLogic {
	return &DownloadImportTemplateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DownloadImportTemplateLogic) DownloadImportTemplate(req *types.EmptyRequest) (resp *types.EmptyResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
