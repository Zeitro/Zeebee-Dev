package organizationlogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/errorx"
)

type CreateLenderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateLenderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLenderLogic {
	return &CreateLenderLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *CreateLenderLogic) CreateLender(in *admin.OrganizationCreateLenderRequest) (*admin.OrganizationDetailResponse, error) {
	if in == nil || in.Data == nil || in.Data.Owner == nil || in.Data.Organization == nil {
		return nil, errorx.New("request data is required")
	}
	owner, organization := in.Data.Owner, in.Data.Organization
	result, err := l.svcCtx.TenantClient.Organization.Create(l.ctx, &tenant.OrganizationCreateRequest{Data: &tenant.OrganizationCreateRequestData{
		UserType:       "lender",
		Name:           organization.Name,
		ContactEmail:   organization.ContactEmail,
		ContactPhone:   organization.ContactPhone,
		ContactName:    organization.ContactName,
		Logo:           organization.Logo,
		Website:        organization.Website,
		Tagline:        organization.Tagline,
		Metadata:       organization.Metadata,
		Address:        organization.Address,
		Status:         "active",
		OwnerType:      "third_party",
		Source:         "manual",
		OwnerEmail:     owner.Email,
		OwnerPassword:  owner.Password,
		OwnerFirstName: owner.FirstName,
		OwnerLastName:  owner.LastName,
	}})
	if err != nil {
		return nil, err
	}
	return &admin.OrganizationDetailResponse{Data: result.Data}, nil
}
