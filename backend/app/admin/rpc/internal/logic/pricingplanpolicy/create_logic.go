package pricingplanpolicylogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/proto"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type CreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.NewBaseLogic(svcCtx),
	}
}

func (l *CreateLogic) Create(in *admin.PricingPlanPolicyCreateRequest) (resp *admin.PricingPlanPolicyDetailResponse, err error) {
	defer func() {
		if err != nil {
			l.svcCtx.Logger.WithContext(l.ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
			).Error(err)
		}
	}()

	dtoData, err := copierx.Trans[admin.PricingPlanPolicyCreateRequestData, dto.PricingPlanPolicyCreateRequestDto](in.Data)
	if err != nil {
		return nil, err
	}

	result, err := l.baseLogic.PricingPlanPolicyLogic.Create(l.ctx, *dtoData)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	protoData, err := copierx.Trans[dto.PricingPlanPolicyEntityDto, proto.PricingPlanPolicyEntity](result, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}

	return &admin.PricingPlanPolicyDetailResponse{Data: protoData}, nil
}
