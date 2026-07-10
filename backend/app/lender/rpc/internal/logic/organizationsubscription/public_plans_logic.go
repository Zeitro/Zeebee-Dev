package organizationsubscriptionlogic

import (
	"context"

	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublicPlansLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
	logx.Logger
}

func NewPublicPlansLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublicPlansLogic {
	return &PublicPlansLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.GetBaseLogic(svcCtx),
		Logger:    logx.WithContext(ctx),
	}
}

func (l *PublicPlansLogic) PublicPlans(in *lender.EmptyRequest) (resp *lender.PublicPricingPlansResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
				logx.Field("service_name", l.svcCtx.Config.Name),
			).Error(err)
		}
	}()

	plans, err := l.baseLogic.PricingPlanLogic.All(l.ctx, constants.UserTypeLender, true)
	if err != nil {
		return nil, err
	}

	data, err := copierx.Trans[[]*dto.PricingPlanEntityDto, []*lender.PricingPlanEntity](&plans)
	if err != nil {
		return nil, err
	}
	return &lender.PublicPricingPlansResponse{Data: *data}, nil
}
