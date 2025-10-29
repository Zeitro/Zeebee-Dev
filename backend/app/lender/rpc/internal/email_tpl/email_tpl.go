package email_tpl

import (
	"bytes"
	"embed"
	"html/template"

	"github.com/pkg/errors"
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

var tpls = map[string]*template.Template{}
var ErrEmailTemplateNotFound = errors.New("email template not found")

func init() {
	tpls = map[string]*template.Template{
		// RegisterVerifyEmail: template.Must(template.New(RegisterVerifyEmail).Parse(SendRegisterVerifyEmailTpl)),
		ResetPwdEmail: template.Must(template.New(ResetPwdEmail).Parse(SendResetPwdEmailTpl)),
		// LoginUrlEmail:       template.Must(template.New(LoginUrlEmail).Parse(SendLoginUrlEmailTpl)),
		SendLoanDocumentRequestEmail: template.Must(template.New(SendLoanDocumentRequestEmail).Parse(SendLoanDocumentRequestEmailTpl)),
	}
}

func GetParsedContent(name string, data map[string]interface{}) (content string, err error) {
	tpl, ok := tpls[name]
	if !ok {
		return "", ErrEmailTemplateNotFound
	}
	buf := new(bytes.Buffer)
	err = tpl.Execute(buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
