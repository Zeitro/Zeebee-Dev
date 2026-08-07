package resourcelogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/domains/zeebee/constants"
)

type ResourceTreeLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewResourceTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResourceTreeLogic {
	return &ResourceTreeLogic{ctx: ctx, svcCtx: svcCtx, baseLogic: logic.GetBaseLogic(svcCtx)}
}

// ResourceTree proxies to tenant; tenant is the RBAC single authority. For an
// organization-scoped tree the caller's org comes from the JWT claim.
func (l *ResourceTreeLogic) ResourceTree(in *admin.ResourceTreeRequest) (resp *admin.ResourceTreeResponse, err error) {
	req := &tenant.ResourceTreeRequest{
		Data: &tenant.ResourceTreeRequestData{
			UserType: string(constants.UserTypeAdmin),
			Type:     in.Data.Type,
			Scope:    in.Data.Scope,
		},
	}
	if constants.ResourceTreeScope(in.Data.Scope) == constants.ResourceTreeScopeOrganization {
		claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
		if err != nil {
			return nil, err
		}
		req.Data.OrganizationId = claim.OrganizationId
	}
	tenantResp, err := l.svcCtx.TenantClient.Resource.ResourceTree(l.ctx, req)
	if err != nil {
		return nil, err
	}
	return &admin.ResourceTreeResponse{Data: tenantResp.Data}, nil
}
