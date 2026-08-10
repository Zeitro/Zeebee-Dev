package rolelogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/pb/borrower"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
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

func (l *UpdateLogic) Update(in *borrower.RoleUpdateRequest) (resp *borrower.RoleUpdateResponse, err error) {

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	input, err := copierx.Trans[borrower.RoleUpdateRequestData, tenant.RoleUpdateRequestData](in.Data)
	if err != nil {
		return nil, err
	}
	input.UserType = string(constants.UserTypeBorrower)
	input.OrganizationId = claim.OrganizationId

	role, err := l.svcCtx.TenantClient.Role.Update(l.ctx, &tenant.RoleUpdateRequest{Data: input})
	if err != nil {
		return nil, err
	}

	return &borrower.RoleUpdateResponse{
		Data: role.Data,
	}, nil
}
