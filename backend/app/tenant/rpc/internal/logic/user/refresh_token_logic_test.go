package userlogic

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"
)

// T024 US-6: admin refresh returns the full TokenSet (ODD-6) signed with the
// admin keys, and the tokens are newly issued.
func TestRefreshToken_AdminFullTokenSet(t *testing.T) {
	env := newTestEnv(t)
	userType := constants.UserTypeAdmin
	email := "admin-refresh@test.local"

	userId, _ := seedLoginUser(t, env, userType, email, testPassword)
	loginResp, err := NewLoginLogic(context.Background(), env.SvcCtx).Login(&tenant.LoginRequest{Data: &tenant.LoginRequestData{
		UserType: string(userType),
		Email:    email,
		Password: testPassword,
	}})
	require.NoError(t, err)
	require.NotNil(t, loginResp.GetData())

	l := NewRefreshTokenLogic(context.Background(), env.SvcCtx)
	resp, err := l.RefreshToken(&tenant.RefreshTokenRequest{Data: &tenant.RefreshTokenRequestData{
		RefreshToken: loginResp.GetData().GetRefreshToken(),
	}})
	require.NoError(t, err)
	require.NotNil(t, resp.GetData())
	assert.NotEmpty(t, resp.GetData().GetAccessToken())
	assert.NotEmpty(t, resp.GetData().GetRefreshToken())
	assert.Greater(t, resp.GetData().GetAccessTokenExpire(), int64(0))
	assert.Greater(t, resp.GetData().GetRefreshTokenExpire(), int64(0))
	assert.NotEqual(t, loginResp.GetData().GetAccessToken(), resp.GetData().GetAccessToken())
	assert.NotEqual(t, loginResp.GetData().GetRefreshToken(), resp.GetData().GetRefreshToken())

	claims, err := utils.ParseToken(resp.GetData().GetAccessToken(), env.SvcCtx.Config.JwtConf.AccessSecret)
	require.NoError(t, err)
	assert.Equal(t, userId, claims.UserId)
	assert.Equal(t, userType, claims.UserType)
}

// T024 US-6: an inactive user cannot refresh.
func TestRefreshToken_InactiveUser(t *testing.T) {
	env := newTestEnv(t)
	userType := constants.UserTypeLender
	email := "inactive-refresh@test.local"

	userId, orgId := seedLoginUser(t, env, userType, email, testPassword)
	require.NoError(t, env.DB.Model(&model.User{}).
		Where("id = ?", userId).
		Update("status", constants.UserStatusInactive).Error)

	var ouId int64
	require.NoError(t, env.DB.Model(&model.OrganizationUser{}).
		Where("user_id = ? AND organization_id = ?", userId, orgId).
		Select("id").Scan(&ouId).Error)
	tok, err := utils.GenerateToken(utils.JwtUser{
		UserId:             userId,
		OrganizationId:     orgId,
		OrganizationUserId: ouId,
		UserType:           userType,
	}, env.SvcCtx.Config.JwtConf.RefreshSecret, env.SvcCtx.Config.JwtConf.RefreshExpire)
	require.NoError(t, err)

	l := NewRefreshTokenLogic(context.Background(), env.SvcCtx)
	_, err = l.RefreshToken(&tenant.RefreshTokenRequest{Data: &tenant.RefreshTokenRequestData{
		RefreshToken: tok,
	}})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "user is not active")
}
