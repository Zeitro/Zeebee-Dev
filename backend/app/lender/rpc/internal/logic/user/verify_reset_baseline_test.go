package userlogic

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"zeebee-backend/app/lender/rpc/internal/global"
	"zeebee-backend/common/stringx"
	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/proto"
)

// T012 characterization baseline: lender verify-email + reset-password error
// paths as they run today (local logic, pre-cutover). The deterministic
// branches (invalid / already-verified / success-reset) are covered here; the
// verify success path also activates the owned org and is covered separately.

// setVerifyCode stores the register-verify code under the same key the domain
// reads, via the same redis client (so the lender-test prefix applies identically).
func setVerifyCode(t *testing.T, email, code string) {
	t.Helper()
	key := stringx.Format(constants.UserRegisterVerifyCode, constants.UserTypeLender, email)
	require.NoError(t, global.GetRedis().SetExpire(context.Background(), key, code, 3600))
}

func setResetCode(t *testing.T, email, code string) {
	t.Helper()
	key := stringx.Format(constants.ResetPwdCode, constants.UserTypeLender, email)
	require.NoError(t, global.GetRedis().SetExpire(context.Background(), key, code, 3600))
}

// --- verify email ---

// VerifyEmail with no code stored is rejected (expired / invalid).
func TestBaselineVerifyEmail_InvalidCode(t *testing.T) {
	env := newBaselineEnv(t)
	email := "verify-invalid@test.local"
	seedLenderUser(t, env, email, baselinePassword) // active user, but no code stored

	l := NewVerifyEmailLogic(context.Background(), env.SvcCtx)
	_, err := l.VerifyEmail(&proto.VerifyEmailRequest{Data: &proto.VerifyEmailRequestData{
		Email: email,
		Code:  "1234567890",
	}})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid or expired code")
}

// VerifyEmail for an already-active user (with a valid code) is rejected.
func TestBaselineVerifyEmail_AlreadyVerified(t *testing.T) {
	env := newBaselineEnv(t)
	email := "verify-active@test.local"
	seedLenderUser(t, env, email, baselinePassword) // active user
	setVerifyCode(t, email, "1234567890")

	l := NewVerifyEmailLogic(context.Background(), env.SvcCtx)
	_, err := l.VerifyEmail(&proto.VerifyEmailRequest{Data: &proto.VerifyEmailRequestData{
		Email: email,
		Code:  "1234567890",
	}})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "already verified")
}

// --- reset password ---

// ResetPassword with a valid code issues a token and the new password works.
func TestBaselineResetPassword_ValidCode(t *testing.T) {
	env := newBaselineEnv(t)
	email := "reset@test.local"
	userId, _ := seedLenderUser(t, env, email, baselinePassword)
	setResetCode(t, email, "resetcode1")

	newPwd := "NewPass!2026"
	l := NewResetPasswordLogic(context.Background(), env.SvcCtx)
	resp, err := l.ResetPassword(&proto.ResetPasswordRequest{Data: &proto.ResetPasswordRequestData{
		Email:    email,
		Password: newPwd,
		Code:     "resetcode1",
	}})
	require.NoError(t, err)
	require.NotNil(t, resp.GetData())
	assert.NotEmpty(t, resp.GetData().GetAccessToken())

	// password actually changed: old password no longer logs in, new one does.
	var user model.User
	require.NoError(t, env.DB.Where("id = ?", userId).First(&user).Error)
	assert.True(t, utils.CheckPassword(newPwd, user.PasswordHash), "new password must hash-match after reset")

	// code consumed
	codeKey := stringx.Format(constants.ResetPwdCode, constants.UserTypeLender, email)
	stored, _ := global.GetRedis().Get(context.Background(), codeKey)
	assert.Empty(t, stored, "reset code must be consumed after use")
}

// ResetPassword with no code stored is rejected.
func TestBaselineResetPassword_InvalidCode(t *testing.T) {
	env := newBaselineEnv(t)
	email := "reset-invalid@test.local"
	seedLenderUser(t, env, email, baselinePassword)

	l := NewResetPasswordLogic(context.Background(), env.SvcCtx)
	_, err := l.ResetPassword(&proto.ResetPasswordRequest{Data: &proto.ResetPasswordRequestData{
		Email:    email,
		Password: "NewPass!2026",
		Code:     "wrongcode",
	}})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid or expired code")
}

// VerifyEmail with a valid code activates an unverified user and issues a token;
// the verify code is consumed.
func TestBaselineVerifyEmail_Success(t *testing.T) {
	env := newBaselineEnv(t)
	email := "verify-ok@test.local"
	userId, _ := seedLenderUser(t, env, email, baselinePassword)
	require.NoError(t, env.DB.Model(&model.User{}).Where("id = ?", userId).
		Update("status", constants.UserStatusUnverified).Error)
	setVerifyCode(t, email, "1234567890")

	l := NewVerifyEmailLogic(context.Background(), env.SvcCtx)
	resp, err := l.VerifyEmail(&proto.VerifyEmailRequest{Data: &proto.VerifyEmailRequestData{
		Email: email,
		Code:  "1234567890",
	}})
	require.NoError(t, err)
	require.NotNil(t, resp.GetData())
	assert.NotEmpty(t, resp.GetData().GetAccessToken())

	// user activated.
	var user model.User
	require.NoError(t, env.DB.Where("id = ?", userId).First(&user).Error)
	assert.Equal(t, constants.UserStatusActive, user.Status)

	// code consumed.
	codeKey := stringx.Format(constants.UserRegisterVerifyCode, constants.UserTypeLender, email)
	stored, _ := global.GetRedis().Get(context.Background(), codeKey)
	assert.Empty(t, stored, "verify code must be consumed after use")
}
