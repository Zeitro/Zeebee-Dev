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

//go:embed send_loan_document_request.html
var SendLoanDocumentRequestEmailTpl string

const (
	// RegisterVerifyEmail = "register_verify_email"
	ResetPwdEmail = "reset_pwd_email"
	// LoginUrlEmail       = "login_url_email"
	SendLoanDocumentRequestEmail = "send_loan_document_request"
)

func init() {
	_ = email.RegisterTemplateWithText(ResetPwdEmail, SendResetPwdEmailTpl)
	_ = email.RegisterTemplateWithText(SendLoanDocumentRequestEmail, SendLoanDocumentRequestEmailTpl)
}
