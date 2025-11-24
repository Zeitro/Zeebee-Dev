package loan

import (
	"context"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/utils"
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
	err = copier.CopyWithOption(&rpcReq.Data, req, copier.Option{
		IgnoreEmpty:   false,
		DeepCopy:      true,
		CaseSensitive: true,
	})
	if err != nil {
		return nil, errorx.GetFailedTransReqErr(err)
	}

	claim, err := utils.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	rpcReq.Data.OrganizationId = claim.OrganizationId
	rpcReq.Data.UserId = claim.UserId

	_, err = l.svcCtx.LoanRpcClient.PendingStatusReview(l.ctx, rpcReq)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	return &types.EmptyResponse{}, nil
}
