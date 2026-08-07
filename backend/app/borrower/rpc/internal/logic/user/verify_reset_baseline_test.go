package userlogic

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"zeebee-backend/app/borrower/rpc/internal/global"
	"zeebee-backend/common/stringx"
	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/proto"
)

// T012 characterization baseline: borrower verify-email (error paths only) +
// reset-password. The verify success path triggers preapproval side-effects
// (handlePreapprovalLoan → pdf_generator / zeitro clients) that the in-memory
// harness cannot satisfy, so only the deterministic branches that return before
// that side-effect are recorded here. The cutover must keep these green.

func setVerifyCode(t *testing.T, email, code string) {
	t.Helper()
	key := stringx.Format(constants.UserRegisterVerifyCode, constants.UserTypeBorrower, email)
	require.NoError(t, global.GetRedis().SetExpire(context.Background(), key, code, 3600))
}

func setBorrowerResetCode(t *testing.T, email, code string) {
	t.Helper()
	key := stringx.Format(constants.ResetPwdCode, constants.UserTypeBorrower, email)
	require.NoError(t, global.GetRedis().SetExpire(context.Background(), key, code, 3600))
}

// --- verify email (error paths only) ---

// VerifyEmail with no code stored is rejected (expired / invalid).
func TestBaselineVerifyEmail_InvalidCode(t *testing.T) {
	env := newBaselineEnv(t)
	email := "verify-invalid@test.local"
	seedBorrowerUser(t, env, email, baselinePassword) // active user, but no code stored

	l := NewVerifyEmailLogic(context.Background(), env.SvcCtx)
	_, err := l.VerifyEmail(&proto.VerifyEmailRequest{Data: &proto.VerifyEmailRequestData{
		Email: email,
		Code:  "1234567890",
	}})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid or expired code")
}

// VerifyEmail for an already-active user (with a valid code) is rejected before
// any preapproval side-effect runs.
func TestBaselineVerifyEmail_AlreadyVerified(t *testing.T) {
	env := newBaselineEnv(t)
	email := "verify-active@test.local"
	seedBorrowerUser(t, env, email, baselinePassword) // active user
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
	userId, _ := seedBorrowerUser(t, env, email, baselinePassword)
	setBorrowerResetCode(t, email, "resetcode1")

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

	// password actually changed: new password hash-matches.
	var user model.User
	require.NoError(t, env.DB.Where("id = ?", userId).First(&user).Error)
	assert.True(t, utils.CheckPassword(newPwd, user.PasswordHash), "new password must hash-match after reset")

	// code consumed
	codeKey := stringx.Format(constants.ResetPwdCode, constants.UserTypeBorrower, email)
	stored, _ := global.GetRedis().Get(context.Background(), codeKey)
	assert.Empty(t, stored, "reset code must be consumed after use")
}

// ResetPassword with no code stored is rejected.
func TestBaselineResetPassword_InvalidCode(t *testing.T) {
	env := newBaselineEnv(t)
	email := "reset-invalid@test.local"
	seedBorrowerUser(t, env, email, baselinePassword)

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
// the verify code is consumed. The preapproval side-effect goroutine is
// fire-and-forget business logic (not the identity contract under test); with no
// draft loan seeded it takes its not-found early return, so it neither affects
// the response nor crashes the test.
func TestBaselineVerifyEmail_Success(t *testing.T) {
	env := newBaselineEnv(t)
	email := "verify-ok@test.local"
	userId, _ := seedBorrowerUser(t, env, email, baselinePassword)
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
	codeKey := stringx.Format(constants.UserRegisterVerifyCode, constants.UserTypeBorrower, email)
	stored, _ := global.GetRedis().Get(context.Background(), codeKey)
	assert.Empty(t, stored, "verify code must be consumed after use")
}
