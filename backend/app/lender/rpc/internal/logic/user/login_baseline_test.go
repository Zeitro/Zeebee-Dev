package userlogic

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/app/lender/rpc/internal/logic/testutil"
	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/proto"
)

// T012 characterization baseline: lender login behavior as it runs today
// (local logic, pre-cutover). After the US-12 cutover routes identity through
// tenant, these baselines must still pass — that is the no-regression gate.

const baselinePassword = "Passw0rd!2026"

func newBaselineEnv(t *testing.T) *testutil.Env {
	t.Helper()
	env := testutil.NewEnv(t)
	logic.ResetBaseLogicForTest()
	return env
}

// seedLenderUser creates an active lender user + org + membership and returns
// (userId, orgId). Mirrors the tenant login fixture.
//
//nolint:unparam // password is constant across callers but kept for clarity
func seedLenderUser(t *testing.T, env *testutil.Env, email, password string) (int64, int64) {
	t.Helper()
	hash, err := utils.HashPassword(password)
	require.NoError(t, err)
	user := &model.User{
		Type:         constants.UserTypeLender,
		Email:        email,
		PasswordHash: hash,
		Status:       constants.UserStatusActive,
	}
	require.NoError(t, env.DB.Create(user).Error)
	org := &model.Organization{
		Name:   "Org-" + email,
		Status: constants.OrganizationStatusActive,
		Type:   constants.UserTypeLender,
	}
	require.NoError(t, env.DB.Create(org).Error)
	require.NoError(t, env.DB.Create(&model.OrganizationUser{
		OrganizationId: org.Id,
		UserId:         user.Id,
		Status:         constants.OrganizationUserStatusActive,
	}).Error)
	require.NoError(t, env.DB.Model(&model.User{}).
		Where("id = ?", user.Id).
		Update("default_organization_id", org.Id).Error)
	return user.Id, org.Id
}

// A lender with valid credentials receives a token set carrying their id and
// userType.
func TestBaselineLogin_Success(t *testing.T) {
	env := newBaselineEnv(t)
	email := "lender@test.local"
	userId, _ := seedLenderUser(t, env, email, baselinePassword)

	l := NewLoginLogic(context.Background(), env.SvcCtx)
	resp, err := l.Login(&proto.UserLoginRequest{Data: &proto.UserLoginRequestData{
		Email:    email,
		Password: baselinePassword,
	}})
	require.NoError(t, err)
	require.NotNil(t, resp.GetData())
	assert.NotEmpty(t, resp.GetData().GetAccessToken())
	assert.NotEmpty(t, resp.GetData().GetRefreshToken())
	assert.Greater(t, resp.GetData().GetAccessTokenExpire(), int64(0))

	claims, err := utils.ParseToken(resp.GetData().GetAccessToken(), env.SvcCtx.Config.JwtConf.AccessSecret)
	require.NoError(t, err)
	assert.Equal(t, userId, claims.UserId)
	assert.Equal(t, constants.UserTypeLender, claims.UserType)
}

// Unknown email is rejected with a not-found error.
func TestBaselineLogin_UnknownUser(t *testing.T) {
	env := newBaselineEnv(t)
	l := NewLoginLogic(context.Background(), env.SvcCtx)
	_, err := l.Login(&proto.UserLoginRequest{Data: &proto.UserLoginRequestData{
		Email:    "nobody@test.local",
		Password: baselinePassword,
	}})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")
}

// Repeated wrong-password attempts are eventually throttled.
func TestBaselineLogin_WrongPasswordThrottled(t *testing.T) {
	env := newBaselineEnv(t)
	email := "throttle@test.local"
	seedLenderUser(t, env, email, baselinePassword)

	l := NewLoginLogic(context.Background(), env.SvcCtx)
	req := &proto.UserLoginRequest{Data: &proto.UserLoginRequestData{
		Email:    email,
		Password: "WrongPass!2026",
	}}

	throttled := false
	for range 15 {
		_, err := l.Login(req)
		require.Error(t, err)
		if strings.Contains(err.Error(), "too many failed login attempts") {
			throttled = true
			break
		}
	}
	assert.True(t, throttled, "expected lender login to throttle after repeated wrong passwords")
}
