package email_tpl

import (
	"embed"

	"zeebee-backend/common/email"
)

// EmailTplFS is the embedded filesystem containing email templates
var EmailTplFS embed.FS

// SendRegisterVerifyEmailTpl is the template for registration verification emails.
//
//go:embed send_register_verify_email.html
var SendRegisterVerifyEmailTpl string

// SendResetPwdEmailTpl is the template for password reset emails.
//
//go:embed send_reset_pwd_email.html
var SendResetPwdEmailTpl string

// SendLoginUrlEmailTpl is the template for login URL emails.
//
//go:embed send_login_url_email.html
var SendLoginUrlEmailTpl string

//go:embed send_self_preapproval_notification.html
var SendSelfPreApprovalNotificationTpl string

//go:embed send_third_party_preapproval_notification.html
var SendThirdPartyPreApprovalNotificationTpl string

const (
	// RegisterVerifyEmail is the key for registration verification email template.
	RegisterVerifyEmail               = "register_verify_email"
	// ResetPwdEmail is the key for password reset email template.
	ResetPwdEmail = "reset_pwd_email"
	// LoginURLEmail is the key for login URL email template.
	LoginUrlEmail = "login_url_email"
	SelfPreApprovalNotification       = "self_preapproval_notification"
	ThirdPartyPreApprovalNotification = "third_party_preapproval_notification"
)

func init() {
	_ = email.RegisterTemplateWithText(RegisterVerifyEmail, SendRegisterVerifyEmailTpl)
	_ = email.RegisterTemplateWithText(ResetPwdEmail, SendResetPwdEmailTpl)
	_ = email.RegisterTemplateWithText(LoginUrlEmail, SendLoginUrlEmailTpl)
	_ = email.RegisterTemplateWithText(SelfPreApprovalNotification, SendSelfPreApprovalNotificationTpl)
	_ = email.RegisterTemplateWithText(ThirdPartyPreApprovalNotification, SendThirdPartyPreApprovalNotificationTpl)
}
