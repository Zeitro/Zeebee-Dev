package email_tpl

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"zeebee-backend/app/tenant/config"
	"zeebee-backend/common/email"
)

// T044 US-11: tenant account templates inject the per-platform brand (logo,
// footer) and platform URLs. Rendering with a production-style borrower
// profile must not leak a test-environment URL.
func TestAccountTemplates_InjectPlatformBrand(t *testing.T) {
	profile := config.PlatformProfile{
		LoginUrl:       "https://app.zeebee.com/login",
		VerifyEmailUrl: "https://app.zeebee.com/verify-email",
		ResetPwdUrl:    "https://app.zeebee.com/reset-password",
		LogoUrl:        "https://cdn.zeebee.com/logo.png",
		FooterAddress:  "Zeebee Inc., 1 Market Street, San Francisco, CA",
	}
	code := "1234567890"
	emailAddr := "borrower@example.com"
	firstName := "Ada"

	cases := []struct {
		name string
		tpl  string
		data map[string]interface{}
		url  string
	}{
		{
			name: "register verify",
			tpl:  RegisterVerifyEmail,
			data: map[string]interface{}{
				"url":            profile.VerifyEmailUrl + "?code=" + code + "&email=" + emailAddr,
				"first_name":     firstName,
				"expire_time":    "24h",
				"logo_url":       profile.LogoUrl,
				"footer_address": profile.FooterAddress,
			},
			url: profile.VerifyEmailUrl,
		},
		{
			name: "login url",
			tpl:  LoginUrlEmail,
			data: map[string]interface{}{
				"url":            profile.LoginUrl,
				"first_name":     firstName,
				"logo_url":       profile.LogoUrl,
				"footer_address": profile.FooterAddress,
			},
			url: profile.LoginUrl,
		},
		{
			name: "reset password",
			tpl:  ResetPwdEmail,
			data: map[string]interface{}{
				"url":            profile.ResetPwdUrl + "?token=abc",
				"first_name":     firstName,
				"expire_time":    "1h",
				"logo_url":       profile.LogoUrl,
				"footer_address": profile.FooterAddress,
			},
			url: profile.ResetPwdUrl,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			body, err := email.ParseTemplate(tc.tpl, tc.data)
			require.NoError(t, err)
			assert.Contains(t, body, profile.LogoUrl)
			assert.Contains(t, body, profile.FooterAddress)
			assert.Contains(t, body, tc.url)
			// Borrower URLs must not point at a test environment.
			assert.NotContains(t, body, "test.local")
			assert.NotContains(t, body, "127.0.0.1")
		})
	}
}
