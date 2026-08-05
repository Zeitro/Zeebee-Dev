package borrowerpropertytransactionlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/timex"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/proto"
)

type ListByBorrowerPropertyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewListByBorrowerPropertyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListByBorrowerPropertyLogic {
	return &ListByBorrowerPropertyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *ListByBorrowerPropertyLogic) ListByBorrowerProperty(in *borrower.BorrowerPropertyTransactionListByPropertyRequest) (resp *borrower.BorrowerPropertyTransactionListByPropertyResponse, err error) {
	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx, constants.UserTypeBorrower)
	if err != nil {
		return nil, err
	}

	transactions, err := l.svcCtx.BorrowerPropertyTransactionRepo.ListByBorrowerPropertyId(l.ctx, in.Data.BorrowerPropertyId, claim.OrganizationId)
	if err != nil {
		return nil, err
	}

	for i := range transactions {
		fillChecklistNodeFiles(l.ctx, l.svcCtx.FileRepo, &transactions[i].Checklist)
	}

	entities, err := copierx.Trans[[]*model.BorrowerPropertyTransaction, []*proto.BorrowerPropertyTransactionEntity](
		&transactions,
		copierx.TimeToRFC3339Str(),
		copierx.SqlTimePtrToStr(timex.LayoutUsStyle),
	)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	if entities == nil {
		return &borrower.BorrowerPropertyTransactionListByPropertyResponse{}, nil
	}

	return &borrower.BorrowerPropertyTransactionListByPropertyResponse{
		Data: *entities,
	}, nil
}
