package userlogic

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"zeebee-backend/app/tenant/rpc/internal/logic/testutil"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"
)

// seedAdminOrg seeds an active Admin-type organization (CreateAdminUser targets
// the admin platform; the org must match that type).
func seedAdminOrg(t *testing.T, env *testutil.Env) *model.Organization {
	t.Helper()
	return testutil.SeedOrganization(t, env.DB, &model.Organization{
		Name:   "AdminOrg",
		Status: constants.OrganizationStatusActive,
		Type:   constants.UserTypeAdmin,
	})
}

// T046 US-14: CreateAdminUser provisions an admin user bound to an existing org
// with credentials and roles. End-to-end exercises the fixed role-binding path
// in OrganizationUserLogic.CreateAdminOrganizationUser.
func TestCreateAdminUser_ProvisionsActiveMemberWithRoles(t *testing.T) {
	env := newTestEnv(t)
	organization := seedAdminOrg(t, env)
	role := testutil.SeedRole(t, env.DB, &model.Role{
		OrganizationId: organization.Id, Name: "AdminRole", Type: constants.UserTypeAdmin, Status: constants.CommonStatusActive,
	})

	resp, err := NewCreateAdminUserLogic(context.Background(), env.SvcCtx).CreateAdminUser(&tenant.CreateAdminUserRequest{Data: &tenant.CreateAdminUserRequestData{
		Email: "newadmin@test.local", Password: testPassword, FirstName: "Ad", LastName: "Min",
		OrganizationId: organization.Id, RoleIds: []int64{role.Id},
	}})
	require.NoError(t, err)
	require.NotNil(t, resp.GetData())
	assert.Greater(t, resp.GetData().GetUserId(), int64(0))

	var user model.User
	require.NoError(t, env.DB.Where("email = ?", "newadmin@test.local").First(&user).Error)
	assert.Equal(t, constants.UserTypeAdmin, user.Type)
	assert.Equal(t, constants.UserStatusActive, user.Status)
	assert.NotEmpty(t, user.PasswordHash)

	var member model.OrganizationUser
	require.NoError(t, env.DB.Where("organization_id = ? AND user_id = ?", organization.Id, user.Id).First(&member).Error)
	assert.Equal(t, constants.OrganizationUserStatusActive, member.Status)

	var roleCount int64
	require.NoError(t, env.DB.Table("organization_user_and_roles").Where("organization_user_id = ?", member.Id).Count(&roleCount).Error)
	assert.Equal(t, int64(1), roleCount)
}

// T046 US-14: provisioning an email already in the organization is rejected.
func TestCreateAdminUser_DuplicateRejected(t *testing.T) {
	env := newTestEnv(t)
	organization := seedAdminOrg(t, env)
	l := NewCreateAdminUserLogic(context.Background(), env.SvcCtx)
	req := &tenant.CreateAdminUserRequest{Data: &tenant.CreateAdminUserRequestData{
		Email: "dup@test.local", Password: testPassword, FirstName: "D", LastName: "Up",
		OrganizationId: organization.Id,
	}}
	_, err := l.CreateAdminUser(req)
	require.NoError(t, err)
	_, err = l.CreateAdminUser(req)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "already in the organization")
}

// T046 US-14: a weak password is rejected.
func TestCreateAdminUser_WeakPasswordRejected(t *testing.T) {
	env := newTestEnv(t)
	organization := seedAdminOrg(t, env)

	_, err := NewCreateAdminUserLogic(context.Background(), env.SvcCtx).CreateAdminUser(&tenant.CreateAdminUserRequest{Data: &tenant.CreateAdminUserRequestData{
		Email: "weak@test.local", Password: "alllowercase", OrganizationId: organization.Id,
	}})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "password is invalid")
}
