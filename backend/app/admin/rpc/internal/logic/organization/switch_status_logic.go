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

type SwitchStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	*logic.BaseLogic
}

func NewSwitchStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SwitchStatusLogic {
	return &SwitchStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		BaseLogic: logic.GetBaseLogic(svcCtx),
	}
}

// SwitchStatus switches organization status between Active and Disabled.
// It updates both the organization status and all associated organization users' status
// in a single transaction for data consistency.
func (l *SwitchStatusLogic) SwitchStatus(in *admin.OrganizationSwitchStatusRequest) (resp *proto.EmptyResponse, err error) {

	return l.svcCtx.TenantClient.Organization.SwitchStatus(l.ctx, &tenant.OrganizationSwitchStatusRequest{Data: &tenant.OrganizationSwitchStatusRequestData{
		UserType: string(constants.UserTypeAdmin),
		Id:       in.Data.Id,
		Status:   in.Data.Status,
	}})
}
