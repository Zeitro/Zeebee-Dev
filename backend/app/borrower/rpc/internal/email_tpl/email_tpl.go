package email_tpl

import (
	"embed"

	"golang.org/x/sync/errgroup"

	"zeebee-backend/common/email"
)

// EmailTplFS is the embedded filesystem containing email templates
var EmailTplFS embed.FS

// SendRegisterVerifyEmailTpl is the template for registration verification emails.
//
//go:embed register_verify_email.html
var RegisterVerifyEmailTpl string

// SendResetPwdEmailTpl is the template for password reset emails.
//
//go:embed reset_pwd_email.html
var ResetPwdEmailTpl string

// SendLoginUrlEmailTpl is the template for login URL emails.
//
//go:embed login_url_email.html
var LoginUrlEmailTpl string

//go:embed user_invite_email.html
var UserInviteEmailTpl string

//go:embed self_preapproval_notification.html
var SelfPreApprovalNotificationTpl string

//go:embed third_party_preapproval_notification.html
var ThirdPartyPreApprovalNotificationTpl string

//go:embed quote_email.html
var QuoteEmailTpl string

const (
	// RegisterVerifyEmail is the key for registration verification email template.
	RegisterVerifyEmail = "register_verify_email"
	// ResetPwdEmail is the key for password reset email template.
	ResetPwdEmail = "reset_pwd_email"
	// LoginURLEmail is the key for login URL email template.
	LoginUrlEmail                     = "login_url_email"
	UserInviteEmail                   = "user_invite_email"
	SelfPreApprovalNotification       = "self_preapproval_notification"
	ThirdPartyPreApprovalNotification = "third_party_preapproval_notification"
	QuoteEmail                        = "quote_email"
)

func init() {

	eg := errgroup.Group{}
	eg.Go(func() error {
		return email.RegisterTemplateWithText(RegisterVerifyEmail, RegisterVerifyEmailTpl)
	})
	eg.Go(func() error {
		return email.RegisterTemplateWithText(ResetPwdEmail, ResetPwdEmailTpl)
	})
	eg.Go(func() error {
		return email.RegisterTemplateWithText(LoginUrlEmail, LoginUrlEmailTpl)
	})
	eg.Go(func() error {
		return email.RegisterTemplateWithText(UserInviteEmail, UserInviteEmailTpl)
	})
	eg.Go(func() error {
		return email.RegisterTemplateWithText(SelfPreApprovalNotification, SelfPreApprovalNotificationTpl)
	})
	eg.Go(func() error {
		return email.RegisterTemplateWithText(ThirdPartyPreApprovalNotification, ThirdPartyPreApprovalNotificationTpl)
	})
	eg.Go(func() error {
		return email.RegisterTemplateWithText(QuoteEmail, QuoteEmailTpl)
	})
	err := eg.Wait()
	if err != nil {
		panic(err)
	}
}
