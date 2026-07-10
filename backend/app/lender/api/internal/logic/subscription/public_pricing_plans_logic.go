package subscription

import (
	"context"

	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublicPricingPlansLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// list public pricing plans
func NewPublicPricingPlansLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublicPricingPlansLogic {
	return &PublicPricingPlansLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PublicPricingPlansLogic) PublicPricingPlans(req *types.EmptyRequest) (resp *types.PublicPricingPlansResponse, err error) {
	rpcResp, err := l.svcCtx.OrganizationSubscriptionRpcClient.PublicPlans(l.ctx, &lender.EmptyRequest{})
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	plans, err := copierx.Trans[[]*lender.PricingPlanEntity, []*types.PublicPricingPlanEntity](&rpcResp.Data)
	if err != nil {
		return nil, err
	}
	return &types.PublicPricingPlansResponse{Plans: *plans}, nil
}
