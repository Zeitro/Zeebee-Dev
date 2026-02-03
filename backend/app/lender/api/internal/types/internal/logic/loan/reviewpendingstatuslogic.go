// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package loan

import (
	"context"

	"zeebee-backend/app/lender/api/internal/types/internal/svc"
	"zeebee-backend/app/lender/api/internal/types/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReviewPendingStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// reveiw loan pending status
func NewReviewPendingStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReviewPendingStatusLogic {
	return &ReviewPendingStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReviewPendingStatusLogic) ReviewPendingStatus(req *types.ReviewLoanPendingStatusRequest) (resp *types.EmptyResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
