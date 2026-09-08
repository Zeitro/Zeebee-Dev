package userlogic

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"

	"zeebee-backend/app/tenant/rpc/internal/logic"
	"zeebee-backend/app/tenant/rpc/internal/logic/testutil"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"
)

const testPassword = "Passw0rd!2026"

func newTestEnv(t *testing.T) *testutil.Env {
	t.Helper()
	env := testutil.NewEnv(t)
	logic.ResetBaseLogicForTest()
	return env
}

func seedLoginUser(t *testing.T, env *testutil.Env, userType constants.UserType, email, password string) (int64, int64) {
	t.Helper()
	hash, err := utils.HashPassword(password)
	require.NoError(t, err)
	user := testutil.SeedUser(t, env.DB, &model.User{
		Type:         userType,
		Email:        email,
		PasswordHash: hash,
		Status:       constants.UserStatusActive,
	})
	org := testutil.SeedOrganization(t, env.DB, &model.Organization{
		Name:   "Org-" + string(userType),
		Status: constants.OrganizationStatusActive,
		Type:   userType,
	})
	testutil.SeedOrganizationUser(t, env.DB, &model.OrganizationUser{
		OrganizationId: org.Id,
		UserId:         user.Id,
		Status:         constants.OrganizationUserStatusActive,
	})
	require.NoError(t, env.DB.Model(&model.User{}).
		Where("id = ?", user.Id).
		Update("default_organization_id", org.Id).Error)
	return user.Id, org.Id
}

// ctxWithUserClaim signs a JWT for the fixture user and carries it in gRPC
// metadata, as a backend would forward it to tenant.
func ctxWithUserClaim(t *testing.T, env *testutil.Env, userType constants.UserType, userId, orgId int64) context.Context {
	t.Helper()
	tok, err := utils.GenerateToken(utils.JwtUser{
		UserId:         userId,
		OrganizationId: orgId,
		UserType:       userType,
	}, env.SvcCtx.Config.JwtConf.AccessSecret, 3600)
	require.NoError(t, err)
	return metadata.NewIncomingContext(context.Background(), metadata.Pairs(
		utils.JwtTokenCtxKey, tok,
		utils.UserTypeCtxKey, string(userType),
	))
}

// T013 US-1: three user types log in and receive a token set signed with the
// matching platform key; wrong password fails; repeated failures are throttled.
func TestLogin_SuccessThreeTypes(t *testing.T) {
	env := newTestEnv(t)

	userTypes := []struct {
		name     string
		userType constants.UserType
		email    string
	}{
		{"borrower", constants.UserTypeBorrower, "borrower@test.local"},
		{"lender", constants.UserTypeLender, "lender@test.local"},
		{"admin", constants.UserTypeAdmin, "admin@test.local"},
	}
	for _, tc := range userTypes {
		t.Run(tc.name, func(t *testing.T) {
			userId, _ := seedLoginUser(t, env, tc.userType, tc.email, testPassword)

			l := NewLoginLogic(context.Background(), env.SvcCtx)
			resp, err := l.Login(&tenant.LoginRequest{Data: &tenant.LoginRequestData{
				UserType: string(tc.userType),
				Email:    tc.email,
				Password: testPassword,
			}})
			require.NoError(t, err)
			require.NotNil(t, resp.GetData())
			assert.NotEmpty(t, resp.GetData().GetAccessToken())
			assert.NotEmpty(t, resp.GetData().GetRefreshToken())
			assert.Greater(t, resp.GetData().GetAccessTokenExpire(), int64(0))

			secret := env.SvcCtx.Config.JwtConf.AccessSecret
			claims, err := utils.ParseToken(resp.GetData().GetAccessToken(), secret)
			require.NoError(t, err)
			assert.Equal(t, userId, claims.UserId)
			assert.Equal(t, tc.userType, claims.UserType)
		})
	}
}

func TestGetClaim_ValidatesTokenType(t *testing.T) {
	env := newTestEnv(t)
	token, err := utils.GenerateToken(utils.JwtUser{
		UserId: 7, OrganizationId: 8, OrganizationUserId: 9, Username: "Ada",
		Role: []string{"admin"}, UserType: constants.UserTypeBorrower,
	}, env.SvcCtx.Config.JwtConf.AccessSecret, 3600)
	require.NoError(t, err)

	resp, err := NewGetClaimLogic(context.Background(), env.SvcCtx).GetClaim(&tenant.GetClaimRequest{Data: &tenant.GetClaimRequestData{
		AccessToken: token, UserType: string(constants.UserTypeBorrower),
	}})
	require.NoError(t, err)
	assert.Equal(t, int64(7), resp.Data.UserId)
	assert.Equal(t, int64(8), resp.Data.OrganizationId)
	assert.Equal(t, []string{"admin"}, resp.Data.Role)

	_, err = NewGetClaimLogic(context.Background(), env.SvcCtx).GetClaim(&tenant.GetClaimRequest{Data: &tenant.GetClaimRequestData{
		AccessToken: token, UserType: string(constants.UserTypeLender),
	}})
	require.Error(t, err)
}

func TestLogin_WrongPasswordAndThrottle(t *testing.T) {
	env := newTestEnv(t)
	email := "throttle@test.local"
	seedLoginUser(t, env, constants.UserTypeBorrower, email, testPassword)

	l := NewLoginLogic(context.Background(), env.SvcCtx)
	req := &tenant.LoginRequest{Data: &tenant.LoginRequestData{
		UserType: string(constants.UserTypeBorrower),
		Email:    email,
		Password: "WrongPass!1",
	}}

	var lastErr error
	for i := 0; i < 10; i++ {
		_, lastErr = l.Login(req)
		require.Error(t, lastErr)
	}
	assert.Contains(t, lastErr.Error(), "too many failed login attempts")
}

func TestLogin_UnknownUser(t *testing.T) {
	env := newTestEnv(t)
	l := NewLoginLogic(context.Background(), env.SvcCtx)
	_, err := l.Login(&tenant.LoginRequest{Data: &tenant.LoginRequestData{
		UserType: string(constants.UserTypeLender),
		Email:    "nobody@test.local",
		Password: testPassword,
	}})
	require.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "user not found"))
}
