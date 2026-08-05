package organizationuserlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
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

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx, constants.UserTypeBorrower)
	if err != nil {
		return nil, err
	}
	err = l.baseLogic.OrganizationUserLogic.SwitchStatus(l.ctx, in.Data.Id, claim.OrganizationId, constants.OrganizationUserStatus(in.Data.Status), claim.UserType)
	if err != nil {
		return nil, err
	}

	return &proto.EmptyResponse{}, nil
}
