package userlogic

import (
	"context"
	"fmt"
	"time"

	"github.com/duke-git/lancet/v2/random"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"

	"zeebee-backend/app/tenant/config"
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

type SendVerifyEmailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewSendVerifyEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendVerifyEmailLogic {
	return &SendVerifyEmailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

// SendVerifyEmail rate-limits, confirms the user is still unverified, then mails
// the verify link using the platform's SMTP client and brand/URL config.
func (l *SendVerifyEmailLogic) SendVerifyEmail(in *tenant.SendVerifyEmailRequest) (resp *proto.SendVerifyEmailResponse, err error) {
	userType, err := logic.ParseUserType(in.Data.UserType)
	if err != nil {
		return nil, err
	}
	emailAddr := in.Data.Email

	key := stringx.Format(constants.RegisterVerifyEmail, userType, emailAddr)
	ok, err := l.svcCtx.LimiterMgr.Limit(l.ctx, limiter.StrategySendVerifyEmail, key)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	if !ok {
		return nil, errorx.Errorf("email %s try too many times", emailAddr)
	}

	user, err := l.svcCtx.UserRepo.GetUserByEmail(l.ctx, emailAddr, userType)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errorx.New("user not found")
	}
	if user.Status != constants.UserStatusUnverified {
		return nil, errorx.New("user is verified")
	}

	threading.GoSafe(func() {
		ctx := context.WithoutCancel(l.ctx)
		if err := sendVerifyEmail(ctx, l.baseLogic, l.svcCtx, userType, emailAddr, user.FirstName); err != nil {
			l.svcCtx.Logger.WithContext(ctx).
				WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).
				Error("fail to send verify email")
		}
	})

	return &proto.SendVerifyEmailResponse{}, nil
}

// sendLoginUrlEmail rate-limits and mails the platform login link to an already
// registered (active) user who re-submitted registration. Mirrors the borrower /
// lender "already activated" path with tenant-owned brand/URL config.
func sendLoginUrlEmail(ctx context.Context, baseLogic *logic.BaseLogic, svcCtx *svc.ServiceContext, userType constants.UserType, emailAddr, firstName string) error {
	key := stringx.Format(constants.LoginEmail, userType, emailAddr)
	ok, err := svcCtx.LimiterMgr.Limit(ctx, limiter.StrategySendLoginEmail, key)
	if err != nil {
		return errorx.WithStack(err)
	}
	if !ok {
		return errorx.Errorf("email %s try too many times", emailAddr)
	}

	profile, err := baseLogic.PlatformFor(userType)
	if err != nil {
		return err
	}
	client, err := baseLogic.EmailClientFor(userType)
	if err != nil {
		return err
	}

	content, err := email.ParseTemplate(email_tpl.LoginUrlEmail, map[string]interface{}{
		"url":            profile.LoginUrl,
		"first_name":     firstName,
		"logo_url":       profile.LogoUrl,
		"footer_address": profile.FooterAddress,
	})
	if err != nil {
		return errorx.WithStack(err)
	}

	return client.SendEmail(email.Message{
		Subject: "Login to your account",
		From:    senderOf(profile),
		To:      []string{emailAddr},
		Body:    content,
		IsHTML:  true,
	})
}

// logSendError records an async account-email failure with its stack for
// observability without failing the request that triggered the send.
func logSendError(svcCtx *svc.ServiceContext, ctx context.Context, msg string, err error) {
	svcCtx.Logger.WithContext(ctx).
		WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).
		Error(msg)
}

// sendVerifyEmail generates (or reuses) the one-time verify code, caches it, and
// mails the verify link using the platform's SMTP client + brand config. Shared
// by SendVerifyEmail and Register.
func sendVerifyEmail(ctx context.Context, baseLogic *logic.BaseLogic, svcCtx *svc.ServiceContext, userType constants.UserType, emailAddr, firstName string) error {
	profile, err := baseLogic.PlatformFor(userType)
	if err != nil {
		return err
	}

	// Generate and cache the one-time code before resolving the SMTP client so a
	// client outage does not block verification (and tests observe the code even
	// when no SMTP is configured).
	codeKey := stringx.Format(constants.UserRegisterVerifyCode, userType, emailAddr)
	code, _ := svcCtx.RedisCli.Get(ctx, codeKey)
	if len(code) == 0 {
		code = random.RandNumeral(10)
		if err := svcCtx.RedisCli.SetExpire(ctx, codeKey, code, profile.VerifyEmailExpireSeconds); err != nil {
			return errorx.WithStack(err)
		}
	}

	client, err := baseLogic.EmailClientFor(userType)
	if err != nil {
		return err
	}

	content, err := email.ParseTemplate(email_tpl.RegisterVerifyEmail, map[string]interface{}{
		"url":            fmt.Sprintf("%s?code=%s&email=%s", profile.VerifyEmailUrl, code, emailAddr),
		"first_name":     firstName,
		"expire_time":    timex.DurationToString(time.Second*time.Duration(profile.VerifyEmailExpireSeconds), true),
		"logo_url":       profile.LogoUrl,
		"footer_address": profile.FooterAddress,
	})
	if err != nil {
		return errorx.WithStack(err)
	}

	return client.SendEmail(email.Message{
		Subject: "Verify your email",
		From:    senderOf(profile),
		To:      []string{emailAddr},
		Body:    content,
		IsHTML:  true,
	})
}

// senderOf prefers the configured sender email, falling back to the SMTP username.
func senderOf(profile config.PlatformProfile) string {
	if profile.SenderEmail != "" {
		return profile.SenderEmail
	}
	return profile.EmailConf.Username
}
