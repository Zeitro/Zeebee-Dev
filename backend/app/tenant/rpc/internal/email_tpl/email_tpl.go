// Package email_tpl holds tenant-owned account-flow email templates
// (register verify, reset password, login url). These consolidate the
// per-platform copies that previously lived in each backend; brand assets
// (logo, footer) are injected per platform via PlatformConf.
package email_tpl

import (
	_ "embed"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/email"
)

//go:embed register_verify_email.html
var registerVerifyEmailTpl string

//go:embed reset_pwd_email.html
var resetPwdEmailTpl string

//go:embed login_url_email.html
var loginUrlEmailTpl string

// Template keys registered with common/email.
const (
	RegisterVerifyEmail = "tenant_register_verify_email"
	ResetPwdEmail       = "tenant_reset_pwd_email"
	LoginUrlEmail       = "tenant_login_url_email"
)

func init() {
	for name, tpl := range map[string]string{
		RegisterVerifyEmail: registerVerifyEmailTpl,
		ResetPwdEmail:       resetPwdEmailTpl,
		LoginUrlEmail:       loginUrlEmailTpl,
	} {
		if err := email.RegisterTemplateWithText(name, tpl); err != nil {
			logx.Must(err)
		}
	}
}
