package apilogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/domains/zeebee/constants"
)

type DetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{ctx: ctx, svcCtx: svcCtx}
}

// Detail proxies to tenant; tenant is the RBAC single authority.
func (l *DetailLogic) Detail(in *admin.ApiDetailRequest) (resp *admin.ApiDetailResponse, err error) {
	tenantResp, err := l.svcCtx.TenantClient.Api.Detail(l.ctx, &tenant.ApiDetailRequest{
		Data: &tenant.ApiDetailRequestData{
			UserType: string(constants.UserTypeAdmin),
			Id:       in.Data.Id,
			Type:     in.Data.Type,
		},
	})
	if err != nil {
		return nil, err
	}
	return &admin.ApiDetailResponse{Data: tenantResp.Data}, nil
}
