package pricingplanlogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/proto"
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

func (l *AllLogic) All(in *admin.PricingPlanAllRequest) (resp *admin.PricingPlanAllResponse, err error) {

	if in.Data == nil {
		return nil, errorx.New("invalid request: data is required")
	}

	plans, err := l.baseLogic.PricingPlanLogic.All(l.ctx, constants.UserType(in.Data.Type), false)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	data, err := copierx.Trans[[]*dto.PricingPlanEntityDto, []*proto.PricingPlanEntity](&plans, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}

	return &admin.PricingPlanAllResponse{Data: *data}, nil
}
