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

type ResetPasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewResetPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetPasswordLogic {
	return &ResetPasswordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *ResetPasswordLogic) ResetPassword(in *proto.ResetPasswordRequest) (resp *proto.UserLoginResponse, err error) {

	input := dto.UserResetPasswordByEmailRequestDto{
		UserType:    constants.UserTypeLender,
		NewPassword: in.Data.Password,
		Email:       in.Data.Email,
		Code:        in.Data.Code,
	}

	token, err := l.baseLogic.UserLogic.ResetPasswordWithEmail(l.ctx, input)
	if err != nil {
		return nil, err
	}

	jwtInfo, err := copierx.Trans[dto.UserJwtTokenDto, proto.UserLoginResponseData](token)
	if err != nil {
		return nil, err
	}

	return &proto.UserLoginResponse{
		Data: jwtInfo,
	}, nil
}
