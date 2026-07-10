package borrowerpropertytransactionlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/timex"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/proto"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type AllLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewAllLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AllLogic {
	return &AllLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *AllLogic) All(in *proto.EmptyRequest) (resp *borrower.BorrowerPropertyTransactionAllResponse, err error) {
	defer func() {
		if err != nil {
			l.svcCtx.Logger.WithContext(l.ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
			).Error(err)
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
	var options *[]*proto.BorrowerPropertyTransactionOptionEntity
	options, err = copierx.Trans[[]*model.BorrowerPropertyTransaction, []*proto.BorrowerPropertyTransactionOptionEntity](&transactions, copierx.SqlTimePtrToStr(timex.LayoutUsStyle), copierx.TimeToRFC3339Str())
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
