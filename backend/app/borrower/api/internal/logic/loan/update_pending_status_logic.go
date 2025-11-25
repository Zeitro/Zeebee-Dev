package loan

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatePendingStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// update pending status
func NewUpdatePendingStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePendingStatusLogic {
	return &UpdatePendingStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdatePendingStatusLogic) UpdatePendingStatus(req *types.LoanUpdatePendingStatusRequest) (resp *types.EmptyResponse, err error) {
	// get borrower id
	claim, err := utils.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	_, err = l.svcCtx.LoanRpcClient.UpdateBorrowerLoanPendingStatus(l.ctx, &borrower.LoanUpdatePendingStatusRequest{
		Data: &borrower.LoanUpdatePendingStatusRequestData{
			Id:             req.Id,
			UserId:         claim.UserId,
			OrganizationId: claim.OrganizationId,
			PendingStatus:  req.PendingStatus,
		},
	})

	return
}
