package userlogic

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/stringx"
	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"
)

const newPassword = "NewPassw0rd!2026"

// T022 US-5: a valid reset code updates the password and issues a TokenSet;
// the old password no longer works and the code is consumed.
func TestResetPasswordByCode_Success(t *testing.T) {
	env := newTestEnv(t)
	userType := constants.UserTypeBorrower
	email := "reset@test.local"

	userId, _ := seedLoginUser(t, env, userType, email, testPassword)
	codeKey := stringx.Format(constants.ResetPwdCode, userType, email)
	require.NoError(t, env.SvcCtx.RedisCli.SetExpire(context.Background(), codeKey, verifyCode, 3600))

	l := NewResetPasswordByCodeLogic(context.Background(), env.SvcCtx)
	resp, err := l.ResetPasswordByCode(&tenant.ResetPasswordByCodeRequest{Data: &tenant.ResetPasswordByCodeRequestData{
		UserType:    string(userType),
		Email:       email,
		Code:        verifyCode,
		NewPassword: newPassword,
	}})
	require.NoError(t, err)
	require.NotNil(t, resp.GetData())
	assert.NotEmpty(t, resp.GetData().GetAccessToken())

	var user model.User
	require.NoError(t, env.DB.First(&user, userId).Error)
	assert.False(t, utils.CheckPassword(testPassword, user.PasswordHash))
	assert.True(t, utils.CheckPassword(newPassword, user.PasswordHash))

	code, err := env.SvcCtx.RedisCli.Get(context.Background(), codeKey)
	require.NoError(t, err)
	assert.Empty(t, code)
}

// T022 US-5: a wrong reset code fails and keeps the old password.
func TestResetPasswordByCode_InvalidCode(t *testing.T) {
	env := newTestEnv(t)
	userType := constants.UserTypeLender
	email := "reset-invalid@test.local"

	userId, _ := seedLoginUser(t, env, userType, email, testPassword)
	codeKey := stringx.Format(constants.ResetPwdCode, userType, email)
	require.NoError(t, env.SvcCtx.RedisCli.SetExpire(context.Background(), codeKey, verifyCode, 3600))

	l := NewResetPasswordByCodeLogic(context.Background(), env.SvcCtx)
	_, err := l.ResetPasswordByCode(&tenant.ResetPasswordByCodeRequest{Data: &tenant.ResetPasswordByCodeRequestData{
		UserType:    string(userType),
		Email:       email,
		Code:        "0000000000",
		NewPassword: newPassword,
	}})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid or expired code")

	var user model.User
	require.NoError(t, env.DB.First(&user, userId).Error)
	assert.True(t, utils.CheckPassword(testPassword, user.PasswordHash))
}

// T022 US-5: reset request email is rate limited per address.
func TestSendResetPasswordEmail_RateLimited(t *testing.T) {
	env := newTestEnv(t)
	userType := constants.UserTypeBorrower
	email := "reset-rate@test.local"

	seedLoginUser(t, env, userType, email, testPassword)
	l := NewSendResetPasswordEmailLogic(context.Background(), env.SvcCtx)
	var lastErr error
	for i := 0; i < 6; i++ { // test rate: 5 per 300s
		_, lastErr = l.SendResetPasswordEmail(&tenant.SendResetPasswordEmailRequest{Data: &tenant.SendResetPasswordEmailRequestData{
			UserType: string(userType),
			Email:    email,
		}})
	}
	require.Error(t, lastErr)
	assert.Contains(t, lastErr.Error(), "try too many times")
}

// T022 US-5: reset request for an unknown user fails.
func TestSendResetPasswordEmail_UnknownUser(t *testing.T) {
	env := newTestEnv(t)
	l := NewSendResetPasswordEmailLogic(context.Background(), env.SvcCtx)
	_, err := l.SendResetPasswordEmail(&tenant.SendResetPasswordEmailRequest{Data: &tenant.SendResetPasswordEmailRequestData{
		UserType: string(constants.UserTypeLender),
		Email:    "nobody@test.local",
	}})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")
}

// T022 US-5: reset request caches a fresh one-time code for an existing user.
func TestSendResetPasswordEmail_CachesCode(t *testing.T) {
	env := newTestEnv(t)
	userType := constants.UserTypeBorrower
	email := "reset-code@test.local"

	seedLoginUser(t, env, userType, email, testPassword)
	l := NewSendResetPasswordEmailLogic(context.Background(), env.SvcCtx)
	_, err := l.SendResetPasswordEmail(&tenant.SendResetPasswordEmailRequest{Data: &tenant.SendResetPasswordEmailRequestData{
		UserType: string(userType),
		Email:    email,
	}})
	require.NoError(t, err)

	codeKey := stringx.Format(constants.ResetPwdCode, userType, email)
	assertEventual(t, func() bool {
		code, err := env.SvcCtx.RedisCli.Get(context.Background(), codeKey)
		return err == nil && len(code) == 10
	})
}

// T022 US-5: change password rejects a wrong old password.
func TestChangePassword_WrongOldPassword(t *testing.T) {
	env := newTestEnv(t)
	userType := constants.UserTypeLender
	email := "change-wrong@test.local"

	userId, _ := seedLoginUser(t, env, userType, email, testPassword)
	l := NewChangePasswordLogic(context.Background(), env.SvcCtx)
	_, err := l.ChangePassword(&tenant.ChangePasswordRequest{Data: &tenant.ChangePasswordRequestData{
		UserType:    string(userType),
		UserId:      userId,
		OldPassword: "WrongPass!1",
		NewPassword: newPassword,
	}})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "old password is incorrect")
}

// T022 US-5: change password with the correct old password updates the hash.
func TestChangePassword_Success(t *testing.T) {
	env := newTestEnv(t)
	userType := constants.UserTypeBorrower
	email := "change-ok@test.local"

	userId, _ := seedLoginUser(t, env, userType, email, testPassword)
	l := NewChangePasswordLogic(context.Background(), env.SvcCtx)
	_, err := l.ChangePassword(&tenant.ChangePasswordRequest{Data: &tenant.ChangePasswordRequestData{
		UserType:    string(userType),
		UserId:      userId,
		OldPassword: testPassword,
		NewPassword: newPassword,
	}})
	require.NoError(t, err)

	var user model.User
	require.NoError(t, env.DB.First(&user, userId).Error)
	assert.False(t, utils.CheckPassword(testPassword, user.PasswordHash))
	assert.True(t, utils.CheckPassword(newPassword, user.PasswordHash))
}
