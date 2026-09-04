package logic

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
)

func TestSyncResource_AddsOnlyMissingApis(t *testing.T) {
	ppLogic, db := setupPricingPlanLogic(t)
	require.NoError(t, db.Exec(`CREATE TABLE apis (
		id INTEGER PRIMARY KEY, created_at DATETIME, updated_at DATETIME,
		is_del INTEGER NOT NULL DEFAULT 0, parent_id INTEGER NOT NULL DEFAULT 0,
		type TEXT NOT NULL, name TEXT NOT NULL DEFAULT '', path TEXT NOT NULL,
		method TEXT NOT NULL, description TEXT NOT NULL DEFAULT '',
		status TEXT NOT NULL DEFAULT 'Active', sort_order INTEGER NOT NULL DEFAULT 0,
		auth_required INTEGER NOT NULL DEFAULT 0, permission_required INTEGER NOT NULL DEFAULT 0
	)`).Error)
	require.NoError(t, db.Exec(`CREATE TABLE resource_and_apis (
		resource_id INTEGER NOT NULL, api_id INTEGER NOT NULL,
		PRIMARY KEY (resource_id, api_id)
	)`).Error)
	require.NoError(t, db.Exec(`INSERT INTO resources (id, type, resource_type, name, tag, status) VALUES (1, 'Borrower', 'Menu', 'Existing', 'Sync:Existing', 'Active')`).Error)
	require.NoError(t, db.Exec(`INSERT INTO apis (id, type, path, method) VALUES (1, 'Borrower', '/existing', 'GET'), (2, 'Borrower', '/new', 'POST')`).Error)
	require.NoError(t, db.Exec(`INSERT INTO resource_and_apis (resource_id, api_id) VALUES (1, 1)`).Error)

	resourceLogic := NewResourceLogic(ppLogic.BaseLogic, map[constants.UserType]string{})
	err := db.Transaction(func(tx *gorm.DB) error {
		return resourceLogic.handleSyncResource(context.Background(), tx, constants.UserTypeBorrower, dto.ResourceSyncItemDto{
			ResourceType: constants.ResourceTypeMenu,
			Tag:          "Sync:Existing",
			Apis: []dto.ApiSyncItemDto{
				{Path: "/existing", Method: "GET"},
				{Path: "/new", Method: "POST"},
			},
		}, 0)
	})
	require.NoError(t, err)

	var apiIDs []int64
	require.NoError(t, db.Raw(`SELECT api_id FROM resource_and_apis WHERE resource_id = 1 ORDER BY api_id`).Scan(&apiIDs).Error)
	assert.Equal(t, []int64{1, 2}, apiIDs)
}

// --- T056: Pricing plan resource intersection filter tests ---

func TestFilterByPricingPlan_NoSubscription_ReturnsEmpty(t *testing.T) {
	logic, _ := setupPricingPlanLogic(t)
	resourceLogic := NewResourceLogic(logic.BaseLogic, map[constants.UserType]string{})
	ctx := context.Background()

	resources := []*model.Resource{
		{Base: model.Base{Id: 1}},
		{Base: model.Base{Id: 2}},
	}

	result, err := resourceLogic.filterByPricingPlan(ctx, 1, resources)
	require.NoError(t, err)
	assert.Empty(t, result)
}

func TestFilterByPricingPlan_WithSubscription_ReturnsIntersection(t *testing.T) {
	ppLogic, db := setupPricingPlanLogic(t)
	ctx := context.Background()

	// Seed resources
	seedResources(t, db, []int64{1, 2, 3, 4, 5})

	// Seed organization
	require.NoError(t, db.Exec("INSERT INTO organizations (id, name, status) VALUES (1, 'Test Org', 'Active')").Error)

	// Seed pricing plan
	plan := &model.PricingPlan{
		Name:     "Test Plan",
		Type:     constants.UserTypeLender,
		PlanType: constants.PricingPlanTypeStandard,
		Status:   constants.CommonStatusActive,
	}
	require.NoError(t, db.Create(plan).Error)

	// Bind resources 1, 2, 3 to plan
	err := ppLogic.BindResources(ctx, plan.Id, []int64{1, 2, 3})
	require.NoError(t, err)

	// Seed subscription
	policy := &model.PricingPlanPolicy{
		PricingPlanId:  plan.Id,
		OriginalAmount: 100,
		ActualAmount:   80,
		Type:           constants.PricingPolicyTypeMonthly,
	}
	require.NoError(t, db.Create(policy).Error)

	// Seed subscription via direct SQL (reuse same db)
	require.NoError(t, db.Exec("INSERT INTO organization_subscriptions (code, organization_id, pricing_plan_id, pricing_plan_policy_id, status, original_amount, actual_amount) VALUES (?, 1, ?, ?, 'Active', 100, 80)", "SU001", plan.Id, policy.Id).Error)

	// Test filter — all resources
	resources := []*model.Resource{
		{Base: model.Base{Id: 1}},
		{Base: model.Base{Id: 2}},
		{Base: model.Base{Id: 3}},
		{Base: model.Base{Id: 4}},
		{Base: model.Base{Id: 5}},
	}

	resourceLogic := NewResourceLogic(ppLogic.BaseLogic, map[constants.UserType]string{})
	result, err := resourceLogic.filterByPricingPlan(ctx, 1, resources)
	require.NoError(t, err)
	assert.Len(t, result, 3)

	resultIds := make(map[int64]bool)
	for _, r := range result {
		resultIds[r.Id] = true
	}
	assert.True(t, resultIds[1])
	assert.True(t, resultIds[2])
	assert.True(t, resultIds[3])
	assert.False(t, resultIds[4])
	assert.False(t, resultIds[5])
}

func TestFilterByPricingPlan_WithSubscription_NoOverlap_ReturnsEmpty(t *testing.T) {
	ppLogic, db := setupPricingPlanLogic(t)
	ctx := context.Background()

	// Seed resources
	seedResources(t, db, []int64{1, 2, 3, 10, 11})

	// Seed organization
	require.NoError(t, db.Exec("INSERT INTO organizations (id, name, status) VALUES (1, 'Test Org', 'Active')").Error)

	// Seed pricing plan with resources 10, 11
	plan := &model.PricingPlan{
		Name:     "Test Plan",
		Type:     constants.UserTypeLender,
		PlanType: constants.PricingPlanTypeStandard,
		Status:   constants.CommonStatusActive,
	}
	require.NoError(t, db.Create(plan).Error)
	err := ppLogic.BindResources(ctx, plan.Id, []int64{10, 11})
	require.NoError(t, err)

	// Seed subscription
	policy := &model.PricingPlanPolicy{
		PricingPlanId:  plan.Id,
		OriginalAmount: 100,
		ActualAmount:   80,
		Type:           constants.PricingPolicyTypeMonthly,
	}
	require.NoError(t, db.Create(policy).Error)

	// Seed subscription via direct SQL (reuse same db)
	require.NoError(t, db.Exec("INSERT INTO organization_subscriptions (code, organization_id, pricing_plan_id, pricing_plan_policy_id, status, original_amount, actual_amount) VALUES (?, 1, ?, ?, 'Active', 100, 80)", "SU002", plan.Id, policy.Id).Error)

	// Filter resources 1, 2, 3 — none overlap with plan's 10, 11
	resources := []*model.Resource{
		{Base: model.Base{Id: 1}},
		{Base: model.Base{Id: 2}},
		{Base: model.Base{Id: 3}},
	}

	resourceLogic := NewResourceLogic(ppLogic.BaseLogic, map[constants.UserType]string{})
	result, err := resourceLogic.filterByPricingPlan(ctx, 1, resources)
	require.NoError(t, err)
	assert.Empty(t, result)
}
