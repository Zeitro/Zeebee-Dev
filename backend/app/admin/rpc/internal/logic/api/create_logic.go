package apilogic

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
func (l *CreateLogic) Create(in *admin.ApiCreateRequest) (resp *admin.ApiDetailResponse, err error) {
	data, err := copierx.Trans[admin.ApiCreateRequestData, tenant.ApiCreateRequestData](in.Data)
	if err != nil {
		return nil, err
	}
	data.UserType = string(constants.UserTypeAdmin)
	tenantResp, err := l.svcCtx.TenantClient.Api.Create(l.ctx, &tenant.ApiCreateRequest{Data: data})
	if err != nil {
		return nil, err
	}
	return &admin.ApiDetailResponse{Data: tenantResp.Data}, nil
}
