package userlogic

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/stringx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"
)

// assertEventual polls cond until it holds or the deadline expires. Used for
// async email/limiter side effects that run on a goroutine after the RPC returns.
func assertEventual(t *testing.T, cond func() bool) {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if cond() {
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatal("condition not met within deadline")
}

// T015 US-2: lender registration builds the identity with ThirdParty owner type,
// lowercases the email, and queues the verify email (code cached in Redis).
func TestRegister_CreatesLenderIdentity(t *testing.T) {
	env := newTestEnv(t)
	l := NewRegisterLogic(context.Background(), env.SvcCtx)

	resp, err := l.Register(&tenant.RegisterRequest{Data: &tenant.RegisterRequestData{
		UserType:         string(constants.UserTypeLender),
		Email:            "Lender.New@Test.LOCAL",
		Password:         testPassword,
		FirstName:        "Lender",
		LastName:         "One",
		OrganizationName: "Lender Org",
	}})
	require.NoError(t, err)
	require.NotNil(t, resp.GetData())
	require.Greater(t, resp.GetData().GetUserId(), int64(0))
	assert.Equal(t, string(constants.UserStatusUnverified), resp.GetData().GetStatus())

	var user model.User
	require.NoError(t, env.DB.First(&user, resp.GetData().GetUserId()).Error)
	assert.Equal(t, "lender.new@test.local", user.Email)
	assert.Equal(t, constants.UserTypeLender, user.Type)
	assert.Equal(t, constants.UserStatusUnverified, user.Status)
	assert.Equal(t, resp.GetData().GetOrganizationId(), user.DefaultOrganizationId)

	var org model.Organization
	require.NoError(t, env.DB.First(&org, resp.GetData().GetOrganizationId()).Error)
	assert.Equal(t, constants.OrganizationOwnerTypeThirdParty, org.OwnerType)
	assert.Equal(t, constants.OrganizationStatusNotActivated, org.Status)

	var ou model.OrganizationUser
	require.NoError(t, env.DB.Where("user_id = ? AND organization_id = ?", user.Id, org.Id).First(&ou).Error)
	assert.Equal(t, constants.OrganizationUserStatusNotActivated, ou.Status)

	codeKey := stringx.Format(constants.UserRegisterVerifyCode, constants.UserTypeLender, user.Email)
	assertEventual(t, func() bool {
		code, err := env.SvcCtx.RedisCli.Get(context.Background(), codeKey)
		return err == nil && len(code) == 10
	})
}

// T015 US-2: borrower registration defaults owner type to None and accepts the
// opaque pre-approval JSON (side effects stay at the borrower backend).
func TestRegister_CreatesBorrowerIdentityWithPreApprovalCtx(t *testing.T) {
	env := newTestEnv(t)
	l := NewRegisterLogic(context.Background(), env.SvcCtx)

	resp, err := l.Register(&tenant.RegisterRequest{Data: &tenant.RegisterRequestData{
		UserType:               string(constants.UserTypeBorrower),
		Email:                  "Borrower@Test.LOCAL",
		Password:               testPassword,
		FirstName:              "Borrower",
		LastName:               "Two",
		OrganizationName:       "Borrower Org",
		OwnerType:              "Self", // must be forced to None for borrower
		BorrowerPreApprovalJson: `{"request_id":"req-123","lender_org_id":42}`,
	}})
	require.NoError(t, err)
	require.NotNil(t, resp.GetData())
	require.Greater(t, resp.GetData().GetUserId(), int64(0))

	var user model.User
	require.NoError(t, env.DB.First(&user, resp.GetData().GetUserId()).Error)
	assert.Equal(t, "borrower@test.local", user.Email)
	assert.Equal(t, constants.UserTypeBorrower, user.Type)

	var org model.Organization
	require.NoError(t, env.DB.First(&org, resp.GetData().GetOrganizationId()).Error)
	assert.Equal(t, constants.OrganizationOwnerTypeNone, org.OwnerType)
}

// T015 US-2: registering an unverified user again resends the verify email and
// returns the existing identity instead of a duplicate-creation error.
func TestRegister_ExistingUnverifiedResendsVerifyEmail(t *testing.T) {
	env := newTestEnv(t)
	userId, _ := seedLoginUser(t, env, constants.UserTypeLender, "existing@test.local", testPassword)
	require.NoError(t, env.DB.Model(&model.User{}).
		Where("id = ?", userId).
		Update("status", constants.UserStatusUnverified).Error)

	l := NewRegisterLogic(context.Background(), env.SvcCtx)
	resp, err := l.Register(&tenant.RegisterRequest{Data: &tenant.RegisterRequestData{
		UserType:         string(constants.UserTypeLender),
		Email:            "EXISTING@test.local",
		Password:         testPassword,
		FirstName:        "Existing",
		LastName:         "User",
		OrganizationName: "Existing Org",
	}})
	require.NoError(t, err)
	require.NotNil(t, resp.GetData())
	assert.Equal(t, userId, resp.GetData().GetUserId())
	assert.Equal(t, string(constants.UserStatusUnverified), resp.GetData().GetStatus())

	var count int64
	require.NoError(t, env.DB.Model(&model.User{}).
		Where("email = ?", "existing@test.local").Count(&count).Error)
	assert.Equal(t, int64(1), count)

	codeKey := stringx.Format(constants.UserRegisterVerifyCode, constants.UserTypeLender, "existing@test.local")
	assertEventual(t, func() bool {
		code, err := env.SvcCtx.RedisCli.Get(context.Background(), codeKey)
		return err == nil && len(code) == 10
	})
}

// T015 US-2: registering an already-active user returns the existing identity and
// sends the login-guidance email (rate-limit key recorded on the async path).
func TestRegister_ExistingActiveSendsLoginGuidance(t *testing.T) {
	env := newTestEnv(t)
	userId, _ := seedLoginUser(t, env, constants.UserTypeBorrower, "active@test.local", testPassword)

	l := NewRegisterLogic(context.Background(), env.SvcCtx)
	resp, err := l.Register(&tenant.RegisterRequest{Data: &tenant.RegisterRequestData{
		UserType:         string(constants.UserTypeBorrower),
		Email:            "active@test.local",
		Password:         testPassword,
		FirstName:        "Active",
		LastName:         "User",
		OrganizationName: "Active Org",
	}})
	require.NoError(t, err)
	require.NotNil(t, resp.GetData())
	assert.Equal(t, userId, resp.GetData().GetUserId())
	assert.Equal(t, string(constants.UserStatusActive), resp.GetData().GetStatus())

	limiterKey := stringx.Format(constants.LoginEmail, constants.UserTypeBorrower, "active@test.local")
	assertEventual(t, func() bool {
		// redis_rate stores rate-limiter state under a "rate:" prefix.
		val, err := env.Redis.Get("rate:" + limiterKey)
		return err == nil && val != ""
	})
}
