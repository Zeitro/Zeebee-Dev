package auth

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
	"gorm.io/gorm"

	"zeebee-backend/common/redisx"
	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/repo"
)

// setupAuthTestDB mirrors the logic-package sqlite fixture for every table the
// rebuild and authorization-check paths read. Each test gets an isolated
// in-memory DB.
func setupAuthTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	require.NoError(t, err)

	for _, s := range []string{
		`CREATE TABLE organizations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			is_del INTEGER DEFAULT 0,
			name TEXT NOT NULL DEFAULT '',
			status TEXT NOT NULL DEFAULT 'Active',
			type TEXT DEFAULT ''
		)`,
		`CREATE TABLE resources (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			is_del INTEGER DEFAULT 0,
			parent_id INTEGER DEFAULT 0,
			type TEXT DEFAULT 'Lender',
			resource_type TEXT DEFAULT 'Menu',
			name TEXT DEFAULT '',
			path TEXT DEFAULT '',
			is_public INTEGER DEFAULT 0,
			status TEXT DEFAULT 'Active',
			sort_order INTEGER DEFAULT 10,
			description TEXT DEFAULT ''
		)`,
		`CREATE TABLE apis (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			is_del INTEGER DEFAULT 0,
			parent_id INTEGER DEFAULT 0,
			type TEXT DEFAULT 'Lender',
			name TEXT DEFAULT '',
			path TEXT DEFAULT '',
			method TEXT DEFAULT '',
			description TEXT DEFAULT '',
			status TEXT DEFAULT 'Active',
			sort_order INTEGER DEFAULT 10,
			auth_required INTEGER DEFAULT 0,
			permission_required INTEGER DEFAULT 0
		)`,
		`CREATE TABLE resource_and_apis (
			resource_id INTEGER NOT NULL,
			api_id INTEGER NOT NULL,
			PRIMARY KEY (resource_id, api_id)
		)`,
		`CREATE TABLE roles (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			is_del INTEGER DEFAULT 0,
			organization_id INTEGER NOT NULL DEFAULT 0,
			name TEXT NOT NULL DEFAULT '',
			type TEXT NOT NULL DEFAULT 'Lender',
			status TEXT NOT NULL DEFAULT 'Active',
			description TEXT NOT NULL DEFAULT ''
		)`,
		`CREATE TABLE role_and_resources (
			role_id INTEGER NOT NULL,
			resource_id INTEGER NOT NULL,
			PRIMARY KEY (role_id, resource_id)
		)`,
		`CREATE TABLE organization_users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			is_del INTEGER DEFAULT 0,
			organization_id INTEGER DEFAULT 0,
			user_id INTEGER DEFAULT 0,
			status TEXT DEFAULT 'Not Activated',
			first_name TEXT DEFAULT '',
			mid_name TEXT DEFAULT '',
			last_name TEXT DEFAULT '',
			full_name TEXT DEFAULT '',
			phone TEXT DEFAULT '',
			remark TEXT DEFAULT ''
		)`,
		`CREATE TABLE organization_user_and_roles (
			organization_user_id INTEGER NOT NULL,
			role_id INTEGER NOT NULL,
			PRIMARY KEY (organization_user_id, role_id)
		)`,
		`CREATE TABLE pricing_plans (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			is_del INTEGER DEFAULT 0,
			name TEXT NOT NULL DEFAULT '',
			status TEXT NOT NULL DEFAULT 'Active',
			type TEXT NOT NULL DEFAULT 'Lender',
			plan_type TEXT NOT NULL DEFAULT 'Standard',
			sort_order INTEGER NOT NULL DEFAULT 6,
			is_public INTEGER NOT NULL DEFAULT 0,
			is_default INTEGER NOT NULL DEFAULT 0,
			trial_days INTEGER NOT NULL DEFAULT 0
		)`,
		`CREATE TABLE pricing_plan_and_resources (
			pricing_plan_id INTEGER NOT NULL DEFAULT 0,
			resource_id INTEGER NOT NULL DEFAULT 0
		)`,
		`CREATE TABLE organization_subscriptions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			is_del INTEGER DEFAULT 0,
			code TEXT NOT NULL DEFAULT '' UNIQUE,
			organization_id INTEGER NOT NULL DEFAULT 0,
			pricing_plan_id INTEGER NOT NULL DEFAULT 0,
			pricing_plan_policy_id INTEGER NOT NULL DEFAULT 0,
			status TEXT NOT NULL DEFAULT 'Active',
			expire_date DATE,
			original_amount REAL NOT NULL DEFAULT 0,
			actual_amount REAL NOT NULL DEFAULT 0,
			is_auto_renew INTEGER NOT NULL DEFAULT 0
		)`,
	} {
		require.NoError(t, db.Exec(s).Error)
	}
	return db
}

// newAuthForTest wires an Auth backed by sqlite repos and a fake in-memory
// enforcer. redisCli is nil: used by rebuild-path tests, which never touch the
// role cache.
func newAuthForTest(t *testing.T) (*Auth, *fakeEnforcer, *gorm.DB) {
	t.Helper()
	db := setupAuthTestDB(t)
	repos := repo.NewAllRepos(db)
	fake := newFakeEnforcer()
	a := newAuthWithEnforcer(repos, nil, nil, fake, "")
	return a, fake, db
}

// newAuthzForTest wires an Auth backed by sqlite + miniredis + fake enforcer.
// miniredis is required by CheckAuthorization: getOrganizationUserRoles
// dereferences l.redisCli.
func newAuthzForTest(t *testing.T) (*Auth, *fakeEnforcer, *gorm.DB, *miniredis.Miniredis) {
	t.Helper()
	db := setupAuthTestDB(t)
	mr := miniredis.RunT(t)
	rds, err := redisx.GetRedis(redisx.Opt{Host: mr.Addr()})
	require.NoError(t, err)
	repos := repo.NewAllRepos(db)
	fake := newFakeEnforcer()
	a := newAuthWithEnforcer(repos, rds, map[constants.UserType]dto.JwtConfig{constants.UserTypeLender: {AccessSecret: "secret"}}, fake, "")
	return a, fake, db, mr
}

type authFixture struct {
	orgID    int64
	roleID   int64
	resID    int64
	planID   int64
	userType constants.UserType
}

// seedGranted seeds org 1, an active lender role, a lender resource, the
// resource bound to the role, an active plan covering the resource, and an
// Active subscription. Baseline under which a rebuild MUST grant (roleId,
// resourceId).
func seedGranted(t *testing.T, db *gorm.DB) authFixture {
	t.Helper()
	const ut = constants.UserTypeLender
	require.NoError(t, db.Exec(`INSERT INTO organizations (id, name, status, type) VALUES (1, 'Org1', 'Active', ?)`, ut).Error)
	require.NoError(t, db.Exec(`INSERT INTO resources (id, type, status, name) VALUES (100, ?, 'Active', 'Res100')`, ut).Error)
	require.NoError(t, db.Exec(`INSERT INTO roles (id, organization_id, name, type, status) VALUES (10, 1, 'LenderRole', ?, 'Active')`, ut).Error)
	require.NoError(t, db.Exec(`INSERT INTO role_and_resources (role_id, resource_id) VALUES (10, 100)`).Error)
	require.NoError(t, db.Exec(`INSERT INTO pricing_plans (id, name, status, type) VALUES (7, 'Plan7', 'Active', ?)`, ut).Error)
	require.NoError(t, db.Exec(`INSERT INTO pricing_plan_and_resources (pricing_plan_id, resource_id) VALUES (7, 100)`).Error)
	require.NoError(t, db.Exec(`INSERT INTO organization_subscriptions (code, organization_id, pricing_plan_id, status, expire_date) VALUES ('SU1', 1, 7, 'Active', ?)`, time.Now().AddDate(0, 0, 10)).Error)
	return authFixture{orgID: 1, roleID: 10, resID: 100, planID: 7, userType: ut}
}

func enforced(t *testing.T, fake *fakeEnforcer, f authFixture) bool {
	t.Helper()
	ok, err := fake.Enforce(BuildInterfacePolicy(f.roleID, f.resID, f.userType)...)
	require.NoError(t, err)
	return ok
}

const (
	authzMethod = "GET"
	authzPath   = "/api/v1/widget"
)

// seedAuthzBaseline builds org 1, resource 100 bound to api (GET /api/v1/widget),
// an admin role (id=10) and a normal staff role (id=11), org-user 50 carrying
// ONLY the staff role, plan 7 covering resource 100, and an Active subscription.
func seedAuthzBaseline(t *testing.T, db *gorm.DB) {
	t.Helper()
	const ut = constants.UserTypeLender
	require.NoError(t, db.Exec(`INSERT INTO organizations (id, name, status, type) VALUES (1, 'Org1', 'Active', ?)`, ut).Error)
	require.NoError(t, db.Exec(`INSERT INTO resources (id, type, status, name) VALUES (100, ?, 'Active', 'Res100')`, ut).Error)
	require.NoError(t, db.Exec(`INSERT INTO apis (id, type, status, name, path, method, permission_required) VALUES (1000, ?, 'Active', 'Widget', ?, ?, 1)`, ut, authzPath, authzMethod).Error)
	require.NoError(t, db.Exec(`INSERT INTO resource_and_apis (resource_id, api_id) VALUES (100, 1000)`).Error)

	require.NoError(t, db.Exec(`INSERT INTO roles (id, organization_id, name, type, status) VALUES (10, 1, ?, ?, 'Active')`, constants.AdminRoleName, ut).Error)
	require.NoError(t, db.Exec(`INSERT INTO roles (id, organization_id, name, type, status) VALUES (11, 1, 'Staff', ?, 'Active')`, ut).Error)

	require.NoError(t, db.Exec(`INSERT INTO organization_users (id, organization_id, user_id, status) VALUES (50, 1, 1, 'Active')`).Error)
	require.NoError(t, db.Exec(`INSERT INTO organization_user_and_roles (organization_user_id, role_id) VALUES (50, 11)`).Error)

	require.NoError(t, db.Exec(`INSERT INTO pricing_plans (id, name, status, type) VALUES (7, 'Plan7', 'Active', ?)`, ut).Error)
	require.NoError(t, db.Exec(`INSERT INTO pricing_plan_and_resources (pricing_plan_id, resource_id) VALUES (7, 100)`).Error)
	require.NoError(t, db.Exec(`INSERT INTO organization_subscriptions (code, organization_id, pricing_plan_id, status, expire_date) VALUES ('SU1', 1, 7, 'Active', ?)`, time.Now().AddDate(0, 0, 10)).Error)
}

// ctxWithClaim builds a grpc-metadata context carrying a signed JWT so
// GetUserClaimFromCtx / IsAdmin resolve inside CheckAuthorization.
func ctxWithClaim(t *testing.T, a *Auth, userType constants.UserType, orgID, orgUserID int64) context.Context {
	t.Helper()
	jc, err := a.JwtConfigFor(userType)
	require.NoError(t, err)
	tok, err := utils.GenerateToken(utils.JwtUser{
		UserId:             1,
		OrganizationId:     orgID,
		OrganizationUserId: orgUserID,
		UserType:           userType,
	}, jc.AccessSecret, 3600)
	require.NoError(t, err)
	return metadata.NewIncomingContext(context.Background(), metadata.Pairs(utils.JwtTokenCtxKey, tok))
}

func makeAdmin(t *testing.T, db *gorm.DB) {
	t.Helper()
	require.NoError(t, db.Exec(`INSERT INTO organization_user_and_roles (organization_user_id, role_id) VALUES (50, 10)`).Error)
}

func grantStaffPolicy(fake *fakeEnforcer) {
	fake.seed([][]string{BuildPolicy(11, 100, constants.UserTypeLender)})
}

func authzReq() dto.AuthCheckRequestDto {
	return dto.AuthCheckRequestDto{Method: authzMethod, Path: authzPath, Type: constants.UserTypeLender}
}

// =============================================================================
// fake enforcer / test-seam unit tests
// =============================================================================

// Test_newAuthWithEnforcer_seam verifies the test seam compiles + a fakeEnforcer
// satisfies the policyEnforcer interface and behaves correctly.
func Test_newAuthWithEnforcer_seam(t *testing.T) {
	fake := newFakeEnforcer()
	a := newAuthWithEnforcer(nil, nil, nil, fake, "")
	require.NotNil(t, a)

	fake.seed([][]string{
		{"Role:1", "Lender", "Resource:5", "allow"},
	})

	ok, err := a.enforcer.Enforce("Role:1", "Lender", "Resource:5", "allow")
	require.NoError(t, err)
	assert.True(t, ok)

	ok, err = a.enforcer.Enforce("Role:2", "Lender", "Resource:5", "allow")
	require.NoError(t, err)
	assert.False(t, ok)
}

// Test_fakeEnforcer_addRemove verifies idempotency + AllPolicies correctness.
func Test_fakeEnforcer_addRemove(t *testing.T) {
	f := newFakeEnforcer()

	p1 := BuildPolicy(1, 10, constants.UserTypeLender)
	p2 := BuildPolicy(1, 20, constants.UserTypeLender)

	require.NoError(t, f.AddPolicies([][]string{p1, p2}))
	all, err := f.AllPolicies()
	require.NoError(t, err)
	assert.Len(t, all, 2)

	require.NoError(t, f.AddPolicies([][]string{p1}))
	all2, _ := f.AllPolicies()
	assert.Len(t, all2, 2)

	require.NoError(t, f.RemovePolicies([][]string{p1}))
	all3, _ := f.AllPolicies()
	assert.Len(t, all3, 1)

	require.NoError(t, f.RemovePolicies([][]string{p1}))
	all4, _ := f.AllPolicies()
	assert.Len(t, all4, 1)
}

// Test_fakeEnforcer_subjectIsolation verifies the subject-scoped filtering logic
// the rebuild diff tests exercise: policies whose subject is outside the scope's
// role-id set are left untouched.
func Test_fakeEnforcer_subjectIsolation(t *testing.T) {
	f := newFakeEnforcer()

	org10Role1 := BuildPolicy(100, 1, constants.UserTypeLender)
	org20Role2 := BuildPolicy(200, 2, constants.UserTypeLender)

	f.seed([][]string{org10Role1, org20Role2})

	all, err := f.AllPolicies()
	require.NoError(t, err)
	assert.Len(t, all, 2)

	org10Set := map[string]int64{"Role:100": 100}
	count := 0
	for _, p := range all {
		if _, ok := org10Set[p[0]]; ok {
			count++
		}
	}
	assert.Equal(t, 1, count, "only one policy should belong to Role:100")
}

// =============================================================================
// RebuildPolicy plan-gate tests
// =============================================================================

// --- baseline: active sub + plan-covered + bound role -> granted ---

func TestRebuild_ActiveSubGranted(t *testing.T) {
	a, fake, db := newAuthForTest(t)
	f := seedGranted(t, db)

	require.NoError(t, a.RebuildPolicyForOrganization(context.Background(), f.orgID, f.userType))
	assert.True(t, enforced(t, fake, f), "active sub + plan-covered resource must be granted")
}

// --- deny-default: no active subscription -> no policy (the original bug) ---

func TestRebuild_NoSubscription_DenyDefault(t *testing.T) {
	a, fake, db := newAuthForTest(t)
	f := seedGranted(t, db)
	require.NoError(t, db.Exec(`DELETE FROM organization_subscriptions`).Error)

	require.NoError(t, a.RebuildPolicyForOrganization(context.Background(), f.orgID, f.userType))
	assert.False(t, enforced(t, fake, f), "org without active subscription must be denied (deny-default)")
}

// --- plan does not cover the resource -> denied ---

func TestRebuild_ResourceNotCoveredByPlan_Denied(t *testing.T) {
	a, fake, db := newAuthForTest(t)
	f := seedGranted(t, db)
	require.NoError(t, db.Exec(`DELETE FROM pricing_plan_and_resources WHERE resource_id = 100`).Error)

	require.NoError(t, a.RebuildPolicyForOrganization(context.Background(), f.orgID, f.userType))
	assert.False(t, enforced(t, fake, f), "resource outside plan scope must be denied")
}

// --- disabled role -> no policy (AllResourcesForAuth preloads only active roles) ---

func TestRebuild_DisabledRole_Denied(t *testing.T) {
	a, fake, db := newAuthForTest(t)
	f := seedGranted(t, db)
	require.NoError(t, db.Exec(`UPDATE roles SET status = 'Disabled' WHERE id = 10`).Error)

	require.NoError(t, a.RebuildPolicyForOrganization(context.Background(), f.orgID, f.userType))
	assert.False(t, enforced(t, fake, f), "disabled role must produce no policy")
}

// --- disabled resource -> no policy (AllResourcesForAuth filters status=Active) ---

func TestRebuild_DisabledResource_Denied(t *testing.T) {
	a, fake, db := newAuthForTest(t)
	f := seedGranted(t, db)
	require.NoError(t, db.Exec(`UPDATE resources SET status = 'Disabled' WHERE id = 100`).Error)

	require.NoError(t, a.RebuildPolicyForOrganization(context.Background(), f.orgID, f.userType))
	assert.False(t, enforced(t, fake, f), "disabled resource must produce no policy")
}

// --- expired subscription -> treated as no active scope -> denied ---

func TestRebuild_ExpiredSubscription_Denied(t *testing.T) {
	a, fake, db := newAuthForTest(t)
	f := seedGranted(t, db)
	require.NoError(t, db.Exec(`UPDATE organization_subscriptions SET status = 'Expired'`).Error)

	require.NoError(t, a.RebuildPolicyForOrganization(context.Background(), f.orgID, f.userType))
	assert.False(t, enforced(t, fake, f), "expired subscription grants nothing")
}

// --- disabled plan -> plan resources excluded -> denied ---

func TestRebuild_DisabledPlan_Denied(t *testing.T) {
	a, fake, db := newAuthForTest(t)
	f := seedGranted(t, db)
	require.NoError(t, db.Exec(`UPDATE pricing_plans SET status = 'Disabled' WHERE id = 7`).Error)

	require.NoError(t, a.RebuildPolicyForOrganization(context.Background(), f.orgID, f.userType))
	assert.False(t, enforced(t, fake, f), "disabled plan must exclude its resources from the gate")
}

// --- becoming effective re-materializes existing bindings ---

func TestRebuild_BecomingEffective_GrantsAfterRebuild(t *testing.T) {
	a, fake, db := newAuthForTest(t)
	f := seedGranted(t, db)

	require.NoError(t, db.Exec(`DELETE FROM organization_subscriptions`).Error)
	require.NoError(t, a.RebuildPolicyForOrganization(context.Background(), f.orgID, f.userType))
	require.False(t, enforced(t, fake, f))

	require.NoError(t, db.Exec(`INSERT INTO organization_subscriptions (code, organization_id, pricing_plan_id, status, expire_date) VALUES ('SU2', 1, 7, 'Active', ?)`, time.Now().AddDate(0, 0, 10)).Error)
	require.NoError(t, a.RebuildPolicyForOrganization(context.Background(), f.orgID, f.userType))
	assert.True(t, enforced(t, fake, f), "becoming effective must re-materialize existing bindings")
}

// --- becoming ineffective: orphan cleanup removes the policy ---

func TestRebuild_BecomingIneffective_RevokesPolicy(t *testing.T) {
	a, fake, db := newAuthForTest(t)
	f := seedGranted(t, db)

	require.NoError(t, a.RebuildPolicyForOrganization(context.Background(), f.orgID, f.userType))
	require.True(t, enforced(t, fake, f))

	require.NoError(t, db.Exec(`UPDATE organization_subscriptions SET status = 'Suspended'`).Error)
	require.NoError(t, a.RebuildPolicyForOrganization(context.Background(), f.orgID, f.userType))
	assert.False(t, enforced(t, fake, f), "becoming ineffective must revoke the policy (orphan cleanup)")
}

// --- org-scoped isolation: rebuilding org 1 never touches org 2's policies ---

func TestRebuild_OrgScopedIsolation(t *testing.T) {
	a, fake, db := newAuthForTest(t)
	const ut = constants.UserTypeLender

	seedGranted(t, db)

	require.NoError(t, db.Exec(`INSERT INTO organizations (id, name, status, type) VALUES (2, 'Org2', 'Active', ?)`, ut).Error)
	require.NoError(t, db.Exec(`INSERT INTO resources (id, type, status, name) VALUES (200, ?, 'Active', 'Res200')`, ut).Error)
	require.NoError(t, db.Exec(`INSERT INTO roles (id, organization_id, name, type, status) VALUES (20, 2, 'Org2Role', ?, 'Active')`, ut).Error)
	require.NoError(t, db.Exec(`INSERT INTO role_and_resources (role_id, resource_id) VALUES (20, 200)`).Error)
	require.NoError(t, db.Exec(`INSERT INTO pricing_plan_and_resources (pricing_plan_id, resource_id) VALUES (7, 200)`).Error)
	require.NoError(t, db.Exec(`INSERT INTO organization_subscriptions (code, organization_id, pricing_plan_id, status, expire_date) VALUES ('SU_O2', 2, 7, 'Active', ?)`, time.Now().AddDate(0, 0, 10)).Error)

	org2Policy := BuildPolicy(20, 200, ut)
	fake.seed([][]string{org2Policy})

	require.NoError(t, a.RebuildPolicyForOrganization(context.Background(), 1, ut))

	ok, err := fake.Enforce(BuildInterfacePolicy(20, 200, ut)...)
	require.NoError(t, err)
	assert.True(t, ok, "org-scoped rebuild must not remove another org's policies")
}

// --- subscription expiry-by-date (not status): past expire_date is ineffective ---

func TestRebuild_PastExpireDate_Denied(t *testing.T) {
	a, fake, db := newAuthForTest(t)
	f := seedGranted(t, db)
	require.NoError(t, db.Exec(`UPDATE organization_subscriptions SET expire_date = ?`, time.Now().AddDate(0, 0, -1)).Error)

	require.NoError(t, a.RebuildPolicyForOrganization(context.Background(), f.orgID, f.userType))
	assert.False(t, enforced(t, fake, f), "subscription past its expire_date must not grant")
}

// --- gate cache invalidation: a second rebuild after a plan change reflects new state ---

func TestRebuild_PlanCacheInvalidatedOnRebuild(t *testing.T) {
	a, fake, db := newAuthForTest(t)
	f := seedGranted(t, db)

	require.NoError(t, a.RebuildPolicyForOrganization(context.Background(), f.orgID, f.userType))
	require.True(t, enforced(t, fake, f))

	require.NoError(t, db.Exec(`INSERT INTO resources (id, type, status, name) VALUES (101, ?, 'Active', 'Res101')`, f.userType).Error)
	require.NoError(t, db.Exec(`INSERT INTO role_and_resources (role_id, resource_id) VALUES (10, 101)`).Error)
	require.NoError(t, db.Exec(`INSERT INTO pricing_plan_and_resources (pricing_plan_id, resource_id) VALUES (7, 101)`).Error)

	require.NoError(t, a.RebuildPolicyForOrganization(context.Background(), f.orgID, f.userType))

	ok, err := fake.Enforce(BuildInterfacePolicy(f.roleID, 101, f.userType)...)
	require.NoError(t, err)
	assert.True(t, ok, "rebuild must recompute the plan gate from fresh DB state")
	assert.True(t, enforced(t, fake, f), "unchanged bindings must keep their policies")
}

// --- resource removal: unbinding a resource from a role drops its policy ---

func TestRebuild_RoleLosesResource_RevokesPolicy(t *testing.T) {
	a, fake, db := newAuthForTest(t)
	f := seedGranted(t, db)

	require.NoError(t, a.RebuildPolicyForOrganization(context.Background(), f.orgID, f.userType))
	require.True(t, enforced(t, fake, f))

	require.NoError(t, db.Exec(`DELETE FROM role_and_resources WHERE resource_id = 100`).Error)
	require.NoError(t, a.RebuildPolicyForOrganization(context.Background(), f.orgID, f.userType))
	assert.False(t, enforced(t, fake, f), "role losing the resource must revoke its policy")
}

// =============================================================================
// CheckAuthorization plan-gate tests
// =============================================================================

// --- admin: API within plan -> authorized ---

func TestCheckAuthorization_AdminInPlan_Authorized(t *testing.T) {
	a, _, db, _ := newAuthzForTest(t)
	seedAuthzBaseline(t, db)
	makeAdmin(t, db)

	ok, isRequired, err := a.CheckAuthorization(ctxWithClaim(t, a, constants.UserTypeLender, 1, 50), authzReq())
	require.NoError(t, err)
	assert.True(t, isRequired)
	assert.True(t, ok, "admin must be authorized when the API's resource is in the plan")
}

// --- admin: API outside plan -> denied (the core change) ---

func TestCheckAuthorization_AdminOutOfPlan_Denied(t *testing.T) {
	a, _, db, _ := newAuthzForTest(t)
	seedAuthzBaseline(t, db)
	makeAdmin(t, db)
	require.NoError(t, db.Exec(`DELETE FROM pricing_plan_and_resources WHERE resource_id = 100`).Error)

	ok, _, err := a.CheckAuthorization(ctxWithClaim(t, a, constants.UserTypeLender, 1, 50), authzReq())
	require.NoError(t, err)
	assert.False(t, ok, "admin must be denied when the API's resource is outside the plan")
}

// --- admin: org without any subscription -> denied (deny-default) ---

func TestCheckAuthorization_AdminNoSubscription_Denied(t *testing.T) {
	a, _, db, _ := newAuthzForTest(t)
	seedAuthzBaseline(t, db)
	makeAdmin(t, db)
	require.NoError(t, db.Exec(`DELETE FROM organization_subscriptions`).Error)

	ok, _, err := a.CheckAuthorization(ctxWithClaim(t, a, constants.UserTypeLender, 1, 50), authzReq())
	require.NoError(t, err)
	assert.False(t, ok, "admin with no active subscription must be denied")
}

// --- non-admin: in-plan + matching role policy -> authorized ---

func TestCheckAuthorization_NonAdminInPlanWithPolicy_Authorized(t *testing.T) {
	a, fake, db, _ := newAuthzForTest(t)
	seedAuthzBaseline(t, db)
	grantStaffPolicy(fake)

	ok, _, err := a.CheckAuthorization(ctxWithClaim(t, a, constants.UserTypeLender, 1, 50), authzReq())
	require.NoError(t, err)
	assert.True(t, ok, "non-admin with a plan-covered resource and matching policy must be authorized")
}

// --- non-admin: in-plan but no matching role policy -> denied ---

func TestCheckAuthorization_NonAdminNoPolicy_Denied(t *testing.T) {
	a, _, db, _ := newAuthzForTest(t)
	seedAuthzBaseline(t, db)

	ok, _, err := a.CheckAuthorization(ctxWithClaim(t, a, constants.UserTypeLender, 1, 50), authzReq())
	require.NoError(t, err)
	assert.False(t, ok, "non-admin without a matching role-resource policy must be denied")
}

// --- non-admin: resource outside plan -> denied even with a matching policy ---

func TestCheckAuthorization_NonAdminOutOfPlan_Denied(t *testing.T) {
	a, fake, db, _ := newAuthzForTest(t)
	seedAuthzBaseline(t, db)
	grantStaffPolicy(fake)
	require.NoError(t, db.Exec(`DELETE FROM pricing_plan_and_resources WHERE resource_id = 100`).Error)

	ok, _, err := a.CheckAuthorization(ctxWithClaim(t, a, constants.UserTypeLender, 1, 50), authzReq())
	require.NoError(t, err)
	assert.False(t, ok, "non-admin must be denied when the resource is outside the plan, even with a policy")
}

// --- API with no resource binding -> denied ---

func TestCheckAuthorization_NoResourceBinding_Denied(t *testing.T) {
	a, _, db, _ := newAuthzForTest(t)
	seedAuthzBaseline(t, db)
	makeAdmin(t, db)
	require.NoError(t, db.Exec(`DELETE FROM resource_and_apis WHERE resource_id = 100`).Error)

	ok, _, err := a.CheckAuthorization(ctxWithClaim(t, a, constants.UserTypeLender, 1, 50), authzReq())
	require.NoError(t, err)
	assert.False(t, ok, "API without any resource binding must be denied")
}

// --- authorization whitelist (permission_required=false) -> not required ---

func TestCheckAuthorization_Whitelist_NotRequired(t *testing.T) {
	a, _, db, _ := newAuthzForTest(t)
	seedAuthzBaseline(t, db)
	require.NoError(t, db.Exec(`UPDATE apis SET permission_required = 0 WHERE id = 1000`).Error)

	ok, isRequired, err := a.CheckAuthorization(ctxWithClaim(t, a, constants.UserTypeLender, 1, 50), authzReq())
	require.NoError(t, err)
	assert.False(t, isRequired, "whitelisted API must not require authorization")
	assert.False(t, ok, "authorized flag stays false when authorization is not required")
}
