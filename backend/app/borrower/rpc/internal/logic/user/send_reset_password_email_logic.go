package userlogic

import (
	"context"
	"fmt"
	"time"

	"github.com/duke-git/lancet/v2/random"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/rpc/internal/email_tpl"
	"zeebee-backend/app/borrower/rpc/internal/global"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/constants"
	email2 "zeebee-backend/common/email"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/limiter"
	"zeebee-backend/common/stringx"
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
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	email := in.Data.Email
	if email == "" {
		return nil, errors.WithStack(errors.New("Email is required"))
	}

	// rate limit
	key := stringx.Format(constants.BorrowerResetPwdEmail, email)
	limitOpt := limiter.WithRedisLimit(
		global.Config.LimiterConf.SendResetPwdEmail.Rate,
		time.Second*time.Duration(global.Config.LimiterConf.SendResetPwdEmail.Period),
	)
	ok, err := global.SendResetPasswordEmailLimiter.Limit(l.ctx, key, limitOpt)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if !ok {
		return nil, errors.Errorf("Email %s try too many times", email)
	}

	// get user
	user, err := l.svcCtx.UserRepo.GetBorrowerUserByEmail(l.ctx, email)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if user == nil || user.Borrower == nil {
		return nil, errors.WithStack(errors.New("User not found"))
	}

	go func() {
		ctx := context.WithoutCancel(l.ctx)
		// send email
		err := l.sendResetPwdEmail(email, user.Borrower.FirstName)
		if err != nil {
			err = errors.Wrap(err, "send reset password email error")
			l.WithContext(ctx).Error(err)
		}
	}()

	return &borrower.SendResetPasswordEmailResponse{}, err
}

func (l *SendResetPasswordEmailLogic) sendResetPwdEmail(email string, firstName string) (err error) {
	ctx := context.WithoutCancel(l.ctx)
	var code string
	redisKey := stringx.Format(constants.BorrowerResetPwdCode, email)
	code, _ = global.GetRedis().Get(ctx, redisKey)
	if len(code) == 0 {
		code = random.RandNumeral(10)
		expireSeconds := global.Config.ResetPwdExpireSeconds
		err = global.GetRedis().SetExpire(ctx, redisKey, code, expireSeconds)
		if err != nil {
			return errors.WithStack(errors.Wrap(err, "Failed to cache reset password code"))
		}
	}

	// send reset pwd email
	data := map[string]interface{}{
		"url":        fmt.Sprintf("%s?code=%s&email=%s", global.Config.ResetPwdUrl, code, email),
		"first_name": firstName,
	}
	content, err := email_tpl.GetParsedContent(email_tpl.ResetPwdEmail, data)
	if err != nil {
		return errors.WithStack(err)
	}
	msg := email2.EmailMessage{
		Subject: "Reset your password",
		From:    global.Config.EmailConf.Username,
		To:      []string{email},
		Body:    content,
		IsHTML:  true,
	}
	err = global.EmailClient.SendEmail(msg)

	return errors.WithStack(err)
}
