package organizationuserlogic

import (
	"context"

	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
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

func (l *AllLogic) All(in *proto.EmptyRequest) (resp *lender.OrganizationUserAllResponse, err error) {

	claimUser, err := l.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	results, err := l.svcCtx.TenantClient.OrganizationUser.AllOrganizationUsers(l.ctx, &tenant.OrganizationUserAllOrganizationUsersRequest{
		Data: &tenant.OrganizationUserAllOrganizationUsersRequestData{
			UserType:       string(constants.UserTypeLender),
			OrganizationId: claimUser.OrganizationId,
		},
	})
	if err != nil {
		return nil, err
	}

	return &lender.OrganizationUserAllResponse{
		Data: results.Data,
	}, nil
}
