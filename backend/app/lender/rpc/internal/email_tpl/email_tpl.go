package email_tpl

import (
	"embed"

	"zeebee-backend/common/email"
)

// EmailTplFS is the embedded filesystem containing email templates.
var EmailTplFS embed.FS

// RegisterVerifyEmailTpl is the template for registration verification emails.
//
//go:embed register_verify_email.html
var RegisterVerifyEmailTpl string

// ResetPwdEmailTpl is the template for password reset emails.
//
//go:embed reset_pwd_email.html
var ResetPwdEmailTpl string

// LoginUrlEmailTpl is the template for login URL emails.
//
//go:embed login_url_email.html
var LoginUrlEmailTpl string

// LoanDocumentRequestEmailTpl is the template for loan document request emails.
//
//go:embed loan_document_request.html
var LoanDocumentRequestEmailTpl string

//go:embed user_invite_email.html
var UserInviteEmailTpl string

const (
	RegisterVerifyEmail          = "register_verify_email"
	ResetPwdEmail                = "reset_pwd_email"
	LoginUrlEmail                = "login_url_email"
	SendLoanDocumentRequestEmail = "send_loan_document_request"
	UserInviteEmail              = "user_invite_email"
)

func init() {
	_ = email.RegisterTemplateWithText(RegisterVerifyEmail, RegisterVerifyEmailTpl)
	_ = email.RegisterTemplateWithText(LoginUrlEmail, LoginUrlEmailTpl)
	_ = email.RegisterTemplateWithText(ResetPwdEmail, ResetPwdEmailTpl)
	_ = email.RegisterTemplateWithText(SendLoanDocumentRequestEmail, LoanDocumentRequestEmailTpl)
	_ = email.RegisterTemplateWithText(UserInviteEmail, UserInviteEmailTpl)
}
