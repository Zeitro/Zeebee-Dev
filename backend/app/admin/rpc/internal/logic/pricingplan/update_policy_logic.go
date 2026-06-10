package pricingplanlogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/data/dto"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatePolicyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewUpdatePolicyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePolicyLogic {
	return &UpdatePolicyLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *UpdatePolicyLogic) UpdatePolicy(in *admin.PricingPlanPolicyUpdateRequest) (resp *admin.PricingPlanPolicyUpdateResponse, err error) {

	dtoData, err := copierx.Trans[admin.PricingPlanPolicyUpdateRequestData, dto.PricingPlanPolicyUpdateDto](in.Data)
	if err != nil {
		return nil, err
	}

	result, err := l.baseLogic.PricingPlanPolicyLogic.Update(l.ctx, in.Data.Id, *dtoData)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	protoData, err := copierx.Trans[dto.PricingPlanPolicyEntityDto, admin.PricingPlanPolicyEntity](result, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}

	return &admin.PricingPlanPolicyUpdateResponse{Data: protoData}, nil
}
