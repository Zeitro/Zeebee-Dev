package borrowerpropertytransactionlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/model"

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
		return nil, err
	}
	var options *[]*borrower.BorrowerPropertyTransactionOptionEntity
	options, err = copierx.Trans[[]*model.BorrowerPropertyTransaction, []*borrower.BorrowerPropertyTransactionOptionEntity](&transactions, copierx.SqlTimePtrToRFC3339Str(), copierx.TimeToRFC3339Str())
	if err != nil {
		return
	}
	if options == nil {
		return &borrower.BorrowerPropertyTransactionAllResponse{}, nil
	}

	return &borrower.BorrowerPropertyTransactionAllResponse{
		Data: *options,
	}, nil
}
