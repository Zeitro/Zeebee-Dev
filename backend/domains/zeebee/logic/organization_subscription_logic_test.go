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
	// monthly term: expire = (now - 1 day) + 1 month, clamped to month end.
	expectedExpiry := addMonthsClamped(beforeCreate.AddDate(0, 0, -1), 1)
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

	renewed, err := subLogic.Renew(ctx, created.Id)
	require.NoError(t, err)
	assert.NotNil(t, renewed.ExpireDate)
	assert.True(t, renewed.ExpireDate.Valid)
}

func TestSubscriptionRenew_NotFound(t *testing.T) {
	subLogic, _ := setupSubscriptionLogic(t)
	ctx := context.Background()

	_, err := subLogic.Renew(ctx, 99999)
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

// TestCalculateExpireDate verifies the calendar-month billing logic:
// expire = (start - 1 day) + N months, with the day clamped to the target
// month's last day. N = 1/3/12 for monthly/quarterly/yearly.
// Cases mirror the validated spec table.
func TestCalculateExpireDate(t *testing.T) {
	parseDate := func(s string) time.Time {
		tt, err := time.ParseInLocation("2006-01-02", s, time.UTC)
		require.NoError(t, err)
		return tt
	}

	cases := []struct {
		name       string
		policyType constants.PricingPolicyType
		start      string
		want       string
	}{
		// monthly (N=1)
		{"monthly/mid-month", constants.PricingPolicyTypeMonthly, "2026-01-15", "2026-02-14"},
		{"monthly/month-start", constants.PricingPolicyTypeMonthly, "2026-01-01", "2026-01-31"},
		{"monthly/month-start-feb", constants.PricingPolicyTypeMonthly, "2026-02-01", "2026-02-28"},
		{"monthly/last-day-jan", constants.PricingPolicyTypeMonthly, "2026-01-31", "2026-02-28"},
		{"monthly/last-day-feb", constants.PricingPolicyTypeMonthly, "2026-02-28", "2026-03-27"},
		// quarterly (N=3)
		{"quarterly/mid-month", constants.PricingPolicyTypeQuarterly, "2026-01-15", "2026-04-14"},
		{"quarterly/month-start", constants.PricingPolicyTypeQuarterly, "2026-01-01", "2026-03-31"},
		{"quarterly/month-start-feb", constants.PricingPolicyTypeQuarterly, "2026-02-01", "2026-04-30"},
		{"quarterly/last-day-jan", constants.PricingPolicyTypeQuarterly, "2026-01-31", "2026-04-30"},
		{"quarterly/last-day-feb", constants.PricingPolicyTypeQuarterly, "2026-02-28", "2026-05-27"},
		{"quarterly/dec-into-short-feb", constants.PricingPolicyTypeQuarterly, "2025-12-01", "2026-02-28"},
		{"quarterly/nov-29-into-feb", constants.PricingPolicyTypeQuarterly, "2025-11-29", "2026-02-28"},
		{"quarterly/nov-30-into-feb", constants.PricingPolicyTypeQuarterly, "2025-11-30", "2026-02-28"},
		// yearly (N=12)
		{"yearly/month-start", constants.PricingPolicyTypeYearly, "2026-01-01", "2026-12-31"},
		{"yearly/mid-month", constants.PricingPolicyTypeYearly, "2026-01-15", "2027-01-14"},
		// Leap day clamps to Feb 28 of the next (non-leap) year.
		{"yearly/leap-day", constants.PricingPolicyTypeYearly, "2024-02-29", "2025-02-28"},
		{"yearly/leap-day2", constants.PricingPolicyTypeYearly, "2023-02-28", "2024-02-27"},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := calculateExpireDate(tc.policyType, parseDate(tc.start))
			require.True(t, got.Valid, "expire date should be valid")
			assert.Equal(t, tc.want, got.Time.Format("2006-01-02"))
		})
	}
}
