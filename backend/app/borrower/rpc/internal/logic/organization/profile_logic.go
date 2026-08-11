package organizationlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/proto"
)

type ProfileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProfileLogic {
	return &ProfileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *ProfileLogic) Profile(in *proto.EmptyRequest) (resp *borrower.OrganizationInfoResponse, err error) {

	claim, err := l.baseLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	organization, err := l.svcCtx.TenantClient.Organization.DetailWithType(l.ctx, &tenant.OrganizationDetailWithTypeRequest{Data: &tenant.OrganizationDetailWithTypeRequestData{
		UserType: string(constants.UserTypeBorrower),
		Id:       claim.OrganizationId,
	}})
	if err != nil {
		return nil, err
	}
	return &borrower.OrganizationInfoResponse{
		Data: organization.Data,
	}, nil
}
