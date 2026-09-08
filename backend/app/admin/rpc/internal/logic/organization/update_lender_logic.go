package organizationlogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/errorx"
)

type UpdateLenderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateLenderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLenderLogic {
	return &UpdateLenderLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *UpdateLenderLogic) UpdateLender(in *admin.OrganizationUpdateLenderRequest) (*admin.OrganizationDetailResponse, error) {
	if in == nil || in.Data == nil {
		return nil, errorx.ErrorEmptyRequestData
	}
	data := in.Data
	request := &tenant.OrganizationUpdateLenderRequest{Data: &tenant.OrganizationUpdateLenderRequestData{
		UserType:     "admin",
		Id:           data.Id,
		ContactEmail: data.ContactEmail,
		ContactPhone: data.ContactPhone,
		ContactName:  data.ContactName,
		Logo:         data.Logo,
		Website:      data.Website,
		Tagline:      data.Tagline,
		Metadata:     data.Metadata,
		Address:      data.Address,
	}}
	if data.Owner != nil {
		request.Data.OwnerEmail = data.Owner.Email
		request.Data.OwnerPassword = data.Owner.Password
		request.Data.OwnerFirstName = data.Owner.FirstName
		request.Data.OwnerLastName = data.Owner.LastName
	}
	result, err := l.svcCtx.TenantClient.Organization.UpdateLender(l.ctx, request)
	if err != nil {
		return nil, err
	}
	return &admin.OrganizationDetailResponse{Data: result.Data}, nil
}
