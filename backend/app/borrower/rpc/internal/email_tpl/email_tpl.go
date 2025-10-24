package email_tpl

import (
	"embed"

	"zeebee-backend/common/email"
)

var EmailTplFS embed.FS

//go:embed send_register_verify_email.html
var SendRegisterVerifyEmailTpl string

//go:embed send_reset_pwd_email.html
var SendResetPwdEmailTpl string

//go:embed send_login_url_email.html
var SendLoginUrlEmailTpl string

//go:embed send_self_preapproval_notification.html
var SendSelfPreApprovalNotificationTpl string

//go:embed send_third_party_preapproval_notification.html
var SendThirdPartyPreApprovalNotificationTpl string

const (
	RegisterVerifyEmail               = "register_verify_email"
	ResetPwdEmail                     = "reset_pwd_email"
	LoginUrlEmail                     = "login_url_email"
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
