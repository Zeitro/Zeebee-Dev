package userlogic

import (
	"context"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/admin/rpc/internal/global"
	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

// user
func (l *LoginLogic) Login(in *admin.UserLoginRequest) (resp *admin.UserLoginResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	loginReq := dto.UserLoginRequestDto{
		Email:    in.Data.Email,
		Password: in.Data.Password,
		UserType: constants.UserTypeAdmin,
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
	data := &admin.UserLoginResponseData{}
	err = copier.CopyWithOption(&data, token, copier.Option{CaseSensitive: true})
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	return &admin.UserLoginResponse{
		Data: data,
	}, nil
}
