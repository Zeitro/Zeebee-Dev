// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package loan

import (
	"context"

	"zeebee-backend/app/lender/api/internal/types/internal/svc"
	"zeebee-backend/app/lender/api/internal/types/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatePendingStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// update loan pending status
func NewUpdatePendingStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePendingStatusLogic {
	return &UpdatePendingStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdatePendingStatusLogic) UpdatePendingStatus(req *types.UpdateLoanPendingStatusRequest) (resp *types.EmptyResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
