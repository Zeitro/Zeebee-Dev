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

type RenewLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewRenewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RenewLogic {
	return &RenewLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *RenewLogic) Renew(in *admin.OrganizationSubscriptionRenewRequest) (resp *proto.EmptyResponse, err error) {
	_, err = l.svcCtx.TenantClient.OrganizationSubscription.Renew(l.ctx, &tenant.OrganizationSubscriptionRenewRequest{Data: &tenant.OrganizationSubscriptionRenewRequestData{
		UserType: string(constants.UserTypeAdmin),
		Id:       in.Data.Id,
	}})
	return &proto.EmptyResponse{}, err
}
