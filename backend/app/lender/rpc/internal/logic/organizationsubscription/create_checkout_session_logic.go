package organizationsubscriptionlogic

import (
	"context"

	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateCheckoutSessionLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
	logx.Logger
}

func NewCreateCheckoutSessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCheckoutSessionLogic {
	return &CreateCheckoutSessionLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.GetBaseLogic(svcCtx),
		Logger:    logx.WithContext(ctx),
	}
}

func (l *CreateCheckoutSessionLogic) CreateCheckoutSession(in *lender.CheckoutSessionCreateRequest) (resp *lender.CheckoutSessionCreateResponse, err error) {
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

	checkoutUrl, err := l.baseLogic.StripePaymentLogic.CreateCheckoutSession(l.ctx, dto.StripeCheckoutRequestDto{
		OrganizationId:      claim.OrganizationId,
		UserId:              claim.UserId,
		PricingPlanId:       in.Data.PricingPlanId,
		PricingPlanPolicyId: in.Data.PricingPlanPolicyId,
		UserType:            constants.UserTypeLender,
	})
	if err != nil {
		return nil, err
	}

	return &lender.CheckoutSessionCreateResponse{CheckoutUrl: checkoutUrl}, nil
}
