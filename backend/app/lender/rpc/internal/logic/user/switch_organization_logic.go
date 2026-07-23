package userlogic

import (
	"context"

	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/proto"
)

type SwitchOrganizationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewSwitchOrganizationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SwitchOrganizationLogic {
	return &SwitchOrganizationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *SwitchOrganizationLogic) SwitchOrganization(in *proto.SwitchOrganizationRequest) (resp *proto.UserLoginResponse, err error) {

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	user, err := l.baseLogic.UserLogic.SwitchOrganization(l.ctx, claim.UserId, in.Data.OrganizationId, constants.UserTypeLender)
	if err != nil {
		return nil, err
	}

	token, err := l.baseLogic.UserLogic.GenerateJwtTokenById(l.ctx, user.Id, user.DefaultOrganizationId, constants.UserTypeLender)
	if err != nil {
		return nil, err
	}
	data, err := copierx.Trans[dto.UserJwtTokenDto, proto.UserLoginResponseData](token)
	if err != nil {
		return nil, err
	}

	return &proto.UserLoginResponse{
		Data: data,
	}, nil
}
