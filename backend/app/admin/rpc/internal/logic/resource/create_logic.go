package resourcelogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/constants"
)

type CreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{ctx: ctx, svcCtx: svcCtx}
}

// Create proxies to tenant; tenant is the RBAC single authority.
func (l *CreateLogic) Create(in *admin.ResourceCreateRequest) (resp *admin.ResourceDetailResponse, err error) {
	data, err := copierx.Trans[admin.ResourceCreateRequestData, tenant.ResourceCreateRequestData](in.Data)
	if err != nil {
		return nil, err
	}
	data.UserType = string(constants.UserTypeAdmin)
	tenantResp, err := l.svcCtx.TenantClient.Resource.Create(l.ctx, &tenant.ResourceCreateRequest{Data: data})
	if err != nil {
		return nil, err
	}
	return &admin.ResourceDetailResponse{Data: tenantResp.Data}, nil
}
