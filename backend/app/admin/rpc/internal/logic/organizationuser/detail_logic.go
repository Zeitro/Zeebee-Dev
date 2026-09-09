package organizationuserlogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
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

func (l *DetailLogic) Detail(in *admin.OrganizationUserDetailRequest) (resp *admin.OrganizationUserDetailResponse, err error) {

	claim, err := l.baseLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	organizationUser, err := l.svcCtx.TenantClient.OrganizationUser.Detail(l.ctx, &tenant.OrganizationUserDetailRequest{Data: &tenant.OrganizationUserDetailRequestData{
		UserType:       string(claim.UserType),
		OrganizationId: claim.OrganizationId,
		Id:             in.Data.Id,
	}})
	if err != nil {
		return nil, err
	}

	return &admin.OrganizationUserDetailResponse{
		Data: organizationUser.Data,
	}, nil
}
