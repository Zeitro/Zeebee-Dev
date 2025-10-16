package email_tpl

import (
	"bytes"
	"embed"
	"html/template"

	"github.com/pkg/errors"
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

const (
	// RegisterVerifyEmail is the key for registration verification email template.
	RegisterVerifyEmail = "register_verify_email"
	// ResetPwdEmail is the key for password reset email template.
	ResetPwdEmail = "reset_pwd_email"
	// LoginURLEmail is the key for login URL email template.
	LoginUrlEmail = "login_url_email"
)

// tpls is a map of template names to parsed templates.
var tpls = map[string]*template.Template{}

// ErrEmailTemplateNotFound is returned when an email template is not found.
var ErrEmailTemplateNotFound = errors.New("email template not found")

func init() {
	tpls = map[string]*template.Template{
		RegisterVerifyEmail: template.Must(template.New(RegisterVerifyEmail).Parse(SendRegisterVerifyEmailTpl)),
		ResetPwdEmail:       template.Must(template.New(ResetPwdEmail).Parse(SendResetPwdEmailTpl)),
		LoginUrlEmail:       template.Must(template.New(LoginUrlEmail).Parse(SendLoginUrlEmailTpl)),
	}
}

// GetParsedContent returns the parsed content of an email template.
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
