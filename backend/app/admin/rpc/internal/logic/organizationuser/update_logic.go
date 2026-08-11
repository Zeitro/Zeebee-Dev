package organizationuserlogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/copierx"
)

type UpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *UpdateLogic) Update(in *admin.OrganizationUserUpdateRequest) (resp *admin.OrganizationUserDetailResponse, err error) {

	claim, err := l.baseLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	input, err := copierx.Trans[admin.OrganizationUserUpdateRequestData, tenant.OrganizationUserUpdateRequestData](in.Data)
	if err != nil {
		return nil, err
	}

	input.UserType = string(claim.UserType)
	input.OrganizationId = claim.OrganizationId

	result, err := l.svcCtx.TenantClient.OrganizationUser.UpdateOrganizationUser(l.ctx, &tenant.OrganizationUserUpdateRequest{Data: input})
	if err != nil {
		return nil, err
	}
	return &admin.OrganizationUserDetailResponse{
		Data: result.Data,
	}, nil
}
