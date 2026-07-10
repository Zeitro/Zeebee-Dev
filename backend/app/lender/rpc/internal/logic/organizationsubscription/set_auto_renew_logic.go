package organizationsubscriptionlogic

import (
	"context"

	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/data/dto"

	"github.com/zeromicro/go-zero/core/logx"
)

type SetAutoRenewLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
	logx.Logger
}

func NewSetAutoRenewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetAutoRenewLogic {
	return &SetAutoRenewLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.GetBaseLogic(svcCtx),
		Logger:    logx.WithContext(ctx),
	}
}

func (l *SetAutoRenewLogic) SetAutoRenew(in *lender.SubscriptionAutoRenewRequest) (resp *lender.EmptyResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
				logx.Field("service_name", l.svcCtx.Config.Name),
			).Error(err)
		}
	}()

	if in.Data == nil {
		return nil, errorx.New("invalid request")
	}

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, errorx.Wrap(err, "fail to get user claim from context")
	}

	err = l.baseLogic.StripePaymentLogic.SetAutoRenew(l.ctx, dto.StripeAutoRenewRequestDto{
		OrganizationId: claim.OrganizationId,
		UserId:         claim.UserId,
		Enabled:        in.Data.Enabled,
	})
	if err != nil {
		return nil, err
	}

	return &lender.EmptyResponse{}, nil
}
