package userlogic

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"

	"zeebee-backend/app/lender/rpc/internal/logic/testutil"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/proto"
)

// T012 characterization baseline: lender register / refresh / switch-organization
// as they run today (local logic, pre-cutover). The cutover must keep these green.

// --- register ---

// Registering a new lender creates an unverified user plus a ThirdParty org.
func TestBaselineRegister_NewLender(t *testing.T) {
	env := newBaselineEnv(t)
	email := "newlender@test.local"
	l := NewRegisterLogic(context.Background(), env.SvcCtx)
	_, err := l.Register(&lender.UserRegisterRequest{Data: &lender.UserRegisterRequestData{
		Email:            email,
		Password:         baselinePassword,
		FirstName:        "New",
		LastName:         "Lender",
		OrganizationName: "New Lender Co",
	}})
	require.NoError(t, err)

	var user model.User
	require.NoError(t, env.DB.Where("email = ? AND type = ?", email, constants.UserTypeLender).First(&user).Error)
	assert.Equal(t, constants.UserStatusUnverified, user.Status)
	assert.Equal(t, constants.UserTypeLender, user.Type)

	var org model.Organization
	require.NoError(t, env.DB.Where("type = ?", constants.UserTypeLender).First(&org).Error)
	assert.Equal(t, constants.OrganizationOwnerTypeThirdParty, org.OwnerType)
	assert.Equal(t, constants.UserTypeLender, org.Type)
}

// Re-registering an existing active user does not duplicate and returns no error.
func TestBaselineRegister_AlreadyActiveNoDuplicate(t *testing.T) {
	env := newBaselineEnv(t)
	email := "active@test.local"
	seedLenderUser(t, env, email, baselinePassword) // seeds an active user

	var before int64
	require.NoError(t, env.DB.Model(&model.User{}).Count(&before).Error)

	l := NewRegisterLogic(context.Background(), env.SvcCtx)
	_, err := l.Register(&lender.UserRegisterRequest{Data: &lender.UserRegisterRequestData{
		Email:            email,
		Password:         baselinePassword,
		FirstName:        "Active",
		LastName:         "User",
		OrganizationName: "Whatever",
	}})
	require.NoError(t, err)

	var after int64
	require.NoError(t, env.DB.Model(&model.User{}).Count(&after).Error)
	assert.Equal(t, before, after, "re-registering an active lender must not create a duplicate user")
}

// --- refresh token ---

// A valid refresh token issues a fresh token set.
func TestBaselineRefreshToken_Valid(t *testing.T) {
	env := newBaselineEnv(t)
	email := "refresh@test.local"
	userId, orgId := seedLenderUser(t, env, email, baselinePassword)

	refresh, err := utils.GenerateToken(utils.JwtUser{
		UserId:         userId,
		OrganizationId: orgId,
		UserType:       constants.UserTypeLender,
	}, env.SvcCtx.Config.JwtConf.RefreshSecret, env.SvcCtx.Config.JwtConf.RefreshExpire)
	require.NoError(t, err)

	l := NewRefreshTokenLogic(context.Background(), env.SvcCtx)
	resp, err := l.RefreshToken(&proto.RefreshTokenRequest{Data: &proto.RefreshTokenRequestData{
		RefreshToken: refresh,
	}})
	require.NoError(t, err)
	require.NotNil(t, resp.GetData())
	assert.NotEmpty(t, resp.GetData().GetAccessToken())

	claims, err := utils.ParseToken(resp.GetData().GetAccessToken(), env.SvcCtx.Config.JwtConf.AccessSecret)
	require.NoError(t, err)
	assert.Equal(t, userId, claims.UserId)
}

// --- switch organization ---

// Switching to an org the user belongs to reissues a token scoped to that org.
func TestBaselineSwitchOrganization_Member(t *testing.T) {
	env := newBaselineEnv(t)
	email := "switch@test.local"
	userId, org1Id, org2Id := seedLenderUserWithTwoOrgs(t, env, email)

	ctx := ctxWithLenderClaim(t, env, userId, org1Id)
	l := NewSwitchOrganizationLogic(ctx, env.SvcCtx)
	resp, err := l.SwitchOrganization(&proto.SwitchOrganizationRequest{Data: &proto.SwitchOrganizationRequestData{
		OrganizationId: org2Id,
	}})
	require.NoError(t, err)
	require.NotNil(t, resp.GetData())

	claims, err := utils.ParseToken(resp.GetData().GetAccessToken(), env.SvcCtx.Config.JwtConf.AccessSecret)
	require.NoError(t, err)
	assert.NotEqual(t, org1Id, claims.OrganizationId, "token must be re-scoped to the switched org")
}

// Switching to an org the user does NOT belong to is rejected.
func TestBaselineSwitchOrganization_NonMember(t *testing.T) {
	env := newBaselineEnv(t)
	email := "switch2@test.local"
	userId, org1Id, _ := seedLenderUserWithTwoOrgs(t, env, email)

	// A third org the user has no membership in.
	other := &model.Organization{Name: "Other-" + email, Status: constants.OrganizationStatusActive, Type: constants.UserTypeLender}
	require.NoError(t, env.DB.Create(other).Error)

	ctx := ctxWithLenderClaim(t, env, userId, org1Id)
	l := NewSwitchOrganizationLogic(ctx, env.SvcCtx)
	_, err := l.SwitchOrganization(&proto.SwitchOrganizationRequest{Data: &proto.SwitchOrganizationRequestData{
		OrganizationId: other.Id,
	}})
	require.Error(t, err)
}

// --- helpers ---

func ctxWithLenderClaim(t *testing.T, env *testutil.Env, userId, orgId int64) context.Context {
	t.Helper()
	tok, err := utils.GenerateToken(utils.JwtUser{
		UserId:         userId,
		OrganizationId: orgId,
		UserType:       constants.UserTypeLender,
	}, env.SvcCtx.Config.JwtConf.AccessSecret, 3600)
	require.NoError(t, err)
	return metadata.NewIncomingContext(context.Background(), metadata.Pairs(utils.JwtTokenCtxKey, tok))
}

// seedLenderUserWithTwoOrgs seeds one lender user belonging to two orgs
// (default = org1) and returns (userId, org1Id, org2Id).
func seedLenderUserWithTwoOrgs(t *testing.T, env *testutil.Env, email string) (int64, int64, int64) {
	t.Helper()
	userId, org1Id := seedLenderUser(t, env, email, baselinePassword)
	org2 := &model.Organization{Name: "Org2-" + email, Status: constants.OrganizationStatusActive, Type: constants.UserTypeLender}
	require.NoError(t, env.DB.Create(org2).Error)
	require.NoError(t, env.DB.Create(&model.OrganizationUser{
		OrganizationId: org2.Id,
		UserId:         userId,
		Status:         constants.OrganizationUserStatusActive,
	}).Error)
	return userId, org1Id, org2.Id
}
