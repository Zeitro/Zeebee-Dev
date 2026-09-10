package organizationemaillogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/proto"
)

type AuthorizeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewAuthorizeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuthorizeLogic {
	return &AuthorizeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *AuthorizeLogic) Authorize(in *borrower.OrganizationEmailAuthorizeRequest) (resp *proto.EmptyResponse, err error) {

	claim, err := l.baseLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	input, err := copierx.Trans[borrower.OrganizationEmailAuthorizeRequestData, tenant.OrganizationEmailAuthorizeRequestData](in.Data)
	if err != nil {
		return nil, err
	}
	input.UserType = string(constants.UserTypeBorrower)
	input.UserId = claim.UserId
	input.OrganizationId = claim.OrganizationId
	_, err = l.svcCtx.TenantClient.OrganizationEmail.Authorize(l.ctx, &tenant.OrganizationEmailAuthorizeRequest{Data: input})
	if err != nil {
		return nil, err
	}
	return &proto.EmptyResponse{}, nil
}
