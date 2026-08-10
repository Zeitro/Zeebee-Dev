package organizationlogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/proto"
)

type UpdateOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewUpdateOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateOrderLogic {
	return &UpdateOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *UpdateOrderLogic) UpdateOrder(in *admin.OrganizationUpdateOrderRequest) (resp *proto.EmptyResponse, err error) {

	return l.svcCtx.TenantClient.Organization.UpdateOrder(l.ctx, &tenant.OrganizationUpdateOrderRequest{Data: &tenant.OrganizationUpdateOrderRequestData{
		UserType: string(constants.UserTypeAdmin),
		Id:       in.Data.Id,
		Order:    in.Data.Order,
	}})
}
