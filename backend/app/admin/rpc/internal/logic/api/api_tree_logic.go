package apilogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/domains/zeebee/constants"
)

type ApiTreeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApiTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiTreeLogic {
	return &ApiTreeLogic{ctx: ctx, svcCtx: svcCtx}
}

// ApiTree proxies to tenant; tenant is the RBAC single authority.
func (l *ApiTreeLogic) ApiTree(in *admin.ApiTreeRequest) (resp *admin.ApiTreeResponse, err error) {
	tenantResp, err := l.svcCtx.TenantClient.Api.ApiTree(l.ctx, &tenant.ApiTreeRequest{
		Data: &tenant.ApiTreeRequestData{
			UserType: string(constants.UserTypeAdmin),
			Type:     in.Data.Type,
		},
	})
	if err != nil {
		return nil, err
	}
	return &admin.ApiTreeResponse{Data: tenantResp.Data}, nil
}
