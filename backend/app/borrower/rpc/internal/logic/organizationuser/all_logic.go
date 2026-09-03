package organizationuserlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/proto"
)

type AllLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	*logic.BaseLogic
}

func NewAllLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AllLogic {
	return &AllLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		BaseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *AllLogic) All(in *proto.EmptyRequest) (resp *borrower.OrganizationUserAllResponse, err error) {

	claimUser, err := l.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	results, err := l.svcCtx.TenantClient.OrganizationUser.AllOrganizationUsers(l.ctx, &tenant.OrganizationUserAllOrganizationUsersRequest{
		Data: &tenant.OrganizationUserAllOrganizationUsersRequestData{
			UserType:       string(constants.UserTypeBorrower),
			OrganizationId: claimUser.OrganizationId,
		},
	})
	if err != nil {
		return nil, err
	}

	return &borrower.OrganizationUserAllResponse{
		Data: results.Data,
	}, nil
}
