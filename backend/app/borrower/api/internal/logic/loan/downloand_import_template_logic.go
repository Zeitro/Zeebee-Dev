package loan

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DownloandImportTemplateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// downloand import template
func NewDownloandImportTemplateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DownloandImportTemplateLogic {
	return &DownloandImportTemplateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DownloandImportTemplateLogic) DownloandImportTemplate(req *types.EmptyRequest) (resp *types.EmptyResponse, err error) {
	// 逻辑处理完成，返回空响应
	// 实际的文件下载在handler中处理
	return &types.EmptyResponse{}, nil
}
