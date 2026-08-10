package pricingplanlogic

import (
	"context"

	"zeebee-backend/app/tenant/rpc/internal/logic"
	"zeebee-backend/app/tenant/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/proto"

	"zeebee-backend/common/errorx"
)

type DeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogic {
	return &DeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *DeleteLogic) Delete(in *tenant.PricingPlanDeleteRequest) (resp *proto.EmptyResponse, err error) {
	if in.Data == nil {
		return nil, errorx.ErrorEmptyRequestData
	}
	if _, err = logic.ParseUserType(in.Data.UserType); err != nil {
		return nil, err
	}
	if err = l.baseLogic.PricingPlanLogic.Delete(l.ctx, in.Data.Id); err != nil {
		return nil, errorx.WithStack(err)
	}
	l.baseLogic.RecordAuditFromClaim(l.ctx, constants.OperationEntityTypePricingPlan, in.Data.Id, constants.OperationActionPlanDelete)

	return &proto.EmptyResponse{}, nil
}
