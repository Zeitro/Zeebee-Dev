package userlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/global"
	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/proto"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *LoginLogic) Login(in *proto.UserLoginRequest) (resp *proto.UserLoginResponse, err error) {

	loginReq := dto.UserLoginRequestDto{
		Email:    in.Data.Email,
		Password: in.Data.Password,
		UserType: constants.UserTypeBorrower,
	}
	limiter := dto.LoginLimiterDto{
		LimiterItem: global.Config.LimiterConf.Login,
		RedisCli:    global.GetRedis(),
	}
	token, err := l.baseLogic.UserLogic.Login(l.ctx, loginReq, limiter)
	if err != nil {
		return nil, err
	}
	data := &proto.UserLoginResponseData{}
	data, err = copierx.Trans[dto.UserJwtTokenDto, proto.UserLoginResponseData](token)
	if err != nil {
		return nil, err
	}
	return &proto.UserLoginResponse{
		Data: data,
	}, nil
}
