package organizationlogic

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

type LenderPageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewLenderPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LenderPageLogic {
	return &LenderPageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *LenderPageLogic) LenderPage(in *proto.CommonPageRequest) (resp *borrower.OrganizationLenderPageResponse, err error) {

	pageReq, err := copierx.Trans[proto.CommonPageRequestData, tenant.OrganizationPageRequestData](in.Data)
	if err != nil {
		return nil, err
	}
	pageReq.UserType = string(constants.UserTypeBorrower)
	pageReq.Filter.OrganizationType = string(constants.UserTypeLender)
	pageReq.Filter.Status = string(constants.OrganizationStatusActive)

	pageResult, err := l.svcCtx.TenantClient.Organization.Page(l.ctx, &tenant.OrganizationPageRequest{Data: pageReq})
	if err != nil {
		return nil, err
	}

	data, err := copierx.Trans[tenant.OrganizationPageResponseData, borrower.OrganizationLenderPageResponseData](pageResult.Data, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}

	return &borrower.OrganizationLenderPageResponse{
		Data: data,
	}, nil
}
