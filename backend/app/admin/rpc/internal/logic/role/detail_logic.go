package rolelogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/domains/zeebee/constants"
)

type DetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *DetailLogic) Detail(in *admin.RoleDetailRequest) (resp *admin.RoleDetailResponse, err error) {

	claim, err := l.baseLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	role, err := l.svcCtx.TenantClient.Role.Detail(l.ctx, &tenant.RoleDetailRequest{Data: &tenant.RoleDetailRequestData{
		UserType:       string(constants.UserTypeAdmin),
		OrganizationId: claim.OrganizationId,
		Id:             in.Data.Id,
	}})
	if err != nil {
		return nil, err
	}

	return &admin.RoleDetailResponse{
		Data: role.Data,
	}, nil
}
