package userlogic

import (
	"context"

	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	tenantpb "zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/proto"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{ctx: ctx, svcCtx: svcCtx}
}

// Register proxies identity creation to tenant. Tenant owns duplicate detection,
// verify-email send, and org seeding; lender only forwards the form fields.
func (l *RegisterLogic) Register(in *lender.UserRegisterRequest) (resp *proto.EmptyResponse, err error) {
	_, err = l.svcCtx.TenantClient.User.Register(l.ctx, &tenantpb.RegisterRequest{
		Data: &tenantpb.RegisterRequestData{
			UserType:         string(constants.UserTypeLender),
			Email:            in.Data.Email,
			Password:         in.Data.Password,
			FirstName:        in.Data.FirstName,
			LastName:         in.Data.LastName,
			OrganizationName: in.Data.OrganizationName,
			OwnerType:        string(constants.OrganizationOwnerTypeThirdParty),
		},
	})
	if err != nil {
		return nil, err
	}
	return &proto.EmptyResponse{}, nil
}
