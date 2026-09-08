package apilogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/proto"
)

type DeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogic {
	return &DeleteLogic{ctx: ctx, svcCtx: svcCtx}
}

// Delete proxies to tenant; tenant is the RBAC single authority.
func (l *DeleteLogic) Delete(in *admin.ApiDeleteRequest) (resp *proto.EmptyResponse, err error) {
	return l.svcCtx.TenantClient.Api.Delete(l.ctx, &tenant.ApiDeleteRequest{
		Data: &tenant.ApiDeleteRequestData{
			UserType: string(constants.UserTypeAdmin),
			Id:       in.Data.Id,
			Type:     in.Data.Type,
		},
	})
}
