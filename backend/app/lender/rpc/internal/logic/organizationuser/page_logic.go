package organizationuserlogic

import (
	"context"

	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/constants"

	"zeebee-backend/common/errorx"
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

func (l *PageLogic) Page(in *lender.OrganizationUserPageRequest) (resp *lender.OrganizationUserPageResponse, err error) {

	claim, err := l.baseLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	pageReq, err := copierx.Trans[lender.OrganizationUserPageRequestData, tenant.OrganizationUserPageRequestData](in.Data)
	// err = copier.Copy(&pageReq, &in.Data)
	if err != nil {
		return nil, err
	}

	pageReq.UserType = string(constants.UserTypeLender)
	pageReq.Filter.OrganizationId = claim.OrganizationId
	pageReq.Filter.UserId = claim.UserId

	results, err := l.svcCtx.TenantClient.OrganizationUser.Page(l.ctx, &tenant.OrganizationUserPageRequest{Data: pageReq})
	if err != nil {
		return nil, err
	}
	data, err := copierx.Trans[tenant.OrganizationUserPageResponseData, lender.OrganizationUserPageResponseData](results.Data, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	return &lender.OrganizationUserPageResponse{
		Data: data,
	}, nil
}
