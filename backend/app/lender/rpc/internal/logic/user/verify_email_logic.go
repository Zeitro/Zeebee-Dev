package userlogic

import (
	"context"

	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/proto"
)

type VerifyEmailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewVerifyEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerifyEmailLogic {
	return &VerifyEmailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *VerifyEmailLogic) VerifyEmail(in *proto.VerifyEmailRequest) (resp *proto.UserLoginResponse, err error) {

	var user *model.User
	email := in.Data.Email
	code := in.Data.Code
	user, err = l.baseLogic.UserLogic.VerifyUserEmail(l.ctx, email, code, constants.UserTypeLender)
	if err != nil {
		return nil, err
	}
	organizationUser, err := l.svcCtx.OrganizationUserRepo.GetOrganizationUserByUserId(l.ctx, user.DefaultOrganizationId, user.Id)
	if err != nil {
		return nil, err
	}

	tokenInfo, err := l.baseLogic.UserLogic.GenerateJwtToken(user, organizationUser)
	if err != nil {
		return nil, err
	}

	data, err := copierx.Trans[dto.UserJwtTokenDto, proto.UserLoginResponseData](tokenInfo)
	if err != nil {
		return nil, err
	}
	return &proto.UserLoginResponse{
		Data: data,
	}, nil
}
