package loan

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateStageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// update stage
func NewUpdateStageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateStageLogic {
	return &UpdateStageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateStageLogic) UpdateStage(req *types.LoanUpdatePendingRequest) (resp *types.EmptyResponse, err error) {
	_, err = l.svcCtx.LoanRpcClient.UpdateBorrowerLoanStage(l.ctx, &borrower.LoanUpdateStageRequest{
		Data: &borrower.LoanUpdateStageRequestData{
			Id:    req.Id,
			Stage: req.Stage,
		},
	})

	return &types.EmptyResponse{}, err
}
