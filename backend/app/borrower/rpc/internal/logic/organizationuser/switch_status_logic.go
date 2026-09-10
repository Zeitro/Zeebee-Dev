package organizationuserlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/proto"
)

type SwitchStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewSwitchStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SwitchStatusLogic {
	return &SwitchStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *SwitchStatusLogic) SwitchStatus(in *borrower.OrganizationUserUpdateStatusRequest) (resp *proto.EmptyResponse, err error) {

	claim, err := l.baseLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	return l.svcCtx.TenantClient.OrganizationUser.SwitchStatus(l.ctx, &tenant.OrganizationUserSwitchStatusRequest{Data: &tenant.OrganizationUserSwitchStatusRequestData{
		UserType:       string(constants.UserTypeBorrower),
		OrganizationId: claim.OrganizationId,
		Id:             in.Data.Id,
		Status:         in.Data.Status,
	}})
}
