package email_tpl

import (
	"embed"

	"zeebee-backend/common/email"
)

// EmailTplFS is the embedded filesystem containing email templates.
var EmailTplFS embed.FS

// SendRegisterVerifyEmailTpl is the template for registration verification emails.
//
//go:embed send_register_verify_email.html
var SendRegisterVerifyEmailTpl string

// SendResetPwdEmailTpl is the template for password reset emails.
// SendResetPwdEmailTpl is the template for password reset emails.
//
//go:embed send_reset_pwd_email.html
var SendResetPwdEmailTpl string

// SendLoginUrlEmailTpl is the template for login URL emails.
//
//go:embed send_login_url_email.html
var SendLoginUrlEmailTpl string

// SendLoanDocumentRequestEmailTpl is the template for loan document request emails.
//
//go:embed send_loan_document_request.html
var SendLoanDocumentRequestEmailTpl string

//go:embed user_invite_email.html
var UserInviteEmailTpl string

const (
	// RegisterVerifyEmail = "register_verify_email"
	// ResetPwdEmail is the key for password reset email template.
	ResetPwdEmail = "reset_pwd_email"
	// LoginUrlEmail       = "login_url_email"
	// SendLoanDocumentRequestEmail is the key for loan document request email
	SendLoanDocumentRequestEmail = "send_loan_document_request"
	// UserInviteEmail is the key for user invite email template
	UserInviteEmail = "user_invite_email"
)

func init() {
	_ = email.RegisterTemplateWithText(ResetPwdEmail, SendResetPwdEmailTpl)
	_ = email.RegisterTemplateWithText(SendLoanDocumentRequestEmail, SendLoanDocumentRequestEmailTpl)
	_ = email.RegisterTemplateWithText(UserInviteEmail, UserInviteEmailTpl)
}
