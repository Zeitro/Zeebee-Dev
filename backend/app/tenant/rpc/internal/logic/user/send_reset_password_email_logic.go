package userlogic

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/duke-git/lancet/v2/random"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"
	"gorm.io/gorm"

	"zeebee-backend/app/tenant/rpc/internal/email_tpl"
	"zeebee-backend/app/tenant/rpc/internal/logic"
	"zeebee-backend/app/tenant/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/email"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/stringx"
	"zeebee-backend/common/timex"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/limiter"
	"zeebee-backend/domains/zeebee/proto"
)

type SendResetPasswordEmailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewSendResetPasswordEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendResetPasswordEmailLogic {
	return &SendResetPasswordEmailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

// SendResetPasswordEmail rate-limits, confirms the user exists, then mails the
// reset-password link using the platform's SMTP client and brand/URL config.
func (l *SendResetPasswordEmailLogic) SendResetPasswordEmail(in *tenant.SendResetPasswordEmailRequest) (resp *proto.SendResetPasswordEmailResponse, err error) {
	userType, err := logic.ParseUserType(in.Data.UserType)
	if err != nil {
		return nil, err
	}
	emailAddr := in.Data.Email
	if emailAddr == "" {
		return nil, errorx.New("email is required")
	}

	key := stringx.Format(constants.ResetPwdEmail, userType, emailAddr)
	ok, err := l.svcCtx.LimiterMgr.Limit(l.ctx, limiter.StrategySendResetPwdEmail, key)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	if !ok {
		return nil, errorx.Errorf("email %s try too many times", emailAddr)
	}

	user, err := l.svcCtx.UserRepo.GetUserByEmail(l.ctx, emailAddr, userType)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorx.New("user not found")
		}
		return nil, errorx.WithStack(err)
	}
	if user == nil {
		return nil, errorx.New("user not found")
	}

	threading.GoSafe(func() {
		ctx := context.WithoutCancel(l.ctx)
		if err := l.sendResetPwdEmail(ctx, userType, emailAddr, user.FirstName); err != nil {
			l.svcCtx.Logger.WithContext(ctx).
				WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).
				Error("fail to send reset password email")
		}
	})

	return &proto.SendResetPasswordEmailResponse{}, nil
}

func (l *SendResetPasswordEmailLogic) sendResetPwdEmail(ctx context.Context, userType constants.UserType, emailAddr, firstName string) error {
	profile, err := l.baseLogic.PlatformFor(userType)
	if err != nil {
		return err
	}

	// Cache the one-time code before resolving the SMTP client so a client
	// outage does not block reset (and tests observe the code without SMTP).
	codeKey := stringx.Format(constants.ResetPwdCode, userType, emailAddr)
	code, _ := l.svcCtx.RedisCli.Get(ctx, codeKey)
	if len(code) == 0 {
		code = random.RandNumeral(10)
		if err := l.svcCtx.RedisCli.SetExpire(ctx, codeKey, code, profile.ResetPwdExpireSeconds); err != nil {
			return errorx.WithStack(err)
		}
	}

	client, err := l.baseLogic.EmailClientFor(userType)
	if err != nil {
		return err
	}

	content, err := email.ParseTemplate(email_tpl.ResetPwdEmail, map[string]interface{}{
		"url":            fmt.Sprintf("%s?code=%s&email=%s", profile.ResetPwdUrl, code, emailAddr),
		"first_name":     firstName,
		"expire_time":    timex.DurationToString(time.Second*time.Duration(profile.ResetPwdExpireSeconds), true),
		"logo_url":       profile.LogoUrl,
		"footer_address": profile.FooterAddress,
	})
	if err != nil {
		return errorx.WithStack(err)
	}

	return client.SendEmail(email.Message{
		Subject: "Reset your password",
		From:    senderOf(profile),
		To:      []string{emailAddr},
		Body:    content,
		IsHTML:  true,
	})
}
