package borrowerpropertytransactionlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"

	"github.com/zeromicro/go-zero/core/logx"
	"zeebee-backend/common/errorx"
)

type AllLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewAllLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AllLogic {
	return &AllLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *AllLogic) All(in *borrower.BorrowerPropertyTransactionAllRequest) (resp *borrower.BorrowerPropertyTransactionAllResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	transactions, err := l.svcCtx.BorrowerPropertyTransactionRepo.AllByOrganization(l.ctx, claim.OrganizationId)
	if err != nil {
		return nil, errorx.Wrap(err, "Failed to get transactions")
	}

	options := make([]*borrower.BorrowerPropertyTransactionOptionEntity, 0, len(transactions))
	for _, t := range transactions {
		options = append(options, &borrower.BorrowerPropertyTransactionOptionEntity{
			Id:                 t.Id,
			Code:               t.Code,
			Type:               string(t.Type),
			Status:             string(t.Status),
			PropertyType:       string(t.PropertyType),
			BorrowerPropertyId: t.BorrowerPropertyId,
		})
	}

	return &borrower.BorrowerPropertyTransactionAllResponse{
		Data: options,
	}, nil
}
