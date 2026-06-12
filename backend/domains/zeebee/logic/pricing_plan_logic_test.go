package logic

import (
	"context"
	"fmt"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/repo"
)

// setupTestDB creates an in-memory SQLite DB with manually created tables
func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	sqls := []string{
		`CREATE TABLE IF NOT EXISTS pricing_plans (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			is_del INTEGER DEFAULT 0,
			name TEXT NOT NULL DEFAULT '',
			subtitle TEXT NOT NULL DEFAULT '',
			description TEXT NOT NULL DEFAULT '',
			status TEXT NOT NULL DEFAULT 'Active',
			trial_days INTEGER NOT NULL DEFAULT 0,
			sort_order INTEGER NOT NULL DEFAULT 6,
			is_public INTEGER NOT NULL DEFAULT 0,
			is_default INTEGER NOT NULL DEFAULT 0,
			type TEXT NOT NULL DEFAULT '',
			plan_type TEXT NOT NULL DEFAULT 'Standard'
		)`,
		`CREATE TABLE IF NOT EXISTS pricing_plan_policies (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			is_del INTEGER DEFAULT 0,
			pricing_plan_id INTEGER NOT NULL DEFAULT 0,
			original_amount REAL NOT NULL DEFAULT 0,
			actual_amount REAL NOT NULL DEFAULT 0,
			type TEXT NOT NULL DEFAULT 'Monthly'
		)`,
		`CREATE TABLE IF NOT EXISTS pricing_plan_and_resources (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			pricing_plan_id INTEGER NOT NULL DEFAULT 0,
			resource_id INTEGER NOT NULL DEFAULT 0
		)`,
		`CREATE TABLE IF NOT EXISTS organization_subscriptions (
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
		`CREATE TABLE IF NOT EXISTS organizations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			is_del INTEGER DEFAULT 0,
			name TEXT NOT NULL DEFAULT '',
			status TEXT NOT NULL DEFAULT 'Active'
		)`,
		`CREATE TABLE IF NOT EXISTS resources (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			is_del INTEGER DEFAULT 0,
			parent_id INTEGER DEFAULT 0,
			type TEXT DEFAULT 'Borrower',
			resource_type TEXT DEFAULT 'Menu',
			name TEXT DEFAULT '',
			path TEXT DEFAULT '',
			icon TEXT DEFAULT '',
			tag TEXT DEFAULT '',
			is_show INTEGER DEFAULT 0,
			status TEXT DEFAULT 'Active',
			sort_order INTEGER DEFAULT 10,
			enable_data_scope INTEGER DEFAULT 0,
			description TEXT DEFAULT '',
			is_public INTEGER DEFAULT 0
		)`,
	}
	for _, s := range sqls {
		require.NoError(t, db.Exec(s).Error)
	}
	return db
}

// setupPricingPlanLogic creates a PricingPlanLogic with test DB
func setupPricingPlanLogic(t *testing.T) (*PricingPlanLogic, *gorm.DB) {
	t.Helper()
	db := setupTestDB(t)
	repos := repo.NewAllRepos(db)
	baseLogic := NewBaseLogic(repos, nil)
	authTopicMap := map[constants.UserType]string{
		constants.UserTypeBorrower: "test-borrower-topic",
		constants.UserTypeLender:   "test-lender-topic",
	}
	return NewPricingPlanLogic(baseLogic, authTopicMap), db
}

func seedPricingPlan(t *testing.T, db *gorm.DB, plan *model.PricingPlan) *model.PricingPlan {
	t.Helper()
	err := db.Create(plan).Error
	require.NoError(t, err)
	return plan
}

func seedResources(t *testing.T, db *gorm.DB, ids []int64) {
	t.Helper()
	for _, id := range ids {
		err := db.Exec("INSERT INTO resources (id, name, status, type) VALUES (?, ?, 'Active', ?)", id, fmt.Sprintf("Resource %d", id), constants.UserTypeLender).Error
		require.NoError(t, err)
	}
}

// --- T020: Create failure cases ---

func TestPricingPlanCreate_DuplicateNameUnderSameType(t *testing.T) {
	logic, _ := setupPricingPlanLogic(t)
	ctx := context.Background()

	input := dto.PricingPlanCreateDto{
		Name:     "Pro Plan",
		Type:     constants.UserTypeLender,
		PlanType: constants.PricingPlanTypeStandard,
	}

	_, err := logic.Create(ctx, input)
	require.NoError(t, err)

	_, err = logic.Create(ctx, input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
}

func TestPricingPlanCreate_DifferentTypeSameName_OK(t *testing.T) {
	logic, _ := setupPricingPlanLogic(t)
	ctx := context.Background()

	input1 := dto.PricingPlanCreateDto{
		Name:     "Pro Plan",
		Type:     constants.UserTypeLender,
		PlanType: constants.PricingPlanTypeStandard,
	}
	input2 := dto.PricingPlanCreateDto{
		Name:     "Pro Plan",
		Type:     constants.UserTypeBorrower,
		PlanType: constants.PricingPlanTypeStandard,
	}

	result1, err := logic.Create(ctx, input1)
	require.NoError(t, err)
	assert.NotNil(t, result1)

	result2, err := logic.Create(ctx, input2)
	require.NoError(t, err)
	assert.NotNil(t, result2)
	assert.NotEqual(t, result1.Id, result2.Id)
}

// --- T021: BindResources failure cases ---

func TestPricingPlanBindResources_PlanNotFound(t *testing.T) {
	logic, _ := setupPricingPlanLogic(t)
	ctx := context.Background()

	err := logic.BindResources(ctx, 99999, []int64{1, 2, 3})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestPricingPlanBindResources_ReplaceExisting(t *testing.T) {
	logic, db := setupPricingPlanLogic(t)
	ctx := context.Background()

	plan := seedPricingPlan(t, db, &model.PricingPlan{
		Name:     "Test Plan",
		Type:     constants.UserTypeLender,
		PlanType: constants.PricingPlanTypeStandard,
		Status:   constants.CommonStatusActive,
	})

	seedResources(t, db, []int64{1, 2, 3, 4, 5})

	err := logic.BindResources(ctx, plan.Id, []int64{1, 2, 3})
	require.NoError(t, err)

	err = logic.BindResources(ctx, plan.Id, []int64{4, 5})
	require.NoError(t, err)

	var count int64
	db.Model(&model.PricingPlanAndResource{}).Where("pricing_plan_id = ?", plan.Id).Count(&count)
	assert.Equal(t, int64(2), count)
}

func TestPricingPlanBindResources_EmptyList(t *testing.T) {
	logic, db := setupPricingPlanLogic(t)
	ctx := context.Background()

	plan := seedPricingPlan(t, db, &model.PricingPlan{
		Name:     "Test Plan",
		Type:     constants.UserTypeLender,
		PlanType: constants.PricingPlanTypeStandard,
		Status:   constants.CommonStatusActive,
	})

	seedResources(t, db, []int64{1})

	err := logic.BindResources(ctx, plan.Id, []int64{1})
	require.NoError(t, err)

	err = logic.BindResources(ctx, plan.Id, []int64{})
	require.NoError(t, err)

	var count int64
	db.Model(&model.PricingPlanAndResource{}).Where("pricing_plan_id = ?", plan.Id).Count(&count)
	assert.Equal(t, int64(0), count)
}

// --- T022: AddPolicy failure cases ---

func TestPricingPlanPolicyCreate_DuplicateType(t *testing.T) {
	logic, _ := setupPricingPlanLogic(t)
	ctx := context.Background()

	planInput := dto.PricingPlanCreateDto{
		Name:     "Test Plan",
		Type:     constants.UserTypeLender,
		PlanType: constants.PricingPlanTypeStandard,
	}
	plan, err := logic.Create(ctx, planInput)
	require.NoError(t, err)

	policyLogic := NewPricingPlanPolicyLogic(logic.BaseLogic)

	input1 := dto.PricingPlanPolicyCreateDto{
		PricingPlanId:  plan.Id,
		OriginalAmount: 100,
		ActualAmount:   80,
		Type:           constants.PricingPolicyTypeMonthly,
	}
	_, err = policyLogic.Create(ctx, input1)
	require.NoError(t, err)

	input2 := dto.PricingPlanPolicyCreateDto{
		PricingPlanId:  plan.Id,
		OriginalAmount: 200,
		ActualAmount:   160,
		Type:           constants.PricingPolicyTypeMonthly,
	}
	_, err = policyLogic.Create(ctx, input2)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
}

func TestPricingPlanPolicyCreate_DifferentType_OK(t *testing.T) {
	logic, _ := setupPricingPlanLogic(t)
	ctx := context.Background()

	planInput := dto.PricingPlanCreateDto{
		Name:     "Test Plan",
		Type:     constants.UserTypeLender,
		PlanType: constants.PricingPlanTypeStandard,
	}
	plan, err := logic.Create(ctx, planInput)
	require.NoError(t, err)

	policyLogic := NewPricingPlanPolicyLogic(logic.BaseLogic)

	input1 := dto.PricingPlanPolicyCreateDto{
		PricingPlanId:  plan.Id,
		OriginalAmount: 100,
		ActualAmount:   80,
		Type:           constants.PricingPolicyTypeMonthly,
	}
	_, err = policyLogic.Create(ctx, input1)
	require.NoError(t, err)

	input2 := dto.PricingPlanPolicyCreateDto{
		PricingPlanId:  plan.Id,
		OriginalAmount: 300,
		ActualAmount:   240,
		Type:           constants.PricingPolicyTypeYearly,
	}
	result, err := policyLogic.Create(ctx, input2)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

// --- US2 tests (T031-T033) ---

// T031: Update failure - Name not changeable
func TestPricingPlanUpdate_NameNotChangeable(t *testing.T) {
	logic, _ := setupPricingPlanLogic(t)
	ctx := context.Background()

	input := dto.PricingPlanCreateDto{
		Name:     "Original Name",
		Type:     constants.UserTypeLender,
		PlanType: constants.PricingPlanTypeStandard,
	}
	created, err := logic.Create(ctx, input)
	require.NoError(t, err)

	updateInput := dto.PricingPlanUpdateDto{
		Id:          created.Id,
		Subtitle:    "Updated Subtitle",
		Description: "Updated Description",
	}

	updated, err := logic.Update(ctx, updateInput)
	require.NoError(t, err)
	assert.Equal(t, "Original Name", updated.Name)
	assert.Equal(t, "Updated Subtitle", updated.Subtitle)
}

// T032: Delete failure - active subscription exists
func TestPricingPlanDelete_ActiveSubscription(t *testing.T) {
	logic, db := setupPricingPlanLogic(t)
	ctx := context.Background()

	plan := seedPricingPlan(t, db, &model.PricingPlan{
		Name:     "Test Plan",
		Type:     constants.UserTypeLender,
		PlanType: constants.PricingPlanTypeStandard,
		Status:   constants.CommonStatusActive,
	})

	err := db.Create(&model.OrganizationSubscription{
		Code:                "SU001",
		OrganizationId:      1,
		PricingPlanId:       plan.Id,
		PricingPlanPolicyId: 1,
		Status:              constants.OrganizationSubscriptionStatusActive,
		OriginalAmount:      100,
		ActualAmount:        80,
	}).Error
	require.NoError(t, err)

	err = logic.Delete(ctx, plan.Id)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "active subscriptions")
}

func TestPricingPlanDelete_ExpiredSubscription_OK(t *testing.T) {
	logic, db := setupPricingPlanLogic(t)
	ctx := context.Background()

	plan := seedPricingPlan(t, db, &model.PricingPlan{
		Name:     "Test Plan",
		Type:     constants.UserTypeLender,
		PlanType: constants.PricingPlanTypeStandard,
		Status:   constants.CommonStatusActive,
	})

	err := db.Create(&model.OrganizationSubscription{
		Code:                "SU001",
		OrganizationId:      1,
		PricingPlanId:       plan.Id,
		PricingPlanPolicyId: 1,
		Status:              constants.OrganizationSubscriptionStatusExpired,
		OriginalAmount:      100,
		ActualAmount:        80,
	}).Error
	require.NoError(t, err)

	err = logic.Delete(ctx, plan.Id)
	assert.NoError(t, err)
}

// T033: UpdatePolicy/DeletePolicy tests
func TestPricingPlanPolicyUpdate_NotFound(t *testing.T) {
	logic, _ := setupPricingPlanLogic(t)
	ctx := context.Background()
	policyLogic := NewPricingPlanPolicyLogic(logic.BaseLogic)

	input := dto.PricingPlanPolicyUpdateDto{
		Id:             99999,
		OriginalAmount: 200,
		ActualAmount:   150,
	}
	_, err := policyLogic.Update(ctx, 99999, input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestPricingPlanPolicyDelete_NotFound(t *testing.T) {
	logic, _ := setupPricingPlanLogic(t)
	ctx := context.Background()
	policyLogic := NewPricingPlanPolicyLogic(logic.BaseLogic)

	err := policyLogic.Delete(ctx, 99999)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestPricingPlanPolicyUpdate_Success(t *testing.T) {
	logic, _ := setupPricingPlanLogic(t)
	ctx := context.Background()

	planInput := dto.PricingPlanCreateDto{
		Name:     "Test Plan",
		Type:     constants.UserTypeLender,
		PlanType: constants.PricingPlanTypeStandard,
	}
	plan, err := logic.Create(ctx, planInput)
	require.NoError(t, err)

	policyLogic := NewPricingPlanPolicyLogic(logic.BaseLogic)

	createInput := dto.PricingPlanPolicyCreateDto{
		PricingPlanId:  plan.Id,
		OriginalAmount: 100,
		ActualAmount:   80,
		Type:           constants.PricingPolicyTypeMonthly,
	}
	policy, err := policyLogic.Create(ctx, createInput)
	require.NoError(t, err)

	updateInput := dto.PricingPlanPolicyUpdateDto{
		Id:             policy.Id,
		OriginalAmount: 200,
		ActualAmount:   150,
	}
	updated, err := policyLogic.Update(ctx, policy.Id, updateInput)
	require.NoError(t, err)
	assert.Equal(t, float64(200), updated.OriginalAmount)
	assert.Equal(t, float64(150), updated.ActualAmount)
}

func TestPricingPlanPolicyDelete_Success(t *testing.T) {
	logic, _ := setupPricingPlanLogic(t)
	ctx := context.Background()

	planInput := dto.PricingPlanCreateDto{
		Name:     "Test Plan",
		Type:     constants.UserTypeLender,
		PlanType: constants.PricingPlanTypeStandard,
	}
	plan, err := logic.Create(ctx, planInput)
	require.NoError(t, err)

	policyLogic := NewPricingPlanPolicyLogic(logic.BaseLogic)

	createInput := dto.PricingPlanPolicyCreateDto{
		PricingPlanId:  plan.Id,
		OriginalAmount: 100,
		ActualAmount:   80,
		Type:           constants.PricingPolicyTypeMonthly,
	}
	policy, err := policyLogic.Create(ctx, createInput)
	require.NoError(t, err)

	err = policyLogic.Delete(ctx, policy.Id)
	assert.NoError(t, err)

	_, err = policyLogic.Update(ctx, policy.Id, dto.PricingPlanPolicyUpdateDto{Id: policy.Id})
	assert.Error(t, err)
}

// --- Additional US1/US2 success path tests ---

func TestPricingPlanCreate_Success(t *testing.T) {
	logic, _ := setupPricingPlanLogic(t)
	ctx := context.Background()

	input := dto.PricingPlanCreateDto{
		Name:        "Pro Plan",
		Subtitle:    "For professionals",
		Description: "Full featured plan",
		Type:        constants.UserTypeLender,
		PlanType:    constants.PricingPlanTypeStandard,
		TrialDays:   14,
	}

	result, err := logic.Create(ctx, input)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Pro Plan", result.Name)
	assert.Equal(t, "For professionals", result.Subtitle)
	assert.Equal(t, int16(14), result.TrialDays)
	assert.NotZero(t, result.Id)
}

func TestPricingPlanDetail_Success(t *testing.T) {
	logic, _ := setupPricingPlanLogic(t)
	ctx := context.Background()

	input := dto.PricingPlanCreateDto{
		Name:     "Detail Plan",
		Type:     constants.UserTypeLender,
		PlanType: constants.PricingPlanTypeStandard,
	}
	created, err := logic.Create(ctx, input)
	require.NoError(t, err)

	detail, err := logic.Detail(ctx, created.Id)
	require.NoError(t, err)
	assert.Equal(t, created.Id, detail.Id)
	assert.Equal(t, "Detail Plan", detail.Name)
}

func TestPricingPlanDetail_NotFound(t *testing.T) {
	logic, _ := setupPricingPlanLogic(t)
	ctx := context.Background()

	_, err := logic.Detail(ctx, 99999)
	assert.Error(t, err)
}

func TestPricingPlanSwitchStatus_Success(t *testing.T) {
	logic, _ := setupPricingPlanLogic(t)
	ctx := context.Background()

	input := dto.PricingPlanCreateDto{
		Name:   "Status Plan",
		Type:   constants.UserTypeLender,
		Status: constants.CommonStatusActive,
	}
	created, err := logic.Create(ctx, input)
	require.NoError(t, err)

	err = logic.SwitchStatus(ctx, created.Id, constants.CommonStatusDisabled)
	assert.NoError(t, err)

	detail, err := logic.Detail(ctx, created.Id)
	require.NoError(t, err)
	assert.Equal(t, string(constants.CommonStatusDisabled), detail.Status)
}

func TestPricingPlanResourceList(t *testing.T) {
	logic, db := setupPricingPlanLogic(t)
	ctx := context.Background()

	plan := seedPricingPlan(t, db, &model.PricingPlan{
		Name:     "Resource Plan",
		Type:     constants.UserTypeLender,
		PlanType: constants.PricingPlanTypeStandard,
		Status:   constants.CommonStatusActive,
	})

	seedResources(t, db, []int64{10, 20, 30})

	err := logic.BindResources(ctx, plan.Id, []int64{10, 20, 30})
	require.NoError(t, err)

	list, err := logic.ResourceList(ctx, plan.Id)
	require.NoError(t, err)
	assert.Len(t, list, 3)
}
