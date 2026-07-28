package borrowerpropertytransactionlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/timex"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"

	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/model"
)

type OperationListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewOperationListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OperationListLogic {
	return &OperationListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *OperationListLogic) OperationList(in *borrower.TransactionOperationListRequest) (resp *borrower.TransactionOperationListResponse, err error) {
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

	input, err := copierx.Trans[borrower.TransactionOperationListRequestData, dto.OperationListRequestDto](in.Data, copierx.StrToSqlTimePtr(timex.LayoutUsStyle))
	if err != nil {
		return nil, err
	}
	input.EntityId = in.Data.TransactionId
	input.EntityType = constants.OperationEntityTypeBorrowerPropertyTransaction
	input.Operation = constants.OperationAction(in.Data.OperationType)

	list, err := l.svcCtx.OperationRepo.GetOperationList(l.ctx, *input)
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
