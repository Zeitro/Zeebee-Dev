package userlogic

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"zeebee-backend/app/admin/rpc/internal/global"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/stringx"
	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"
)

// T012 characterization baseline: admin refresh-token + reset-password error
// paths as they run today (local logic, pre-cutover). Admin has no self-register,
// switch, or verify-email flow — users are provisioned — so only these two are
// characterized. The cutover must keep these green.

func setResetCode(t *testing.T, email, code string) {
	t.Helper()
	key := stringx.Format(constants.ResetPwdCode, constants.UserTypeAdmin, email)
	require.NoError(t, global.GetRedis().SetExpire(context.Background(), key, code, 3600))
}

// --- refresh token ---

// A valid refresh token issues a fresh access token carrying the same user id.
func TestBaselineRefreshToken_Valid(t *testing.T) {
	env := newBaselineEnv(t)
	email := "refresh@test.local"
	userId, orgId := seedAdminUser(t, env, email, baselinePassword)

	refresh, err := utils.GenerateToken(utils.JwtUser{
		UserId:         userId,
		OrganizationId: orgId,
		UserType:       constants.UserTypeAdmin,
	}, env.SvcCtx.Config.JwtConf.RefreshSecret, env.SvcCtx.Config.JwtConf.RefreshExpire)
	require.NoError(t, err)

	l := NewRefreshTokenLogic(context.Background(), env.SvcCtx)
	resp, err := l.RefreshToken(&admin.RefreshTokenRequest{Data: &admin.RefreshTokenRequestData{
		RefreshToken: refresh,
	}})
	require.NoError(t, err)
	require.NotNil(t, resp.GetData())
	assert.NotEmpty(t, resp.GetData().GetAccessToken())

	claims, err := utils.ParseToken(resp.GetData().GetAccessToken(), env.SvcCtx.Config.JwtConf.AccessSecret)
	require.NoError(t, err)
	assert.Equal(t, userId, claims.UserId)
}

// --- reset password ---

// ResetPassword with a valid code issues a token and the new password works.
func TestBaselineResetPassword_ValidCode(t *testing.T) {
	env := newBaselineEnv(t)
	email := "reset@test.local"
	userId, _ := seedAdminUser(t, env, email, baselinePassword)
	setResetCode(t, email, "resetcode1")

	newPwd := "NewPass!2026"
	l := NewResetPasswordLogic(context.Background(), env.SvcCtx)
	resp, err := l.ResetPassword(&admin.ResetPasswordRequest{Data: &admin.ResetPasswordRequestData{
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
	codeKey := stringx.Format(constants.ResetPwdCode, constants.UserTypeAdmin, email)
	stored, _ := global.GetRedis().Get(context.Background(), codeKey)
	assert.Empty(t, stored, "reset code must be consumed after use")
}

// ResetPassword with no code stored is rejected.
func TestBaselineResetPassword_InvalidCode(t *testing.T) {
	env := newBaselineEnv(t)
	email := "reset-invalid@test.local"
	seedAdminUser(t, env, email, baselinePassword)

	l := NewResetPasswordLogic(context.Background(), env.SvcCtx)
	_, err := l.ResetPassword(&admin.ResetPasswordRequest{Data: &admin.ResetPasswordRequestData{
		Email:    email,
		Password: "NewPass!2026",
		Code:     "wrongcode",
	}})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid or expired code")
}
