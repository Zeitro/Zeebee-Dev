package rolelogic

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

func (l *PageLogic) Page(in *proto.CommonPageRequest) (resp *borrower.RolePageResponse, err error) {

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	pageReq, err := copierx.Trans[proto.CommonPageRequestData, tenant.RolePageRequestData](in.Data)
	if err != nil {
		return nil, err
	}
	pageReq.UserType = string(constants.UserTypeBorrower)
	pageReq.OrganizationId = claim.OrganizationId

	result, err := l.svcCtx.TenantClient.Role.Page(l.ctx, &tenant.RolePageRequest{Data: pageReq})
	if err != nil {
		return nil, err
	}
	data, err := copierx.Trans[tenant.RolePageResponseData, borrower.RolePageResponseData](result.Data, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}

	return &borrower.RolePageResponse{
		Data: data,
	}, nil
}
