package userlogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/lender/rpc/internal/global"
	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
)

type LoginLogic struct {
	baseLogic *logic.BaseLogic
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.GetBaseLogic(svcCtx),
		Logger:    logx.WithContext(ctx),
	}
}

// user
func (l *LoginLogic) Login(in *lender.UserLoginRequest) (resp *lender.UserJwtResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
				logx.Field("service_name", l.svcCtx.Config.Name),
			).Error(err)
		}
	}()
	loginReq := dto.UserLoginRequestDto{
		Email:    in.Data.Email,
		Password: in.Data.Password,
		UserType: constants.UserTypeLender,
	}
	limiter := dto.LoginLimiterDto{
		Rate:     global.Config.LimiterConf.Login.Rate,
		Period:   global.Config.LimiterConf.Login.Period,
		RedisCli: global.GetRedis(),
	}
	token, err := l.baseLogic.UserLogic.Login(l.ctx, loginReq, limiter)
	if err != nil {
		return nil, err
	}

	data, err := copierx.Trans[dto.UserJwtTokenDto, lender.UserJwtResponseData](token)
	if err != nil {
		return nil, err
	}

	return &lender.UserJwtResponse{
		Data: data,
	}, nil
}
