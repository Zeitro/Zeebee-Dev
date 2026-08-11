package organizationuserlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
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

func (l *CreateLogic) Create(in *borrower.OrganizationUserCreateRequest) (resp *borrower.OrganizationUserDetailResponse, err error) {

	claimUser, err := l.baseLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	input := &tenant.OrganizationUserInviteRequestData{
		UserType:       string(constants.UserTypeBorrower),
		OrganizationId: claimUser.OrganizationId,
		Email:          in.Data.Email,
		Phone:          in.Data.Phone,
		Remark:         in.Data.Remark,
		RoleIds:        in.Data.RoleIds,
		Operator:       claimUser.Username,
	}
	result, err := l.svcCtx.TenantClient.OrganizationUser.InviteOrganizationUser(l.ctx, &tenant.OrganizationUserInviteRequest{Data: input})
	if err != nil {
		return nil, err
	}

	return &borrower.OrganizationUserDetailResponse{
		Data: result.Data,
	}, nil
}
