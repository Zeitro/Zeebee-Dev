package apilogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/constants"
)

type UpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{ctx: ctx, svcCtx: svcCtx}
}

// Update proxies to tenant; tenant is the RBAC single authority.
func (l *UpdateLogic) Update(in *admin.ApiUpdateRequest) (resp *admin.ApiDetailResponse, err error) {
	data, err := copierx.Trans[admin.ApiUpdateRequestData, tenant.ApiUpdateRequestData](in.Data)
	if err != nil {
		return nil, err
	}
	data.UserType = string(constants.UserTypeAdmin)
	tenantResp, err := l.svcCtx.TenantClient.Api.Update(l.ctx, &tenant.ApiUpdateRequest{Data: data})
	if err != nil {
		return nil, err
	}
	return &admin.ApiDetailResponse{Data: tenantResp.Data}, nil
}
