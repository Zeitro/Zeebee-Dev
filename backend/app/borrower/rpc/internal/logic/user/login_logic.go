package userlogic

import (
	"context"
	"strconv"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"zeebee-backend/app/borrower/rpc/internal/global"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/constants"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/stringx"
	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/model"

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
	redisKey := stringx.Format(constants.BorrowerLoginFailedTimes, in.Data.Email)
	var loginFailedTimes int
	if global.Config.LimiterConf.Login.Rate > 0 {
		loginFailedTimes, err = global.GetRedis().GetInt(l.ctx, redisKey)
		if err != nil {
			return nil, errorx.WithStack(err)
		}
		if loginFailedTimes >= global.Config.LimiterConf.Login.Rate {
			return nil, errorx.New("Login Failed Times Exceeded, Please Try Again Later")
		}
	}

	user, err := l.svcCtx.UserRepo.GetBorrowerUserByEmail(l.ctx, in.Data.Email)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	if user.Borrower == nil {
		return nil, errorx.New("User Not Found.")
	}

	// check pwd
	if !utils.CheckPassword(in.Data.Password, user.PasswordHash) {
		// add input wrong pwd times to redis
		_, err = global.GetRedis().IncreaseExpire(l.ctx, redisKey, global.Config.LimiterConf.Login.Period)
		if err != nil {
			return nil, errorx.WithStack(err)
		}
		return nil, errorx.WithStack(errorx.New("The Email And Password You Entered Don't Match."))
	}

	if global.Config.LimiterConf.Login.Rate > 0 {
		err = global.GetRedis().Del(l.ctx, redisKey)
		if err != nil {
			return nil, errorx.WithStack(err)
		}
	}

	// check status
	if user.Status == constants.UserStatusInactive {
		return nil, errorx.New("User Is Not Active")
	}
	if user.Status == constants.UserStatusUnverified {
		err = status.Error(codes.PermissionDenied, "User Is Unverified")
		return nil, err
	}

	// generate token
	accessToken, err := utils.GenerateToken(
		strconv.Itoa(int(user.Id)),
		strconv.Itoa(int(user.Borrower.Id)),
		0,
		user.Email,
		"",
		string(constants.UserTypeBorrower),
		l.svcCtx.Config.JwtAuth.AccessSecret,
		l.svcCtx.Config.JwtAuth.AccessExpire,
	)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	// refresh token
	refreshToken, err := utils.GenerateToken(
		strconv.Itoa(int(user.Id)),
		strconv.Itoa(int(user.Borrower.Id)),
		0,
		user.Email,
		"",
		string(constants.UserTypeBorrower),
		l.svcCtx.Config.JwtAuth.RefreshSecret,
		l.svcCtx.Config.JwtAuth.RefreshExpire,
	)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	return &borrower.UserLoginResponse{
		Data: &borrower.UserLoginResponseData{
			AccessToken:        accessToken,
			RefreshToken:       refreshToken,
			AccessTokenExpire:  l.svcCtx.Config.JwtAuth.AccessExpire,
			RefreshTokenExpire: l.svcCtx.Config.JwtAuth.RefreshExpire,
		},
	}, nil
}

func generateBorrowerToken(user *model.User, svcCtx *svc.ServiceContext) (token *borrower.UserLoginResponseData, err error) {
	if user.Borrower == nil {
		return nil, errorx.WithStack(errorx.New("Borrower Is Empty."))
	}
	accessToken, err := utils.GenerateToken(
		strconv.Itoa(int(user.Id)),
		strconv.Itoa(int(user.Borrower.Id)),
		0,
		user.Email,
		"",
		string(constants.UserTypeBorrower),
		svcCtx.Config.JwtAuth.AccessSecret,
		svcCtx.Config.JwtAuth.AccessExpire,
	)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	refreshToken, err := utils.GenerateToken(
		strconv.Itoa(int(user.Id)),
		strconv.Itoa(int(user.Borrower.Id)),
		0,
		user.Email,
		"",
		string(constants.UserTypeBorrower),
		svcCtx.Config.JwtAuth.RefreshSecret,
		svcCtx.Config.JwtAuth.RefreshExpire,
	)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	token = &borrower.UserLoginResponseData{
		AccessToken:        accessToken,
		RefreshToken:       refreshToken,
		AccessTokenExpire:  svcCtx.Config.JwtAuth.AccessExpire,
		RefreshTokenExpire: svcCtx.Config.JwtAuth.RefreshExpire,
	}
	return

}
