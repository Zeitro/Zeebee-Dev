package resourcelogic

import (
	"context"

	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
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

// ResourceTree proxies to tenant; tenant is the RBAC single authority. Lender
// resources are always scoped to the caller's organization.
func (l *ResourceTreeLogic) ResourceTree(in *lender.ResourceTreeRequest) (resp *lender.ResourceTreeResponse, err error) {
	claim, err := l.baseLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	tenantResp, err := l.svcCtx.TenantClient.Resource.ResourceTree(l.ctx, &tenant.ResourceTreeRequest{
		Data: &tenant.ResourceTreeRequestData{
			UserType:       string(constants.UserTypeLender),
			Type:           string(constants.UserTypeLender),
			Scope:          string(constants.ResourceTreeScopeOrganization),
			OrganizationId: claim.OrganizationId,
		},
	})
	if err != nil {
		return nil, err
	}
	return &lender.ResourceTreeResponse{Data: tenantResp.Data}, nil
}
