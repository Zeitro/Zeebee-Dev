package organizationlogic

import (
	"context"

	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/proto"
)

type UpdateProfileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateProfileLogic {
	return &UpdateProfileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateProfileLogic) UpdateProfile(in *lender.OrganizationUpdateProfileRequest) (resp *lender.OrganizationProfileResponse, err error) {

	data, err := copierx.Trans[lender.OrganizationUpdateProfileRequestData, tenant.OrganizationUpdateProfileRequestData](in.Data)
	if err != nil {
		return nil, err
	}
	data.UserType = "lender"
	_, err = l.svcCtx.TenantClient.Organization.UpdateProfile(l.ctx, &tenant.OrganizationUpdateProfileRequest{Data: data})
	if err != nil {
		return nil, err
	}
	profileLogic := NewProfileLogic(l.ctx, l.svcCtx)
	return profileLogic.Profile(&proto.EmptyRequest{})
}
