package subscription

import (
	"context"

	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type SubscriptionCheckoutLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// create stripe checkout session
func NewSubscriptionCheckoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SubscriptionCheckoutLogic {
	return &SubscriptionCheckoutLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SubscriptionCheckoutLogic) SubscriptionCheckout(req *types.SubscriptionCheckoutRequest) (resp *types.SubscriptionCheckoutResponse, err error) {
	rpcResp, err := l.svcCtx.OrganizationSubscriptionRpcClient.CreateCheckoutSession(l.ctx, &lender.CheckoutSessionCreateRequest{
		Data: &lender.CheckoutSessionCreateRequestData{
			PricingPlanId:       req.PricingPlanId,
			PricingPlanPolicyId: req.PricingPlanPolicyId,
		},
	})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	return &types.SubscriptionCheckoutResponse{CheckoutUrl: rpcResp.CheckoutUrl}, nil
}
