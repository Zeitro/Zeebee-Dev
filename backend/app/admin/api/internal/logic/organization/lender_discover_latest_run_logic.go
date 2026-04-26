package organization

import (
	"context"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LenderDiscoverLatestRunLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// lender discover latest run
func NewLenderDiscoverLatestRunLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LenderDiscoverLatestRunLogic {
	return &LenderDiscoverLatestRunLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LenderDiscoverLatestRunLogic) LenderDiscoverLatestRun(req *types.EmptyRequest) (resp *types.LenderDiscoverRunResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
