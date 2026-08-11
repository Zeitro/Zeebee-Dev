package resourcelogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
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

// ResourceTree proxies to tenant; tenant is the RBAC single authority. Borrower
// resources are always scoped to the caller's organization.
func (l *ResourceTreeLogic) ResourceTree(in *borrower.ResourceTreeRequest) (resp *borrower.ResourceTreeResponse, err error) {
	claim, err := l.baseLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	tenantResp, err := l.svcCtx.TenantClient.Resource.ResourceTree(l.ctx, &tenant.ResourceTreeRequest{
		Data: &tenant.ResourceTreeRequestData{
			UserType:       string(constants.UserTypeBorrower),
			Type:           string(constants.UserTypeBorrower),
			Scope:          string(constants.ResourceTreeScopeOrganization),
			OrganizationId: claim.OrganizationId,
		},
	})
	if err != nil {
		return nil, err
	}
	return &borrower.ResourceTreeResponse{Data: tenantResp.Data}, nil
}
