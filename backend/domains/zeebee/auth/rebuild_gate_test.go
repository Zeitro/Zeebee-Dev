package auth

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/repo"
)

// setupAuthTestDB mirrors the logic-package sqlite fixture for the subset of
// tables the rebuild gate reads. Each test gets an isolated in-memory DB.
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

// newAuthForTest wires an Auth backed by the sqlite repos and a fake in-memory
// enforcer. redisCli is nil: only the org-scoped rebuild path (which never
// self-publishes) is exercised by these tests.
func newAuthForTest(t *testing.T) (*Auth, *fakeEnforcer, *gorm.DB) {
	t.Helper()
	db := setupAuthTestDB(t)
	repos := repo.NewAllRepos(db)
	fake := newFakeEnforcer()
	a := newAuthWithEnforcer(repos, nil, dto.JwtConfig{}, fake, "")
	return a, fake, db
}

type authFixture struct {
	orgID    int64
	roleID   int64
	resID    int64
	planID   int64
	userType constants.UserType
}

// seedGranted seeds org 1, an active lender role, a lender resource, the
// resource bound to the role, an active plan whose plan covers the resource,
// and an Active subscription. This is the baseline under which a rebuild
// MUST grant a policy (roleId, resourceId).
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
	// remove the subscription: org has a bound role but no effective plan scope
	require.NoError(t, db.Exec(`DELETE FROM organization_subscriptions`).Error)

	require.NoError(t, a.RebuildPolicyForOrganization(context.Background(), f.orgID, f.userType))
	assert.False(t, enforced(t, fake, f), "org without active subscription must be denied (deny-default)")
}

// --- plan does not cover the resource -> denied ---

func TestRebuild_ResourceNotCoveredByPlan_Denied(t *testing.T) {
	a, fake, db := newAuthForTest(t)
	f := seedGranted(t, db)
	// unbind the resource from the plan: role still bound, but plan no longer covers it
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
	// status Expired is not in EffectiveStatuses
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

// --- the fix: becoming effective re-materializes existing bindings ---

func TestRebuild_BecomingEffective_GrantsAfterRebuild(t *testing.T) {
	a, fake, db := newAuthForTest(t)
	f := seedGranted(t, db)

	// first rebuild with NO subscription: denied (deny-default)
	require.NoError(t, db.Exec(`DELETE FROM organization_subscriptions`).Error)
	require.NoError(t, a.RebuildPolicyForOrganization(context.Background(), f.orgID, f.userType))
	require.False(t, enforced(t, fake, f))

	// subscription created (becoming effective): rebuild must now grant.
	// Existing role-resource bindings are re-materialized through the gate.
	require.NoError(t, db.Exec(`INSERT INTO organization_subscriptions (code, organization_id, pricing_plan_id, status, expire_date) VALUES ('SU2', 1, 7, 'Active', ?)`, time.Now().AddDate(0, 0, 10)).Error)
	require.NoError(t, a.RebuildPolicyForOrganization(context.Background(), f.orgID, f.userType))
	assert.True(t, enforced(t, fake, f), "becoming effective must re-materialize existing bindings")
}

// --- becoming ineffective: orphan cleanup removes the policy ---

func TestRebuild_BecomingIneffective_RevokesPolicy(t *testing.T) {
	a, fake, db := newAuthForTest(t)
	f := seedGranted(t, db)

	// granted baseline
	require.NoError(t, a.RebuildPolicyForOrganization(context.Background(), f.orgID, f.userType))
	require.True(t, enforced(t, fake, f))

	// subscription suspended (becoming ineffective): rebuild must revoke.
	require.NoError(t, db.Exec(`UPDATE organization_subscriptions SET status = 'Suspended'`).Error)
	require.NoError(t, a.RebuildPolicyForOrganization(context.Background(), f.orgID, f.userType))
	assert.False(t, enforced(t, fake, f), "becoming ineffective must revoke the policy (orphan cleanup)")
}

// --- org-scoped isolation: rebuilding org 1 never touches org 2's policies ---

func TestRebuild_OrgScopedIsolation(t *testing.T) {
	a, fake, db := newAuthForTest(t)
	const ut = constants.UserTypeLender

	// org 1: granted baseline
	seedGranted(t, db)

	// org 2: separate role/resource/plan/subscription
	require.NoError(t, db.Exec(`INSERT INTO organizations (id, name, status, type) VALUES (2, 'Org2', 'Active', ?)`, ut).Error)
	require.NoError(t, db.Exec(`INSERT INTO resources (id, type, status, name) VALUES (200, ?, 'Active', 'Res200')`, ut).Error)
	require.NoError(t, db.Exec(`INSERT INTO roles (id, organization_id, name, type, status) VALUES (20, 2, 'Org2Role', ?, 'Active')`, ut).Error)
	require.NoError(t, db.Exec(`INSERT INTO role_and_resources (role_id, resource_id) VALUES (20, 200)`).Error)
	require.NoError(t, db.Exec(`INSERT INTO pricing_plan_and_resources (pricing_plan_id, resource_id) VALUES (7, 200)`).Error)
	require.NoError(t, db.Exec(`INSERT INTO organization_subscriptions (code, organization_id, pricing_plan_id, status, expire_date) VALUES ('SU_O2', 2, 7, 'Active', ?)`, time.Now().AddDate(0, 0, 10)).Error)

	// seed org 2's policy into the enforcer directly, then rebuild ONLY org 1.
	org2Policy := BuildPolicy(20, 200, ut)
	fake.seed([][]string{org2Policy})

	require.NoError(t, a.RebuildPolicyForOrganization(context.Background(), 1, ut))

	// org 2's policy survives: subject Role:20 is outside org 1's subject set.
	ok, err := fake.Enforce(BuildInterfacePolicy(20, 200, ut)...)
	require.NoError(t, err)
	assert.True(t, ok, "org-scoped rebuild must not remove another org's policies")
}

// --- subscription expiry-by-date (not status): past expire_date is ineffective ---

func TestRebuild_PastExpireDate_Denied(t *testing.T) {
	a, fake, db := newAuthForTest(t)
	f := seedGranted(t, db)
	// Active status but expire_date in the past -> isSubscriptionExpired true
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

	// add a second resource, bind it to role + plan, rebuild -> new policy appears
	require.NoError(t, db.Exec(`INSERT INTO resources (id, type, status, name) VALUES (101, ?, 'Active', 'Res101')`, f.userType).Error)
	require.NoError(t, db.Exec(`INSERT INTO role_and_resources (role_id, resource_id) VALUES (10, 101)`).Error)
	require.NoError(t, db.Exec(`INSERT INTO pricing_plan_and_resources (pricing_plan_id, resource_id) VALUES (7, 101)`).Error)

	require.NoError(t, a.RebuildPolicyForOrganization(context.Background(), f.orgID, f.userType))

	ok, err := fake.Enforce(BuildInterfacePolicy(f.roleID, 101, f.userType)...)
	require.NoError(t, err)
	assert.True(t, ok, "rebuild must recompute the plan gate from fresh DB state")
	// original resource still granted
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
