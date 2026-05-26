package borrowerpropertytransactionlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/domains/zeebee/constants"

	"github.com/zeromicro/go-zero/core/logx"
	"zeebee-backend/common/errorx"
)

type OperationStatsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewOperationStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OperationStatsLogic {
	return &OperationStatsLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *OperationStatsLogic) OperationStats(in *borrower.TransactionOperationStatsRequest) (resp *borrower.TransactionOperationStatsResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, errorx.Wrap(err, "failed to get user claim")
	}

	// validate transaction ownership
	_, err = l.svcCtx.BorrowerPropertyTransactionRepo.GetOrganizationTransaction(l.ctx, in.Data.TransactionId, claim.OrganizationId)
	if err != nil {
		return nil, errorx.Wrap(err, "transaction not found")
	}

	items, total, err := l.svcCtx.OperationRepo.GetOperationStats(l.ctx, in.Data.TransactionId, constants.OperationEntityTypeBorrowerPropertyTransaction)
	if err != nil {
		return nil, errorx.Wrap(err, "failed to get operation stats")
	}

	respItems := make([]*borrower.OperationStatsItem, 0, len(items))
	for _, item := range items {
		respItems = append(respItems, &borrower.OperationStatsItem{
			Type:  item.Type,
			Count: item.Count,
		})
	}

	return &borrower.TransactionOperationStatsResponse{
		Data: &borrower.TransactionOperationStatsResponseData{
			Total: total,
			Items: respItems,
		},
	}, nil
}
