package organizationlogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/proto"
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

// Delete deletes an organization and all its associated data.
// Deletion order (respecting foreign key constraints):
// 1. organization_users (via DeleteByOrganizationIdWithTx)
// 2. roles (via GORM clause.Associations cascade delete)
// 3. organization
// Note: Users are NOT deleted as they may belong to other organizations
func (l *DeleteLogic) Delete(in *proto.IdRequest) (resp *proto.EmptyResponse, err error) {

	return l.svcCtx.TenantClient.Organization.Delete(l.ctx, &tenant.OrganizationDeleteRequest{Data: &tenant.OrganizationDeleteRequestData{
		UserType: string(constants.UserTypeAdmin),
		Id:       in.Data,
	}})
}
