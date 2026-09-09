package organizationlogic

import (
	"context"

	"zeebee-backend/app/tenant/rpc/internal/logic"
	"zeebee-backend/app/tenant/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/proto"
)

type UpdateOrderLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewUpdateOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateOrderLogic {
	return &UpdateOrderLogic{ctx: ctx, svcCtx: svcCtx, baseLogic: logic.GetBaseLogic(svcCtx)}
}

// UpdateOrder sets the sort order of an organization.
func (l *UpdateOrderLogic) UpdateOrder(in *tenant.OrganizationUpdateOrderRequest) (resp *proto.EmptyResponse, err error) {
	if in.Data == nil {
		return nil, errorx.ErrorEmptyRequestData
	}
	if _, err = logic.ParseUserType(in.Data.UserType); err != nil {
		return
	}

	if err = l.baseLogic.OrganizationLogic.UpdateOrder(l.ctx, in.Data.Id, int(in.Data.Order)); err != nil {
		return
	}
	return &proto.EmptyResponse{}, nil
}
