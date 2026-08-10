package organizationuserlogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/copierx"
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

func (l *CreateLogic) Create(in *admin.OrganizationUserCreateRequest) (resp *admin.OrganizationUserDetailResponse, err error) {

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	input, err := copierx.Trans[admin.OrganizationUserCreateRequestData, tenant.OrganizationUserCreateAdminRequestData](in.Data)
	if err != nil {
		return nil, err
	}

	input.UserType = string(claim.UserType)
	input.OrganizationId = claim.OrganizationId
	user, err := l.svcCtx.TenantClient.OrganizationUser.CreateAdminOrganizationUser(l.ctx, &tenant.OrganizationUserCreateAdminRequest{Data: input})
	if err != nil {
		return nil, err
	}

	return &admin.OrganizationUserDetailResponse{
		Data: user.Data,
	}, nil
}
