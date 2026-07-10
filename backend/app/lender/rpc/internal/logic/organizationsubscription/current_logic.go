package organizationsubscriptionlogic

import (
	"context"

	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/timex"
	"zeebee-backend/domains/zeebee/data/dto"

	"github.com/zeromicro/go-zero/core/logx"
)

type CurrentLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
	logx.Logger
}

func NewCurrentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CurrentLogic {
	return &CurrentLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.GetBaseLogic(svcCtx),
		Logger:    logx.WithContext(ctx),
	}
}

func (l *CurrentLogic) Current(in *lender.EmptyRequest) (resp *lender.SubscriptionCurrentResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
				logx.Field("service_name", l.svcCtx.Config.Name),
			).Error(err)
		}
	}()

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, errorx.Wrap(err, "fail to get user claim from context")
	}

	subscription, err := l.baseLogic.StripePaymentLogic.Current(l.ctx, claim.OrganizationId)
	if err != nil {
		return nil, err
	}
	if subscription == nil {
		return &lender.SubscriptionCurrentResponse{}, nil
	}

	data, err := copierx.Trans[dto.OrganizationSubscriptionEntity, lender.OrganizationSubscriptionEntity](subscription, copierx.TimeToRFC3339Str(), copierx.SqlTimePtrToStr(timex.LayoutUsStyle))
	if err != nil {
		return nil, err
	}
	return &lender.SubscriptionCurrentResponse{Data: data}, nil
}
