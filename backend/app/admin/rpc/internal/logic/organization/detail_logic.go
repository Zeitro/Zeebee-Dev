package organizationlogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/proto"
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

// Detail retrieves organization details by ID.
// It fetches the organization with preloaded relations (Owner, OrganizationEmails).
func (l *DetailLogic) Detail(in *proto.IdRequest) (resp *admin.OrganizationDetailResponse, err error) {

	organization, err := l.svcCtx.TenantClient.Organization.Detail(l.ctx, &tenant.OrganizationDetailRequest{Data: &tenant.OrganizationDetailRequestData{
		UserType: string(constants.UserTypeAdmin),
		Id:       in.Data,
	}})
	if err != nil {
		return nil, err
	}
	return &admin.OrganizationDetailResponse{
		Data: organization.Data,
	}, nil
}
