package organization

import (
	"context"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LenderDiscoverPreviewLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// lender discover preview (OpenAI + dedupe, no insert)
func NewLenderDiscoverPreviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LenderDiscoverPreviewLogic {
	return &LenderDiscoverPreviewLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LenderDiscoverPreviewLogic) LenderDiscoverPreview(req *types.EmptyRequest) (resp *types.LenderDiscoverPreviewResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
