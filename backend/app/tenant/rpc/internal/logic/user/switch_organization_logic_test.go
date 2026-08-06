package userlogic

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"zeebee-backend/app/tenant/rpc/internal/logic/testutil"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"
)

// seedMemberInOrgs creates a user who is an active member of the given
// organizations and sets the first one as default. Membership is seeded via
// organization_users rows, which GetUserOrganizations reads from.
func seedMemberInOrgs(t *testing.T, env *testutil.Env, userType constants.UserType, email string, orgs ...*model.Organization) (userId, defaultOrgId int64) {
	t.Helper()
	hash, err := utils.HashPassword(testPassword)
	require.NoError(t, err)
	user := testutil.SeedUser(t, env.DB, &model.User{
		Type:         userType,
		Email:        email,
		PasswordHash: hash,
		Status:       constants.UserStatusActive,
	})
	for _, org := range orgs {
		testutil.SeedOrganizationUser(t, env.DB, &model.OrganizationUser{
			OrganizationId: org.Id,
			UserId:         user.Id,
			Status:         constants.OrganizationUserStatusActive,
		})
	}
	require.NoError(t, env.DB.Model(&model.User{}).Where("id = ?", user.Id).
		Update("default_organization_id", orgs[0].Id).Error)
	return user.Id, orgs[0].Id
}

func seedOrg(t *testing.T, env *testutil.Env, userType constants.UserType, name string, status constants.OrganizationStatus) *model.Organization {
	t.Helper()
	return testutil.SeedOrganization(t, env.DB, &model.Organization{
		Name:   name,
		Status: status,
		Type:   userType,
	})
}

// T032 US-13: switching to a member org re-issues a token carrying the new org.
func TestSwitchOrganization_ReissuesTokenForMember(t *testing.T) {
	env := newTestEnv(t)
	userType := constants.UserTypeLender
	orgA := seedOrg(t, env, userType, "OrgA", constants.OrganizationStatusActive)
	orgB := seedOrg(t, env, userType, "OrgB", constants.OrganizationStatusActive)
	userId, _ := seedMemberInOrgs(t, env, userType, "switch@test.local", orgA, orgB)

	resp, err := NewSwitchOrganizationLogic(ctxWithUserClaim(t, env, userType, userId, orgA.Id), env.SvcCtx).SwitchOrganization(&tenant.SwitchOrganizationRequest{Data: &tenant.SwitchOrganizationRequestData{
		UserId:         userId,
		OrganizationId: orgB.Id,
	}})
	require.NoError(t, err)
	require.NotNil(t, resp.GetData())
	assert.NotEmpty(t, resp.GetData().GetAccessToken())

	claims, err := utils.ParseToken(resp.GetData().GetAccessToken(), env.SvcCtx.Config.JwtConf.AccessSecret)
	require.NoError(t, err)
	assert.Equal(t, userId, claims.UserId)
	assert.Equal(t, orgB.Id, claims.OrganizationId)
	assert.Equal(t, userType, claims.UserType)

	var user model.User
	require.NoError(t, env.DB.First(&user, userId).Error)
	assert.Equal(t, orgB.Id, user.DefaultOrganizationId)
}

// T032 US-13: switching to an organization the user is not a member of is rejected.
func TestSwitchOrganization_RejectsNonMember(t *testing.T) {
	env := newTestEnv(t)
	userType := constants.UserTypeLender
	orgA := seedOrg(t, env, userType, "OrgA", constants.OrganizationStatusActive)
	notMember := seedOrg(t, env, userType, "OrgC", constants.OrganizationStatusActive) // user never joins this one
	userId, _ := seedMemberInOrgs(t, env, userType, "nonmember@test.local", orgA)

	_, err := NewSwitchOrganizationLogic(ctxWithUserClaim(t, env, userType, userId, orgA.Id), env.SvcCtx).SwitchOrganization(&tenant.SwitchOrganizationRequest{Data: &tenant.SwitchOrganizationRequestData{
		UserId:         userId,
		OrganizationId: notMember.Id,
	}})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not member")
}

// T032 US-13: switching to an inactive organization is rejected.
func TestSwitchOrganization_RejectsInactiveOrg(t *testing.T) {
	env := newTestEnv(t)
	userType := constants.UserTypeLender
	orgA := seedOrg(t, env, userType, "OrgA", constants.OrganizationStatusActive)
	inactive := seedOrg(t, env, userType, "OrgInactive", constants.OrganizationStatusDisabled)
	userId, _ := seedMemberInOrgs(t, env, userType, "inactive@test.local", orgA, inactive)

	_, err := NewSwitchOrganizationLogic(ctxWithUserClaim(t, env, userType, userId, orgA.Id), env.SvcCtx).SwitchOrganization(&tenant.SwitchOrganizationRequest{Data: &tenant.SwitchOrganizationRequestData{
		UserId:         userId,
		OrganizationId: inactive.Id,
	}})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not active")
}
