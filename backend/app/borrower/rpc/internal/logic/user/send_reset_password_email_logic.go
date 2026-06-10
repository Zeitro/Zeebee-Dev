package userlogic

import (
	"context"
	"fmt"

	"github.com/duke-git/lancet/v2/random"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/rpc/internal/email_tpl"
	"zeebee-backend/app/borrower/rpc/internal/global"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	email2 "zeebee-backend/common/email"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/stringx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/limiter"
)

type SendResetPasswordEmailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendResetPasswordEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendResetPasswordEmailLogic {
	return &SendResetPasswordEmailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SendResetPasswordEmailLogic) SendResetPasswordEmail(in *borrower.SendResetPasswordEmailRequest) (resp *borrower.SendResetPasswordEmailResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
				logx.Field("service_name", l.svcCtx.Config.Name),
			).Error(err)
		}
	}()

	email := in.Data.Email
	if email == "" {
		return nil, errorx.WithStack(errorx.New("Email is required"))
	}

	// rate limit
	key := stringx.Format(constants.ResetPwdEmail, constants.UserTypeBorrower, email)
	ok, err := global.LimiterMgr.Limit(l.ctx, limiter.StrategySendResetPwdEmail, key)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	if !ok {
		return nil, errorx.Errorf("Email %s try too many times", email)
	}

	// get user
	user, err := l.svcCtx.UserRepo.GetUserByEmail(l.ctx, email, constants.UserTypeBorrower)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	if user == nil {
		return nil, errorx.WithStack(errorx.New("User not found"))
	}

	go func() {
		ctx := context.WithoutCancel(l.ctx)
		// send email
		err := l.sendResetPwdEmail(email, user.FirstName)
		if err != nil {
			err = errorx.Wrap(err, "Send Reset Password Email Error")
			l.WithContext(ctx).Error(err)
		}
	}()

	return &borrower.SendResetPasswordEmailResponse{}, err
}

func (l *SendResetPasswordEmailLogic) sendResetPwdEmail(email string, firstName string) (err error) {
	ctx := context.WithoutCancel(l.ctx)
	var code string
	redisKey := stringx.Format(constants.ResetPwdCode, constants.UserTypeBorrower, email)
	code, _ = global.GetRedis().Get(ctx, redisKey)
	if len(code) == 0 {
		code = random.RandNumeral(10)
		expireSeconds := global.Config.ResetPwdExpireSeconds
		err = global.GetRedis().SetExpire(ctx, redisKey, code, expireSeconds)
		if err != nil {
			return errorx.WithStack(errorx.Wrap(err, "Failed To Cache Reset Password Code"))
		}
	}

	// send reset pwd email
	data := map[string]interface{}{
		"url":        fmt.Sprintf("%s?code=%s&email=%s", global.Config.ResetPwdUrl, code, email),
		"first_name": firstName,
	}
	content, err := email2.ParseTemplate(email_tpl.ResetPwdEmail, data)
	if err != nil {
		return errorx.WithStack(err)
	}
	msg := email2.Message{
		Subject: "Reset your password",
		From:    global.Config.EmailConf.Username,
		To:      []string{email},
		Body:    content,
		IsHTML:  true,
	}
	err = global.EmailClient.SendEmail(msg)

	return errorx.WithStack(err)
}
