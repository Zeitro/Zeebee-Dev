package rolelogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/constants"
)

type CreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *CreateLogic) Create(in *admin.RoleCreateRequest) (resp *admin.RoleDetailResponse, err error) {

	claim, err := l.baseLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	input, err := copierx.Trans[admin.RoleCreateRequestData, tenant.RoleCreateRequestData](in.Data)
	if err != nil {
		return nil, err
	}

	input.UserType = string(constants.UserTypeAdmin)
	input.OrganizationId = claim.OrganizationId

	role, err := l.svcCtx.TenantClient.Role.Create(l.ctx, &tenant.RoleCreateRequest{Data: input})
	if err != nil {
		return nil, err
	}

	return &admin.RoleDetailResponse{
		Data: role.Data,
	}, nil
}
