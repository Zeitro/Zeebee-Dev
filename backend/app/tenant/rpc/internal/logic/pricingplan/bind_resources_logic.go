package pricingplanlogic

import (
	"context"

	"zeebee-backend/app/tenant/rpc/internal/logic"
	"zeebee-backend/app/tenant/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/domains/zeebee/proto"

	"zeebee-backend/common/errorx"
)

type BindResourcesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewBindResourcesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BindResourcesLogic {
	return &BindResourcesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *BindResourcesLogic) BindResources(in *tenant.PricingPlanBindResourcesRequest) (resp *proto.EmptyResponse, err error) {
	if in.Data == nil {
		return nil, errorx.ErrorEmptyRequestData
	}
	if _, err = logic.ParseUserType(in.Data.UserType); err != nil {
		return nil, err
	}
	if err = l.baseLogic.PricingPlanLogic.BindResources(l.ctx, in.Data.Id, in.Data.ResourceIds); err != nil {
		return nil, errorx.WithStack(err)
	}

	return &proto.EmptyResponse{}, nil
}
