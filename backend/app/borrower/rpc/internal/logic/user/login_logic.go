package userlogic

import (
	"context"
	"strconv"

	"github.com/pkg/errors"

	innerconstants "zeebee-backend/app/borrower/rpc/internal/constants"
	"zeebee-backend/app/borrower/rpc/internal/global"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
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

func (l *LoginLogic) Login(in *borrower.UserLoginRequest) (resp *borrower.UserLoginResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	// check input wrong pwd times
	redisKey := stringx.Formate(innerconstants.BorrowerLoginFailedTimes, in.Email)
	var loginFailedTimes int
	if global.Config.MaxFailedLoginTimes > 0 {
		loginFailedTimes, err = global.GetRedis().GetInt(l.ctx, redisKey)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		if loginFailedTimes >= global.Config.MaxFailedLoginTimes {
			return nil, errors.WithStack(errors.New("login failed times exceeded, please try again later"))
		}
	}

	user, err := l.svcCtx.UserRepo.FindOneByColumns(l.ctx,
		map[string]interface{}{
			"email": in.Email,
			"type":  constants.UserTypeBorrower,
		},
	)
	// check pwd
	if !utils.CheckPassword(in.Password, user.PasswordHash) {
		// add input wrong pwd times to redis
		_, err = global.GetRedis().IncreaseExpire(l.ctx, redisKey, global.Config.ForbiddenLoginSeconds)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		return nil, errors.New("password is incorrect")
	}

	if global.Config.MaxFailedLoginTimes > 0 {
		err = global.GetRedis().Del(l.ctx, redisKey)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}

	// check status
	if user.Status != constants.UserStatusActive {
		return nil, errors.WithStack(errors.New("user is not active"))
	}
	// generate token
	accessToken, err := utils.GenerateToken(strconv.Itoa(int(user.ID)), user.Email, "", l.svcCtx.Config.JwtAuth.AccessSecret, l.svcCtx.Config.JwtAuth.AccessExpire)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// refresh token
	refreshToken, err := utils.GenerateToken(strconv.Itoa(int(user.ID)), user.Email, "", l.svcCtx.Config.JwtAuth.RefreshSecret, l.svcCtx.Config.JwtAuth.RefreshExpire)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &borrower.UserLoginResponse{
		AccessToken:        accessToken,
		RefreshToken:       refreshToken,
		AccessTokenExpire:  l.svcCtx.Config.JwtAuth.AccessExpire,
		RefreshTokenExpire: l.svcCtx.Config.JwtAuth.RefreshExpire,
	}, nil
}
