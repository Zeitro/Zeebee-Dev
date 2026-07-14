package subscription

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/proto"

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
	rpcResp, err := l.svcCtx.OrganizationSubscriptionRpcClient.PublicPlans(l.ctx, &proto.EmptyRequest{})
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	plans, err := copierx.Trans[[]*borrower.PricingPlanEntity, []*types.PublicPricingPlanEntity](&rpcResp.Data)
	if err != nil {
		return nil, err
	}
	return &types.PublicPricingPlansResponse{Plans: *plans}, nil
}
