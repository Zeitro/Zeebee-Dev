package organization

import (
	"context"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LenderDiscoverUpdateSettingsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// lender discover update settings
func NewLenderDiscoverUpdateSettingsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LenderDiscoverUpdateSettingsLogic {
	return &LenderDiscoverUpdateSettingsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LenderDiscoverUpdateSettingsLogic) LenderDiscoverUpdateSettings(req *types.LenderDiscoverUpdateSettingsRequest) (resp *types.LenderDiscoverSettingsResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
