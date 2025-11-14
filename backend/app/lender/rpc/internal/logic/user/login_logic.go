package userlogic

import (
	"context"
	"strconv"

	"zeebee-backend/app/lender/rpc/internal/logic"

	"zeebee-backend/app/lender/rpc/internal/global"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/constants"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/stringx"
	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/model"

	"github.com/zeromicro/go-zero/core/logx"
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
		baseLogic: logic.NewBaseLogic(ctx, svcCtx),
		Logger:    logx.WithContext(ctx),
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
			return nil, errorx.WithStack(err)
		}
		if loginFailedTimes >= global.Config.LimiterConf.Login.Rate {
			return nil, errorx.WithStack(errorx.New("Login Failed Times Exceeded, Please Try Again Later"))
		}
	}

	user, err := l.svcCtx.UserRepo.GetLenderUserByEmail(l.ctx, in.Data.Email)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	if user.LenderEmployee == nil {
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
		return nil, errorx.WithStack(errorx.New("User Is Not Active"))
	}
	if user.Status == constants.UserStatusUnverified {
		return nil, errorx.WithStack(errorx.New("User Is Unverified"))
	}

	// check company
	if user.LenderEmployee == nil {
		return nil, errorx.WithStack(errorx.New("LenderEmployee Is Empty."))
	}

	employee, err := l.baseLogic.CheckAndUpdateActiveLenderCompanyByEmployee(user.LenderEmployee)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	user.LenderEmployee = employee

	// token
	token, err := generateLenderToken(user, l.svcCtx)
	if err != nil {
		return nil, errorx.WithStack(err)
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
		return nil, errorx.WithStack(errorx.New("LenderEmployee Is Empty."))
	}
	accessToken, err := utils.GenerateToken(
		strconv.Itoa(int(user.Id)),
		strconv.Itoa(int(user.LenderEmployee.Id)),
		user.LenderEmployee.DefaultLenderCompanyId,
		user.LenderEmployee.FullName,
		"",
		string(constants.UserTypeLender),
		svcCtx.Config.JwtAuth.AccessSecret,
		svcCtx.Config.JwtAuth.AccessExpire,
	)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	refreshToken, err := utils.GenerateToken(
		strconv.Itoa(int(user.Id)),
		strconv.Itoa(int(user.LenderEmployee.Id)),
		user.LenderEmployee.DefaultLenderCompanyId,
		user.LenderEmployee.FullName,
		"",
		string(constants.UserTypeLender),
		svcCtx.Config.JwtAuth.RefreshSecret,
		svcCtx.Config.JwtAuth.RefreshExpire,
	)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	token = &lender.UserLoginResponseData{
		AccessToken:        accessToken,
		RefreshToken:       refreshToken,
		AccessTokenExpire:  svcCtx.Config.JwtAuth.AccessExpire,
		RefreshTokenExpire: svcCtx.Config.JwtAuth.RefreshExpire,
	}
	return
}
