package rolelogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
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

func (l *PageLogic) Page(in *admin.RolePageRequest) (resp *admin.RolePageResponse, err error) {

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	pageReq, err := copierx.Trans[admin.RolePageRequestData, tenant.RolePageRequestData](in.Data)
	if err != nil {
		return nil, err
	}

	pageReq.UserType = string(constants.UserTypeAdmin)
	pageReq.OrganizationId = claim.OrganizationId

	result, err := l.svcCtx.TenantClient.Role.Page(l.ctx, &tenant.RolePageRequest{Data: pageReq})
	if err != nil {
		return nil, err
	}

	data, err := copierx.Trans[tenant.RolePageResponseData, admin.RolePageResponseData](result.Data, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}

	return &admin.RolePageResponse{
		Data: data,
	}, nil
}
