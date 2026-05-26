package borrowerpropertytransactionlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/domains/zeebee/constants"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/model"
)

type OperationListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewOperationListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OperationListLogic {
	return &OperationListLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *OperationListLogic) OperationList(in *borrower.TransactionOperationListRequest) (resp *borrower.TransactionOperationListResponse, err error) {
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
		err = errorx.Wrap(err, "failed to find transaction")
		return
	}

	list, err := l.svcCtx.OperationRepo.GetOperationList(l.ctx, in.Data.TransactionId, constants.OperationEntityTypeBorrowerPropertyTransaction, in.Data.OperationType)
	if err != nil {
		return nil, errorx.Wrap(err, "failed to get operation list")
	}

	respItems, err := copierx.Trans[[]*model.Operation, []*borrower.OperationListItem](&list, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, errorx.Wrap(err, "failed to convert operation list")
	}

	return &borrower.TransactionOperationListResponse{
		Data: &borrower.TransactionOperationListResponseData{
			List: *respItems,
		},
	}, nil
}
