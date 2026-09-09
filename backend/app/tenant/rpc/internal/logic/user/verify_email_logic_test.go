package userlogic

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/stringx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"
)

const verifyCode = "1234567890"

// T017 US-3: a correct one-time code activates the user, its owned organization
// and the org-user binding, then issues a TokenSet and consumes the code.
func TestVerifyEmail_Success(t *testing.T) {
	env := newTestEnv(t)
	userType := constants.UserTypeBorrower
	email := "verify@test.local"

	userId, orgId := seedLoginUser(t, env, userType, email, testPassword)
	require.NoError(t, env.DB.Model(&model.User{}).
		Where("id = ?", userId).
		Updates(map[string]interface{}{
			"status": constants.UserStatusUnverified,
		}).Error)
	require.NoError(t, env.DB.Model(&model.Organization{}).
		Where("id = ?", orgId).
		Updates(map[string]interface{}{
			"status":   constants.OrganizationStatusNotActivated,
			"owner_id": userId,
		}).Error)
	require.NoError(t, env.DB.Model(&model.OrganizationUser{}).
		Where("user_id = ? AND organization_id = ?", userId, orgId).
		Update("status", constants.OrganizationUserStatusNotActivated).Error)

	codeKey := stringx.Format(constants.UserRegisterVerifyCode, userType, email)
	require.NoError(t, env.SvcCtx.RedisCli.SetExpire(context.Background(), codeKey, verifyCode, 3600))

	l := NewVerifyEmailLogic(context.Background(), env.SvcCtx)
	resp, err := l.VerifyEmail(&tenant.VerifyEmailRequest{Data: &tenant.VerifyEmailRequestData{
		UserType: string(userType),
		Email:    email,
		Code:     verifyCode,
	}})
	require.NoError(t, err)
	require.NotNil(t, resp.GetData())
	assert.NotEmpty(t, resp.GetData().GetAccessToken())

	var user model.User
	require.NoError(t, env.DB.First(&user, userId).Error)
	assert.Equal(t, constants.UserStatusActive, user.Status)

	var org model.Organization
	require.NoError(t, env.DB.First(&org, orgId).Error)
	assert.Equal(t, constants.OrganizationStatusActive, org.Status)

	var ou model.OrganizationUser
	require.NoError(t, env.DB.Where("user_id = ? AND organization_id = ?", userId, orgId).First(&ou).Error)
	assert.Equal(t, constants.OrganizationUserStatusActive, ou.Status)

	// Code is consumed after successful verification.
	code, err := env.SvcCtx.RedisCli.Get(context.Background(), codeKey)
	require.NoError(t, err)
	assert.Empty(t, code)
}

// T017 US-3: a wrong code fails and leaves the user unverified.
func TestVerifyEmail_WrongCode(t *testing.T) {
	env := newTestEnv(t)
	userType := constants.UserTypeLender
	email := "wrong@test.local"

	userId, _ := seedLoginUser(t, env, userType, email, testPassword)
	require.NoError(t, env.DB.Model(&model.User{}).
		Where("id = ?", userId).
		Update("status", constants.UserStatusUnverified).Error)

	codeKey := stringx.Format(constants.UserRegisterVerifyCode, userType, email)
	require.NoError(t, env.SvcCtx.RedisCli.SetExpire(context.Background(), codeKey, verifyCode, 3600))

	l := NewVerifyEmailLogic(context.Background(), env.SvcCtx)
	_, err := l.VerifyEmail(&tenant.VerifyEmailRequest{Data: &tenant.VerifyEmailRequestData{
		UserType: string(userType),
		Email:    email,
		Code:     "0000000000",
	}})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid or expired code")

	var user model.User
	require.NoError(t, env.DB.First(&user, userId).Error)
	assert.Equal(t, constants.UserStatusUnverified, user.Status)
}

// T017 US-3: an expired (absent) code fails with the same error.
func TestVerifyEmail_ExpiredCode(t *testing.T) {
	env := newTestEnv(t)
	userType := constants.UserTypeBorrower
	email := "expired@test.local"

	userId, _ := seedLoginUser(t, env, userType, email, testPassword)
	require.NoError(t, env.DB.Model(&model.User{}).
		Where("id = ?", userId).
		Update("status", constants.UserStatusUnverified).Error)

	l := NewVerifyEmailLogic(context.Background(), env.SvcCtx)
	_, err := l.VerifyEmail(&tenant.VerifyEmailRequest{Data: &tenant.VerifyEmailRequestData{
		UserType: string(userType),
		Email:    email,
		Code:     verifyCode,
	}})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid or expired code")
}

// T017 US-3: verifying an already-active user is rejected.
func TestVerifyEmail_AlreadyVerified(t *testing.T) {
	env := newTestEnv(t)
	userType := constants.UserTypeLender
	email := "active-verify@test.local"

	seedLoginUser(t, env, userType, email, testPassword) // seeds Active user

	codeKey := stringx.Format(constants.UserRegisterVerifyCode, userType, email)
	require.NoError(t, env.SvcCtx.RedisCli.SetExpire(context.Background(), codeKey, verifyCode, 3600))

	l := NewVerifyEmailLogic(context.Background(), env.SvcCtx)
	_, err := l.VerifyEmail(&tenant.VerifyEmailRequest{Data: &tenant.VerifyEmailRequestData{
		UserType: string(userType),
		Email:    email,
		Code:     verifyCode,
	}})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "email is already verified")
}

// T017 US-3: SendVerifyEmail caches a fresh one-time code for an unverified user.
func TestSendVerifyEmail_Success(t *testing.T) {
	env := newTestEnv(t)
	userType := constants.UserTypeBorrower
	email := "send-verify@test.local"

	userId, _ := seedLoginUser(t, env, userType, email, testPassword)
	require.NoError(t, env.DB.Model(&model.User{}).
		Where("id = ?", userId).
		Update("status", constants.UserStatusUnverified).Error)

	l := NewSendVerifyEmailLogic(context.Background(), env.SvcCtx)
	_, err := l.SendVerifyEmail(&tenant.SendVerifyEmailRequest{Data: &tenant.SendVerifyEmailRequestData{
		UserType: string(userType),
		Email:    email,
	}})
	require.NoError(t, err)

	codeKey := stringx.Format(constants.UserRegisterVerifyCode, userType, email)
	assertEventual(t, func() bool {
		code, err := env.SvcCtx.RedisCli.Get(context.Background(), codeKey)
		return err == nil && len(code) == 10
	})
}

// T017 US-3: SendVerifyEmail rejects verified users.
func TestSendVerifyEmail_AlreadyVerified(t *testing.T) {
	env := newTestEnv(t)
	userType := constants.UserTypeLender
	email := "send-verify-active@test.local"

	seedLoginUser(t, env, userType, email, testPassword) // Active

	l := NewSendVerifyEmailLogic(context.Background(), env.SvcCtx)
	_, err := l.SendVerifyEmail(&tenant.SendVerifyEmailRequest{Data: &tenant.SendVerifyEmailRequestData{
		UserType: string(userType),
		Email:    email,
	}})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "user is verified")
}

// T017 US-3: SendVerifyEmail is rate limited per email.
func TestSendVerifyEmail_RateLimited(t *testing.T) {
	env := newTestEnv(t)
	userType := constants.UserTypeBorrower
	email := "rate-limited@test.local"

	userId, _ := seedLoginUser(t, env, userType, email, testPassword)
	require.NoError(t, env.DB.Model(&model.User{}).
		Where("id = ?", userId).
		Update("status", constants.UserStatusUnverified).Error)

	l := NewSendVerifyEmailLogic(context.Background(), env.SvcCtx)
	var lastErr error
	for i := 0; i < 6; i++ { // test rate: 5 per 300s
		_, lastErr = l.SendVerifyEmail(&tenant.SendVerifyEmailRequest{Data: &tenant.SendVerifyEmailRequestData{
			UserType: string(userType),
			Email:    email,
		}})
	}
	require.Error(t, lastErr)
	assert.Contains(t, lastErr.Error(), "try too many times")
}
