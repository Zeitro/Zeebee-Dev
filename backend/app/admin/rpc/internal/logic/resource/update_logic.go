package resourcelogic

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
func (l *UpdateLogic) Update(in *admin.ResourceUpdateRequest) (resp *admin.ResourceDetailResponse, err error) {
	data, err := copierx.Trans[admin.ResourceUpdateRequestData, tenant.ResourceUpdateRequestData](in.Data)
	if err != nil {
		return nil, err
	}
	data.UserType = string(constants.UserTypeAdmin)
	tenantResp, err := l.svcCtx.TenantClient.Resource.Update(l.ctx, &tenant.ResourceUpdateRequest{Data: data})
	if err != nil {
		return nil, err
	}
	return &admin.ResourceDetailResponse{Data: tenantResp.Data}, nil
}
