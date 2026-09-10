package organizationsubscriptionlogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/proto"
)

type UpdateExpireDateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewUpdateExpireDateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateExpireDateLogic {
	return &UpdateExpireDateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *UpdateExpireDateLogic) UpdateExpireDate(in *admin.OrganizationSubscriptionUpdateExpireDateRequest) (resp *proto.EmptyResponse, err error) {

	return l.svcCtx.TenantClient.OrganizationSubscription.UpdateExpireDate(l.ctx, &tenant.OrganizationSubscriptionUpdateExpireDateRequest{Data: &tenant.OrganizationSubscriptionUpdateExpireDateRequestData{
		UserType:   string(constants.UserTypeAdmin),
		Id:         in.Data.Id,
		ExpireDate: in.Data.ExpireDate,
	}})
}
