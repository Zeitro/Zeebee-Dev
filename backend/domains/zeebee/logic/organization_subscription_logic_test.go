package logic

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/repo"
)

func setupSubscriptionLogic(t *testing.T) (*OrganizationSubscriptionLogic, *gorm.DB) {
	t.Helper()
	db := setupTestDB(t)

	repos := repo.NewAllRepos(db)
	baseLogic := NewBaseLogic(repos, nil)
	authTopicMap := map[constants.UserType]string{
		constants.UserTypeBorrower: "test-borrower-topic",
		constants.UserTypeLender:   "test-lender-topic",
	}
	return NewOrganizationSubscriptionLogic(baseLogic, authTopicMap), db
}

// --- T043: Create failure cases ---

func TestSubscriptionCreate_PlanDisabled(t *testing.T) {
	subLogic, db := setupSubscriptionLogic(t)
	ctx := context.Background()

	plan := &model.PricingPlan{
		Name:     "Disabled Plan",
		Type:     constants.UserTypeLender,
		PlanType: constants.PricingPlanTypeStandard,
		Status:   constants.CommonStatusDisabled,
	}
	require.NoError(t, db.Create(plan).Error)

	policy := &model.PricingPlanPolicy{
		PricingPlanId:  plan.Id,
		OriginalAmount: 100,
		ActualAmount:   80,
		Type:           constants.PricingPolicyTypeMonthly,
	}
	require.NoError(t, db.Create(policy).Error)

	input := dto.OrganizationSubscriptionCreateDto{
		OrganizationId:      1,
		PricingPlanId:       plan.Id,
		PricingPlanPolicyId: policy.Id,
	}

	_, err := subLogic.Create(ctx, input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "disabled")
}

func TestSubscriptionCreate_PolicyNotFound(t *testing.T) {
	subLogic, db := setupSubscriptionLogic(t)
	ctx := context.Background()

	plan := &model.PricingPlan{
		Name:     "Active Plan",
		Type:     constants.UserTypeLender,
		PlanType: constants.PricingPlanTypeStandard,
		Status:   constants.CommonStatusActive,
	}
	require.NoError(t, db.Create(plan).Error)

	input := dto.OrganizationSubscriptionCreateDto{
		OrganizationId:      1,
		PricingPlanId:       plan.Id,
		PricingPlanPolicyId: 99999,
	}

	_, err := subLogic.Create(ctx, input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "policy not found")
}

// --- T044: Create success cases ---

func TestSubscriptionCreate_CodeFormat(t *testing.T) {
	subLogic, db := setupSubscriptionLogic(t)
	ctx := context.Background()

	plan := &model.PricingPlan{
		Name:     "Active Plan",
		Type:     constants.UserTypeLender,
		PlanType: constants.PricingPlanTypeStandard,
		Status:   constants.CommonStatusActive,
	}
	require.NoError(t, db.Create(plan).Error)

	policy := &model.PricingPlanPolicy{
		PricingPlanId:  plan.Id,
		OriginalAmount: 100,
		ActualAmount:   80,
		Type:           constants.PricingPolicyTypeMonthly,
	}
	require.NoError(t, db.Create(policy).Error)

	input := dto.OrganizationSubscriptionCreateDto{
		OrganizationId:      1,
		PricingPlanId:       plan.Id,
		PricingPlanPolicyId: policy.Id,
	}

	result, err := subLogic.Create(ctx, input)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, len(result.Code) > 2, "Code should have prefix + digits")
	assert.Equal(t, "SU", result.Code[:2])
}

func TestSubscriptionCreate_ExpireDateCalculated(t *testing.T) {
	subLogic, db := setupSubscriptionLogic(t)
	ctx := context.Background()

	plan := &model.PricingPlan{
		Name:     "Active Plan",
		Type:     constants.UserTypeLender,
		PlanType: constants.PricingPlanTypeStandard,
		Status:   constants.CommonStatusActive,
	}
	require.NoError(t, db.Create(plan).Error)

	policy := &model.PricingPlanPolicy{
		PricingPlanId:  plan.Id,
		OriginalAmount: 100,
		ActualAmount:   80,
		Type:           constants.PricingPolicyTypeMonthly,
	}
	require.NoError(t, db.Create(policy).Error)

	input := dto.OrganizationSubscriptionCreateDto{
		OrganizationId:      1,
		PricingPlanId:       plan.Id,
		PricingPlanPolicyId: policy.Id,
	}

	beforeCreate := time.Now()
	result, err := subLogic.Create(ctx, input)
	require.NoError(t, err)

	assert.NotNil(t, result.ExpireDate)
	assert.True(t, result.ExpireDate.Valid)
	expectedExpiry := beforeCreate.AddDate(0, 0, 30)
	tolerance := 5 * time.Second
	assert.WithinDuration(t, expectedExpiry, result.ExpireDate.Time, tolerance)
}

func TestSubscriptionCreate_TrialStatus(t *testing.T) {
	subLogic, db := setupSubscriptionLogic(t)
	ctx := context.Background()

	plan := &model.PricingPlan{
		Name:      "Trial Plan",
		Type:      constants.UserTypeLender,
		PlanType:  constants.PricingPlanTypeStandard,
		Status:    constants.CommonStatusActive,
		TrialDays: 14,
		IsDefault: true,
	}
	require.NoError(t, db.Create(plan).Error)

	// SubscribeDefaultPlan creates a trial subscription
	err := subLogic.SubscribeDefaultPlan(ctx, 1, constants.UserTypeLender)
	require.NoError(t, err)

	var sub model.OrganizationSubscription
	require.NoError(t, db.First(&sub, "organization_id = ?", 1).Error)
	assert.Equal(t, constants.OrganizationSubscriptionStatusTrialing, sub.Status)
	assert.Equal(t, float64(0), sub.OriginalAmount)
	assert.Equal(t, float64(0), sub.ActualAmount)
	assert.NotNil(t, sub.ExpireDate)
}

func TestSubscriptionCreate_OldSubscriptionExpired(t *testing.T) {
	subLogic, db := setupSubscriptionLogic(t)
	ctx := context.Background()

	plan := &model.PricingPlan{
		Name:     "Active Plan",
		Type:     constants.UserTypeLender,
		PlanType: constants.PricingPlanTypeStandard,
		Status:   constants.CommonStatusActive,
	}
	require.NoError(t, db.Create(plan).Error)

	policy := &model.PricingPlanPolicy{
		PricingPlanId:  plan.Id,
		OriginalAmount: 100,
		ActualAmount:   80,
		Type:           constants.PricingPolicyTypeMonthly,
	}
	require.NoError(t, db.Create(policy).Error)

	oldSub := &model.OrganizationSubscription{
		Code:                "SU_OLD",
		OrganizationId:      1,
		PricingPlanId:       plan.Id,
		PricingPlanPolicyId: policy.Id,
		Status:              constants.OrganizationSubscriptionStatusActive,
		ExpireDate:          &sql.NullTime{Time: time.Now().AddDate(0, 0, 10), Valid: true},
		OriginalAmount:      100,
		ActualAmount:        80,
	}
	require.NoError(t, db.Create(oldSub).Error)

	input := dto.OrganizationSubscriptionCreateDto{
		OrganizationId:      1,
		PricingPlanId:       plan.Id,
		PricingPlanPolicyId: policy.Id,
	}

	result, err := subLogic.Create(ctx, input)
	require.NoError(t, err)
	assert.NotNil(t, result)

	var oldSubResult model.OrganizationSubscription
	require.NoError(t, db.First(&oldSubResult, oldSub.Id).Error)
	assert.Equal(t, string(constants.OrganizationSubscriptionStatusExpired), string(oldSubResult.Status))
}

func TestSubscriptionCreate_AmountCopied(t *testing.T) {
	subLogic, db := setupSubscriptionLogic(t)
	ctx := context.Background()

	plan := &model.PricingPlan{
		Name:     "Active Plan",
		Type:     constants.UserTypeLender,
		PlanType: constants.PricingPlanTypeStandard,
		Status:   constants.CommonStatusActive,
	}
	require.NoError(t, db.Create(plan).Error)

	policy := &model.PricingPlanPolicy{
		PricingPlanId:  plan.Id,
		OriginalAmount: 299.99,
		ActualAmount:   199.99,
		Type:           constants.PricingPolicyTypeYearly,
	}
	require.NoError(t, db.Create(policy).Error)

	input := dto.OrganizationSubscriptionCreateDto{
		OrganizationId:      1,
		PricingPlanId:       plan.Id,
		PricingPlanPolicyId: policy.Id,
	}

	result, err := subLogic.Create(ctx, input)
	require.NoError(t, err)
	assert.Equal(t, 299.99, result.OriginalAmount)
	assert.Equal(t, 199.99, result.ActualAmount)
}

// --- T045: Renew/UpdateStatus cases ---

func TestSubscriptionRenew_Success(t *testing.T) {
	subLogic, db := setupSubscriptionLogic(t)
	ctx := context.Background()

	plan := &model.PricingPlan{
		Name:     "Active Plan",
		Type:     constants.UserTypeLender,
		PlanType: constants.PricingPlanTypeStandard,
		Status:   constants.CommonStatusActive,
	}
	require.NoError(t, db.Create(plan).Error)

	policy := &model.PricingPlanPolicy{
		PricingPlanId:  plan.Id,
		OriginalAmount: 100,
		ActualAmount:   80,
		Type:           constants.PricingPolicyTypeMonthly,
	}
	require.NoError(t, db.Create(policy).Error)

	createInput := dto.OrganizationSubscriptionCreateDto{
		OrganizationId:      1,
		PricingPlanId:       plan.Id,
		PricingPlanPolicyId: policy.Id,
	}
	created, err := subLogic.Create(ctx, createInput)
	require.NoError(t, err)

	newExpiry := &sql.NullTime{Time: time.Now().AddDate(0, 1, 0), Valid: true}
	renewInput := dto.OrganizationSubscriptionRenewDto{
		Id:         created.Id,
		ExpireDate: newExpiry,
	}

	renewed, err := subLogic.Renew(ctx, created.Id, renewInput)
	require.NoError(t, err)
	assert.NotNil(t, renewed.ExpireDate)
	assert.True(t, renewed.ExpireDate.Valid)
}

func TestSubscriptionRenew_NotFound(t *testing.T) {
	subLogic, _ := setupSubscriptionLogic(t)
	ctx := context.Background()

	renewInput := dto.OrganizationSubscriptionRenewDto{
		Id:         99999,
		ExpireDate: &sql.NullTime{Time: time.Now().AddDate(0, 1, 0), Valid: true},
	}

	_, err := subLogic.Renew(ctx, 99999, renewInput)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestSubscriptionSwitchStatus_Success(t *testing.T) {
	subLogic, db := setupSubscriptionLogic(t)
	ctx := context.Background()

	plan := &model.PricingPlan{
		Name:     "Active Plan",
		Type:     constants.UserTypeLender,
		PlanType: constants.PricingPlanTypeStandard,
		Status:   constants.CommonStatusActive,
	}
	require.NoError(t, db.Create(plan).Error)

	policy := &model.PricingPlanPolicy{
		PricingPlanId:  plan.Id,
		OriginalAmount: 100,
		ActualAmount:   80,
		Type:           constants.PricingPolicyTypeMonthly,
	}
	require.NoError(t, db.Create(policy).Error)

	createInput := dto.OrganizationSubscriptionCreateDto{
		OrganizationId:      1,
		PricingPlanId:       plan.Id,
		PricingPlanPolicyId: policy.Id,
	}
	created, err := subLogic.Create(ctx, createInput)
	require.NoError(t, err)

	err = subLogic.SwitchStatus(ctx, created.Id, constants.OrganizationSubscriptionStatusSuspended)
	assert.NoError(t, err)

	var sub model.OrganizationSubscription
	require.NoError(t, db.First(&sub, created.Id).Error)
	assert.Equal(t, string(constants.OrganizationSubscriptionStatusSuspended), string(sub.Status))
}

func TestSubscriptionSwitchStatus_SameStatus_NoOp(t *testing.T) {
	subLogic, db := setupSubscriptionLogic(t)
	ctx := context.Background()

	plan := &model.PricingPlan{
		Name:     "Active Plan",
		Type:     constants.UserTypeLender,
		PlanType: constants.PricingPlanTypeStandard,
		Status:   constants.CommonStatusActive,
	}
	require.NoError(t, db.Create(plan).Error)

	policy := &model.PricingPlanPolicy{
		PricingPlanId:  plan.Id,
		OriginalAmount: 100,
		ActualAmount:   80,
		Type:           constants.PricingPolicyTypeMonthly,
	}
	require.NoError(t, db.Create(policy).Error)

	sub := &model.OrganizationSubscription{
		Code:                "SU_TEST",
		OrganizationId:      1,
		PricingPlanId:       plan.Id,
		PricingPlanPolicyId: policy.Id,
		Status:              constants.OrganizationSubscriptionStatusActive,
		OriginalAmount:      100,
		ActualAmount:        80,
	}
	require.NoError(t, db.Create(sub).Error)

	err := subLogic.SwitchStatus(ctx, sub.Id, constants.OrganizationSubscriptionStatusActive)
	assert.NoError(t, err)
}
