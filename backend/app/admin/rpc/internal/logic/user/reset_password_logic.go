package userlogic

import (
	"context"

	"github.com/jinzhu/copier"

	"zeebee-backend/app/admin/rpc/internal/global"
	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/constants"
	"zeebee-backend/common/stringx"
	logic2 "zeebee-backend/domains/zeebee/logic"
	"zeebee-backend/domains/zeebee/model"

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

// rpc VerifyEmail(VerifyEmailRequest) returns(VerifyEmailResponse);
func (l *ResetPasswordLogic) ResetPassword(in *admin.ResetPasswordRequest) (resp *admin.ResetPasswordResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	var orgUser *model.OrganizationUser

	email := in.Data.Email
	code := in.Data.Code
	newPwd := in.Data.Password
	if email == "" || code == "" || newPwd == "" {
		return nil, errorx.New("Email, Code And New Password Are Required")
	}
	redisKey := stringx.Format(constants.AdminResetPwdCode, email)
	cacheToken, err := global.GetRedis().Get(l.ctx, redisKey)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	if cacheToken == "" || cacheToken != code {
		return nil, errorx.WithStack(errorx.New("Invalid Or Expired Code"))
	}

	user, err := l.svcCtx.UserRepo.GetUserByEmail(l.ctx, email, constants.UserTypeAdmin)
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

	user, orgUser, err = l.baseLogic.UserLogic.CheckAndUpdateActiveOrganizationByUser(l.ctx, user)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	token, err := l.baseLogic.UserLogic.GenerateJwtToken(user, orgUser, logic2.JwtConfig{
		AccessSecret:  l.svcCtx.Config.JwtAuth.AccessSecret,
		AccessExpire:  l.svcCtx.Config.JwtAuth.AccessExpire,
		RefreshSecret: l.svcCtx.Config.JwtAuth.RefreshSecret,
		RefreshExpire: l.svcCtx.Config.JwtAuth.RefreshExpire,
	}, constants.UserTypeAdmin)
	if err != nil {
		return nil, err
	}
	data := &admin.UserLoginResponseData{}
	err = copier.CopyWithOption(&data, token, copier.Option{CaseSensitive: true})
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	// delete redis key
	_ = global.GetRedis().Del(l.ctx, redisKey)
	return &admin.ResetPasswordResponse{
		Data: data,
	}, nil
}
