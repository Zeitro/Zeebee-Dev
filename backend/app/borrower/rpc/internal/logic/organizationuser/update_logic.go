package organizationuserlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/constants"
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

func (l *UpdateLogic) Update(in *borrower.OrganizationUserUpdateRequest) (resp *borrower.OrganizationUserDetailResponse, err error) {

	claim, err := l.baseLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	input, err := copierx.Trans[borrower.OrganizationUserUpdateRequestData, tenant.OrganizationUserUpdateRequestData](in.Data)
	if err != nil {
		return nil, err
	}
	input.UserType = string(constants.UserTypeBorrower)
	input.OrganizationId = claim.OrganizationId

	result, err := l.svcCtx.TenantClient.OrganizationUser.UpdateOrganizationUser(l.ctx, &tenant.OrganizationUserUpdateRequest{Data: input})
	if err != nil {
		return nil, err
	}
	return &borrower.OrganizationUserDetailResponse{
		Data: result.Data,
	}, nil
}
