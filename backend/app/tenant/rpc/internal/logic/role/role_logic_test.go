package rolelogic

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"zeebee-backend/app/tenant/rpc/internal/logic"
	"zeebee-backend/app/tenant/rpc/internal/logic/testutil"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"
)

func newTestEnv(t *testing.T) *testutil.Env {
	t.Helper()
	env := testutil.NewEnv(t)
	logic.ResetBaseLogicForTest()
	return env
}

// seedActiveOrganization builds an active organization (with an owner member)
// that roles can belong to.
func seedActiveOrganization(t *testing.T, env *testutil.Env, userType constants.UserType) *model.Organization {
	t.Helper()
	owner := testutil.SeedUser(t, env.DB, &model.User{
		Type:   userType,
		Email:  "owner-" + string(userType) + "@test.local",
		Status: constants.UserStatusActive,
	})
	organization := testutil.SeedOrganization(t, env.DB, &model.Organization{
		Name:    "Org-" + string(userType),
		Status:  constants.OrganizationStatusActive,
		OwnerId: owner.Id,
		Type:    userType,
	})
	testutil.SeedOrganizationUser(t, env.DB, &model.OrganizationUser{
		OrganizationId: organization.Id,
		UserId:         owner.Id,
		Status:         constants.OrganizationUserStatusActive,
	})
	return organization
}

func seedResource(t *testing.T, env *testutil.Env, userType constants.UserType, name string) *model.Resource {
	t.Helper()
	return testutil.SeedResource(t, env.DB, &model.Resource{
		Type:     userType,
		Status:   constants.CommonStatusActive,
		Name:     name,
		IsPublic: true,
	})
}

// resourceBindings counts role_and_resources rows for a role.
func resourceBindings(t *testing.T, env *testutil.Env, roleId int64) int64 {
	t.Helper()
	var count int64
	require.NoError(t, env.DB.Model(&model.RoleAndResource{}).Where("role_id = ?", roleId).Count(&count).Error)
	return count
}

// createRole is the common Create call used across tests; returns the created role id.
func createRole(t *testing.T, env *testutil.Env, organization *model.Organization, name string, resourceIds []int64) int64 {
	t.Helper()
	relations := make([]*tenant.RoleAndResourceRelationEntity, 0, len(resourceIds))
	for _, rid := range resourceIds {
		relations = append(relations, &tenant.RoleAndResourceRelationEntity{
			ResourceId: rid, DataScope: string(constants.DataScopeAll),
		})
	}
	resp, err := NewCreateLogic(context.Background(), env.SvcCtx).Create(&tenant.RoleCreateRequest{Data: &tenant.RoleCreateRequestData{
		UserType:         string(organization.Type),
		OrganizationId:   organization.Id,
		Name:             name,
		Status:           string(constants.CommonStatusActive),
		RoleAndResources: relations,
	}})
	require.NoError(t, err)
	require.NotNil(t, resp.GetData())
	return resp.GetData().GetId()
}

// T029 US-8: Create adds a role and binds only the explicitly listed resources.
func TestCreate_BindsResources(t *testing.T) {
	env := newTestEnv(t)
	organization := seedActiveOrganization(t, env, constants.UserTypeLender)
	r1 := seedResource(t, env, constants.UserTypeLender, "Res1")
	r2 := seedResource(t, env, constants.UserTypeLender, "Res2")

	roleId := createRole(t, env, organization, "Manager", []int64{r1.Id, r2.Id})

	var role model.Role
	require.NoError(t, env.DB.First(&role, roleId).Error)
	assert.Equal(t, "Manager", role.Name)
	assert.Equal(t, constants.CommonStatusActive, role.Status)
	assert.Equal(t, int64(2), resourceBindings(t, env, roleId))
}

// T029 US-8: a duplicate role name in the same organization is rejected.
func TestCreate_DuplicateNameRejected(t *testing.T) {
	env := newTestEnv(t)
	organization := seedActiveOrganization(t, env, constants.UserTypeLender)
	createRole(t, env, organization, "Manager", nil)

	_, err := NewCreateLogic(context.Background(), env.SvcCtx).Create(&tenant.RoleCreateRequest{Data: &tenant.RoleCreateRequestData{
		UserType:       string(constants.UserTypeLender),
		OrganizationId: organization.Id,
		Name:           "Manager",
		Status:         string(constants.CommonStatusActive),
	}})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
}

func TestCreate_InvalidResourceLeavesNoRole(t *testing.T) {
	env := newTestEnv(t)
	organization := seedActiveOrganization(t, env, constants.UserTypeLender)
	_, err := NewCreateLogic(context.Background(), env.SvcCtx).Create(&tenant.RoleCreateRequest{Data: &tenant.RoleCreateRequestData{
		UserType:       string(constants.UserTypeLender),
		OrganizationId: organization.Id,
		Name:           "Invalid",
		Status:         string(constants.CommonStatusActive),
		RoleAndResources: []*tenant.RoleAndResourceRelationEntity{{
			ResourceId: 99999, DataScope: string(constants.DataScopeAll),
		}},
	}})
	require.Error(t, err)
	var count int64
	require.NoError(t, env.DB.Model(&model.Role{}).Where("organization_id = ?", organization.Id).Count(&count).Error)
	assert.Zero(t, count)
}

func TestCreate_InvalidResourceVariantsLeaveNoRole(t *testing.T) {
	for _, tc := range []struct {
		name     string
		resource *model.Resource
		id       int64
	}{
		{"missing", nil, 99999},
		{"wrong type", &model.Resource{Type: constants.UserTypeBorrower, Name: "Borrower", Status: constants.CommonStatusActive}, 0},
		{"disabled", &model.Resource{Type: constants.UserTypeLender, Name: "Disabled", Status: constants.CommonStatusDisabled}, 0},
		{"outside plan", &model.Resource{Type: constants.UserTypeLender, Name: "Outside", Status: constants.CommonStatusActive}, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			env := newTestEnv(t)
			organization := seedActiveOrganization(t, env, constants.UserTypeLender)
			id := tc.id
			if tc.resource != nil {
				id = testutil.SeedResource(t, env.DB, tc.resource).Id
			}
			_, err := NewCreateLogic(context.Background(), env.SvcCtx).Create(&tenant.RoleCreateRequest{Data: &tenant.RoleCreateRequestData{
				UserType: string(constants.UserTypeLender), OrganizationId: organization.Id, Name: "Invalid", Status: string(constants.CommonStatusActive),
				RoleAndResources: []*tenant.RoleAndResourceRelationEntity{{ResourceId: id, DataScope: string(constants.DataScopeAll)}},
			}})
			require.Error(t, err)
			var count int64
			require.NoError(t, env.DB.Model(&model.Role{}).Where("organization_id = ?", organization.Id).Count(&count).Error)
			assert.Zero(t, count)
		})
	}
}

// T029 US-8: Detail returns the role (org + platform scoped).
func TestDetail_ReturnsRole(t *testing.T) {
	env := newTestEnv(t)
	organization := seedActiveOrganization(t, env, constants.UserTypeLender)
	roleId := createRole(t, env, organization, "Detail", nil)

	resp, err := NewDetailLogic(context.Background(), env.SvcCtx).Detail(&tenant.RoleDetailRequest{Data: &tenant.RoleDetailRequestData{
		UserType: string(constants.UserTypeLender), OrganizationId: organization.Id, Id: roleId,
	}})
	require.NoError(t, err)
	require.NotNil(t, resp.GetData())
	assert.Equal(t, roleId, resp.GetData().GetId())
	assert.Equal(t, "Detail", resp.GetData().GetName())
}

// T029 US-8: Update edits fields and rebinds resources (drop to one).
func TestUpdate_RebindsResources(t *testing.T) {
	env := newTestEnv(t)
	organization := seedActiveOrganization(t, env, constants.UserTypeLender)
	r1 := seedResource(t, env, constants.UserTypeLender, "UR1")
	r2 := seedResource(t, env, constants.UserTypeLender, "UR2")
	roleId := createRole(t, env, organization, "Updatable", []int64{r1.Id, r2.Id})
	require.Equal(t, int64(2), resourceBindings(t, env, roleId))

	_, err := NewUpdateLogic(context.Background(), env.SvcCtx).Update(&tenant.RoleUpdateRequest{Data: &tenant.RoleUpdateRequestData{
		UserType: string(constants.UserTypeLender), OrganizationId: organization.Id, Id: roleId,
		Name: "Renamed", Description: "desc",
		RoleAndResources: []*tenant.RoleAndResourceRelationEntity{{
			ResourceId: r1.Id, DataScope: string(constants.DataScopeAll),
		}},
	}})
	require.NoError(t, err)

	var role model.Role
	require.NoError(t, env.DB.First(&role, roleId).Error)
	assert.Equal(t, "Renamed", role.Name)
	assert.Equal(t, "desc", role.Description)
	assert.Equal(t, int64(1), resourceBindings(t, env, roleId), "resource set should shrink to one")
}

// T029 US-8: SwitchStatus disables a role (direction B: next AuthzCheck reflects
// it without invalidation broadcast; the live-query effect is covered by the
// auth package's TestAuthzCheck_ImmediateEffectAfterRoleChange).
func TestSwitchStatus_DisablesRole(t *testing.T) {
	env := newTestEnv(t)
	organization := seedActiveOrganization(t, env, constants.UserTypeLender)
	roleId := createRole(t, env, organization, "Toggle", nil)

	_, err := NewSwitchStatusLogic(context.Background(), env.SvcCtx).SwitchStatus(&tenant.RoleSwitchStatusRequest{Data: &tenant.RoleSwitchStatusRequestData{
		UserType: string(constants.UserTypeLender), OrganizationId: organization.Id, Id: roleId,
		Status: string(constants.CommonStatusDisabled),
	}})
	require.NoError(t, err)

	var role model.Role
	require.NoError(t, env.DB.First(&role, roleId).Error)
	assert.Equal(t, constants.CommonStatusDisabled, role.Status)
}

// T029 US-8: Page lists the organization's roles.
func TestPage_ListsRoles(t *testing.T) {
	env := newTestEnv(t)
	organization := seedActiveOrganization(t, env, constants.UserTypeLender)
	createRole(t, env, organization, "Alpha", nil)
	createRole(t, env, organization, "Beta", nil)

	resp, err := NewPageLogic(context.Background(), env.SvcCtx).Page(&tenant.RolePageRequest{Data: &tenant.RolePageRequestData{
		UserType: string(constants.UserTypeLender), OrganizationId: organization.Id, Page: 1, Limit: 10,
	}})
	require.NoError(t, err)
	require.NotNil(t, resp.GetData())
	assert.Equal(t, int64(2), resp.GetData().GetTotal())
	assert.Len(t, resp.GetData().GetList(), 2)
}

// T029 US-8: All returns every role of an organization.
func TestAll_ReturnsOrgRoles(t *testing.T) {
	env := newTestEnv(t)
	organization := seedActiveOrganization(t, env, constants.UserTypeLender)
	createRole(t, env, organization, "One", nil)
	createRole(t, env, organization, "Two", nil)

	resp, err := NewAllLogic(context.Background(), env.SvcCtx).All(&tenant.RoleAllRequest{Data: &tenant.RoleAllRequestData{
		UserType: string(constants.UserTypeLender), OrganizationId: organization.Id,
	}})
	require.NoError(t, err)
	assert.Len(t, resp.GetData(), 2)
}

// T029 US-8: AssignOrganizationUsers binds members onto a role.
func TestAssignOrganizationUsers_BindsMembers(t *testing.T) {
	env := newTestEnv(t)
	organization := seedActiveOrganization(t, env, constants.UserTypeLender)
	roleId := createRole(t, env, organization, "Assignable", nil)
	u1 := testutil.SeedUser(t, env.DB, &model.User{Type: constants.UserTypeLender, Email: "au1@test.local", Status: constants.UserStatusActive})
	m1 := testutil.SeedOrganizationUser(t, env.DB, &model.OrganizationUser{OrganizationId: organization.Id, UserId: u1.Id, Status: constants.OrganizationUserStatusActive})
	u2 := testutil.SeedUser(t, env.DB, &model.User{Type: constants.UserTypeLender, Email: "au2@test.local", Status: constants.UserStatusActive})
	m2 := testutil.SeedOrganizationUser(t, env.DB, &model.OrganizationUser{OrganizationId: organization.Id, UserId: u2.Id, Status: constants.OrganizationUserStatusActive})

	_, err := NewAssignOrganizationUsersLogic(context.Background(), env.SvcCtx).AssignOrganizationUsers(&tenant.RoleAssignOrganizationUsersRequest{Data: &tenant.RoleAssignOrganizationUsersRequestData{
		UserType: string(constants.UserTypeLender), OrganizationId: organization.Id, Id: roleId,
		OrganizationUserIds: []int64{m1.Id, m2.Id},
	}})
	require.NoError(t, err)

	var count int64
	require.NoError(t, env.DB.Table("organization_user_and_roles").Where("role_id = ?", roleId).Count(&count).Error)
	assert.Equal(t, int64(2), count)
}

// T029 US-8: Delete removes a role.
func TestDelete_RemovesRole(t *testing.T) {
	env := newTestEnv(t)
	organization := seedActiveOrganization(t, env, constants.UserTypeLender)
	roleId := createRole(t, env, organization, "Deletable", nil)
	require.NoError(t, env.DB.Create(&model.RoleResourceSnapshot{OrganizationId: organization.Id, PricingPlanId: 1, RoleId: roleId, Resources: []model.RoleResourceSnapshotResource{{ResourceId: 1, DataScope: constants.DataScopeAll}}}).Error)

	_, err := NewDeleteLogic(context.Background(), env.SvcCtx).Delete(&tenant.RoleDeleteRequest{Data: &tenant.RoleDeleteRequestData{
		UserType: string(constants.UserTypeLender), OrganizationId: organization.Id, Id: roleId,
	}})
	require.NoError(t, err)

	var count int64
	require.NoError(t, env.DB.Model(&model.Role{}).Where("id = ?", roleId).Count(&count).Error)
	assert.Zero(t, count)
	require.NoError(t, env.DB.Model(&model.RoleResourceSnapshot{}).Where("role_id = ?", roleId).Count(&count).Error)
	assert.Zero(t, count)
}
