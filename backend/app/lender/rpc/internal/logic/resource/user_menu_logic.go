package resourcelogic

import (
	"context"

	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/proto"
)

type UserMenuLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewUserMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserMenuLogic {
	return &UserMenuLogic{ctx: ctx, svcCtx: svcCtx, baseLogic: logic.GetBaseLogic(svcCtx)}
}

// UserMenu proxies to tenant; the caller's organization/membership come from the
// locally verified JWT claim.
func (l *UserMenuLogic) UserMenu(in *proto.EmptyRequest) (resp *lender.ResourceTreeResponse, err error) {
	claim, err := l.baseLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	tenantResp, err := l.svcCtx.TenantClient.Resource.UserMenu(l.ctx, &tenant.UserMenuRequest{
		Data: &tenant.UserMenuRequestData{
			UserType:           string(constants.UserTypeLender),
			OrganizationId:     claim.OrganizationId,
			OrganizationUserId: claim.OrganizationUserId,
		},
	})
	if err != nil {
		return nil, err
	}
	return &lender.ResourceTreeResponse{Data: tenantResp.Data}, nil
}
