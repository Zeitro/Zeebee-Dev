package organization

import (
	"context"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LenderDiscoverGetSettingsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// lender discover get settings
func NewLenderDiscoverGetSettingsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LenderDiscoverGetSettingsLogic {
	return &LenderDiscoverGetSettingsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LenderDiscoverGetSettingsLogic) LenderDiscoverGetSettings(req *types.EmptyRequest) (resp *types.LenderDiscoverSettingsResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
