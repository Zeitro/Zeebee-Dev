package userlogic

import (
	"context"

	"github.com/jinzhu/copier"

	"zeebee-backend/app/borrower/rpc/internal/global"
	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/constants"
	"zeebee-backend/common/stringx"
	logic2 "zeebee-backend/domains/zeebee/logic"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"

	"zeebee-backend/common/utils"
)

type ResetPasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewResetPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetPasswordLogic {
	return &ResetPasswordLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.NewBaseLogic(ctx, svcCtx),
	}
}

func (l *ResetPasswordLogic) ResetPassword(in *borrower.ResetPasswordRequest) (resp *borrower.ResetPasswordResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	email := in.Data.Email
	code := in.Data.Code
	newPwd := in.Data.Password
	if email == "" || code == "" || newPwd == "" {
		return nil, errorx.New("Email, Code And New Password Are Required")
	}
	redisKey := stringx.Format(constants.BorrowerResetPwdCode, email)
	cacheToken, err := global.GetRedis().Get(l.ctx, redisKey)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	if cacheToken == "" || cacheToken != code {
		return nil, errorx.WithStack(errorx.New("Invalid Or Expired Code"))
	}

	user, err := l.svcCtx.UserRepo.GetUserByEmail(l.ctx, email, constants.UserTypeBorrower)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errorx.New("User Not Found")
	}

	// hash pwd
	pwdHash, err := utils.HashPassword(newPwd)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	// update pwd
	err = l.svcCtx.UserRepo.UpdateUserPassword(l.ctx, user.Id, pwdHash)
	if err != nil {
		return nil, err
	}

	user, err = l.baseLogic.UserLogic.CheckAndUpdateActiveOrganizationByUser(l.ctx, user)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	token, err := l.baseLogic.UserLogic.GenerateJwtToken(user, logic2.JwtConfig{
		AccessSecret:  l.svcCtx.Config.JwtAuth.AccessSecret,
		AccessExpire:  l.svcCtx.Config.JwtAuth.AccessExpire,
		RefreshSecret: l.svcCtx.Config.JwtAuth.RefreshSecret,
		RefreshExpire: l.svcCtx.Config.JwtAuth.RefreshExpire,
	}, constants.UserTypeBorrower)
	if err != nil {
		return nil, err
	}
	data := &borrower.UserLoginResponseData{}
	err = copier.CopyWithOption(&data, token, copier.Option{CaseSensitive: true})
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	// delete redis key
	_ = global.GetRedis().Del(l.ctx, redisKey)
	return &borrower.ResetPasswordResponse{
		Data: data,
	}, nil
}
