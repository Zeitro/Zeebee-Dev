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

type AllLenderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewAllLenderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AllLenderLogic {
	return &AllLenderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *AllLenderLogic) AllLender(in *proto.EmptyRequest) (*borrower.OrganizationAllLenderResponse, error) {
	result, err := l.svcCtx.TenantClient.Organization.AllForOption(l.ctx, &tenant.OrganizationAllForOptionRequest{Data: &tenant.OrganizationAllForOptionRequestData{
		UserType: string(constants.UserTypeLender),
		Status:   string(constants.OrganizationStatusActive),
	}})
	if err != nil {
		return nil, err
	}

	return &borrower.OrganizationAllLenderResponse{
		Data: result.Data,
	}, nil
}
