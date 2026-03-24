package loan

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
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

	rpcReq := &lender.LoanPendingStatusReviewRequest{
		Data: &lender.LoanPendingStatusReviewRequestData{},
	}
	rpcReq.Data, err = copierx.Trans[types.ReviewLoanPendingStatusRequest, lender.LoanPendingStatusReviewRequestData](req)
	if err != nil {
		return nil, err
	}
	_, err = l.svcCtx.LoanRpcClient.PendingStatusReview(l.ctx, rpcReq)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	return &types.EmptyResponse{}, nil
}
