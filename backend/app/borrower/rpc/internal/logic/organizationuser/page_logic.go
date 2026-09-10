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

type PageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PageLogic {
	return &PageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *PageLogic) Page(in *borrower.OrganizationUserPageRequest) (resp *borrower.OrganizationUserPageResponse, err error) {

	claim, err := l.baseLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	pageReq, err := copierx.Trans[borrower.OrganizationUserPageRequestData, tenant.OrganizationUserPageRequestData](in.Data)
	if err != nil {
		return nil, err
	}
	pageReq.UserType = string(constants.UserTypeBorrower)
	pageReq.Filter.OrganizationId = claim.OrganizationId
	pageReq.Filter.UserId = claim.UserId

	results, err := l.svcCtx.TenantClient.OrganizationUser.Page(l.ctx, &tenant.OrganizationUserPageRequest{Data: pageReq})
	if err != nil {
		return nil, err
	}
	data, err := copierx.Trans[tenant.OrganizationUserPageResponseData, borrower.OrganizationUserPageResponseData](results.Data, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}

	return &borrower.OrganizationUserPageResponse{
		Data: data,
	}, nil
}
