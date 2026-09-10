package organizationlogic

import (
	"context"

	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/copierx"
)

type CreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLogic) Create(in *lender.OrganizationCreateRequest) (resp *lender.OrganizationDetailResponse, err error) {

	data, err := copierx.Trans[lender.OrganizationCreateRequestData, tenant.OrganizationCreateRequestData](in.Data)
	if err != nil {
		return nil, err
	}
	data.UserType = "lender"
	result, err := l.svcCtx.TenantClient.Organization.Create(l.ctx, &tenant.OrganizationCreateRequest{Data: data})
	if err != nil {
		return nil, err
	}
	return &lender.OrganizationDetailResponse{Data: result.Data}, nil
}
