package organizationuserlogic

import (
	"context"

	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
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

func (l *UpdateLogic) Update(in *lender.OrganizationUserUpdateRequest) (resp *lender.OrganizationUserDetailResponse, err error) {

	claim, err := l.baseLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	input, err := copierx.Trans[lender.OrganizationUserUpdateRequestData, tenant.OrganizationUserUpdateRequestData](in.Data)
	if err != nil {
		return nil, err
	}
	input.UserType = string(constants.UserTypeLender)
	input.OrganizationId = claim.OrganizationId

	result, err := l.svcCtx.TenantClient.OrganizationUser.UpdateOrganizationUser(l.ctx, &tenant.OrganizationUserUpdateRequest{Data: input})
	if err != nil {
		return nil, err
	}
	return &lender.OrganizationUserDetailResponse{
		Data: result.Data,
	}, nil
}
