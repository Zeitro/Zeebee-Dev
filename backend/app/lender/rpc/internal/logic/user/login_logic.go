package userlogic

import (
	"context"
	"strconv"

	"github.com/pkg/errors"

	"zeebee-backend/app/lender/rpc/internal/global"
	"zeebee-backend/app/lender/rpc/internal/repo/model"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/constants"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/stringx"
	"zeebee-backend/common/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// user
func (l *LoginLogic) Login(in *lender.UserLoginRequest) (resp *lender.UserLoginResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	// check input wrong pwd times
	redisKey := stringx.Format(constants.LenderLoginFailedTimes, in.Data.Email)
	var loginFailedTimes int
	if global.Config.LimiterConf.Login.Rate > 0 {
		loginFailedTimes, err = global.GetRedis().GetInt(l.ctx, redisKey)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		if loginFailedTimes >= global.Config.LimiterConf.Login.Rate {
			return nil, errors.WithStack(errors.New("Login failed times exceeded, please try again later"))
		}
	}

	user, err := l.svcCtx.UserRepo.GetLenderUserByEmail(l.ctx, in.Data.Email)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if user.LenderEmployee == nil {
		return nil, errors.New("User not found.")
	}

	// check pwd
	if !utils.CheckPassword(in.Data.Password, user.PasswordHash) {
		// add input wrong pwd times to redis
		_, err = global.GetRedis().IncreaseExpire(l.ctx, redisKey, global.Config.LimiterConf.Login.Period)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		return nil, errors.WithStack(errors.New("The email and password you entered don't match."))
	}

	if global.Config.LimiterConf.Login.Rate > 0 {
		err = global.GetRedis().Del(l.ctx, redisKey)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}

	// check status
	if user.Status == constants.UserStatusInactive {
		return nil, errors.WithStack(errors.New("User is not active"))
	}
	if user.Status == constants.UserStatusUnverified {
		return nil, errors.WithStack(errors.New("User is unverified"))
	}

	// token
	token, err := generateLenderToken(user, l.svcCtx)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &lender.UserLoginResponse{
		Data: &lender.UserLoginResponseData{
			AccessToken:        token.AccessToken,
			RefreshToken:       token.RefreshToken,
			AccessTokenExpire:  token.AccessTokenExpire,
			RefreshTokenExpire: token.RefreshTokenExpire,
		},
	}, nil
}

func generateLenderToken(user *model.User, svcCtx *svc.ServiceContext) (token *lender.UserLoginResponseData, err error) {
	if user.LenderEmployee == nil {
		return nil, errors.WithStack(errors.New("LenderEmployee is empty."))
	}
	accessToken, err := utils.GenerateToken(
		strconv.Itoa(int(user.Id)),
		strconv.Itoa(int(user.LenderEmployee.Id)),
		user.LenderEmployee.LenderCompanyId,
		user.LenderEmployee.FullName,
		"",
		string(constants.UserTypeLender),
		svcCtx.Config.JwtAuth.AccessSecret,
		svcCtx.Config.JwtAuth.AccessExpire,
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	refreshToken, err := utils.GenerateToken(
		strconv.Itoa(int(user.Id)),
		strconv.Itoa(int(user.LenderEmployee.Id)),
		user.LenderEmployee.LenderCompanyId,
		user.LenderEmployee.FullName,
		"",
		string(constants.UserTypeLender),
		svcCtx.Config.JwtAuth.RefreshSecret,
		svcCtx.Config.JwtAuth.RefreshExpire,
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	token = &lender.UserLoginResponseData{
		AccessToken:        accessToken,
		RefreshToken:       refreshToken,
		AccessTokenExpire:  svcCtx.Config.JwtAuth.AccessExpire,
		RefreshTokenExpire: svcCtx.Config.JwtAuth.RefreshExpire,
	}
	return

}
