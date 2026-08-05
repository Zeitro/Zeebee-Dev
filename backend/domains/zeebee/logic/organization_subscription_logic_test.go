package logic

import (
	"context"
	"database/sql"
	"fmt"
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
	baseLogic := NewBaseLogic(repos, nil, nil, nil)
	authTopicMap := map[constants.UserType]string{
		constants.UserTypeBorrower: "test-borrower-topic",
		constants.UserTypeLender:   "test-lender-topic",
	}
	// Create (H4) validates the org's type against the plan type, and every test
	// in this suite subscribes org 1 to a Lender plan, so seed a matching org.
	seedOrgWithType(t, db, 1, constants.UserTypeLender)
	return NewOrganizationSubscriptionLogic(baseLogic, authTopicMap), db
}

// seedOrgWithType inserts an organization with the given id and type. Create flows
// now validate the org's type against the plan type (H4), so tests exercising
// Create/Renew/SwitchStatus must seed a matching organization first.
func seedOrgWithType(t *testing.T, db *gorm.DB, id int64, orgType constants.UserType) {
	t.Helper()
	org := &model.Organization{
		Name:   fmt.Sprintf("Org %d", id),
		Status: constants.OrganizationStatusActive,
		Type:   orgType,
	}
	require.NoError(t, db.Create(org).Error)
	// force the deterministic id the tests pass as OrganizationId
	require.NoError(t, db.Model(&model.Organization{}).Where("id = ?", org.Id).Update("id", id).Error)
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
		Type:                constants.UserTypeLender,
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
		Type:                constants.UserTypeLender,
	}

	_, err := subLogic.Create(ctx, input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "pricing plan policy")
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
		Type:                constants.UserTypeLender,
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
		Type:                constants.UserTypeLender,
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
		Type:                constants.UserTypeLender,
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
		Type:                constants.UserTypeLender,
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
		Type:                constants.UserTypeLender,
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
		Type:                constants.UserTypeLender,
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

// --- H4: organization type must match plan type ---

// Create rejects a subscription whose organization type differs from the plan
// type, so plan resources cannot be bound to the wrong user-type roles.
func TestSubscriptionCreate_OrgTypeMismatch(t *testing.T) {
	subLogic, db := setupSubscriptionLogic(t)
	ctx := context.Background()

	// setupSubscriptionLogic seeds org 1 as Lender; create a Borrower plan.
	plan := &model.PricingPlan{
		Name:     "Borrower Plan",
		Type:     constants.UserTypeBorrower,
		PlanType: constants.PricingPlanTypeStandard,
		Status:   constants.CommonStatusActive,
	}
	require.NoError(t, db.Create(plan).Error)
	policy := &model.PricingPlanPolicy{
		PricingPlanId:  plan.Id,
		OriginalAmount: 10,
		ActualAmount:   5,
		Type:           constants.PricingPolicyTypeMonthly,
	}
	require.NoError(t, db.Create(policy).Error)

	input := dto.OrganizationSubscriptionCreateDto{
		OrganizationId:      1,
		PricingPlanId:       plan.Id,
		PricingPlanPolicyId: policy.Id,
		Type:                constants.UserTypeBorrower,
	}

	_, err := subLogic.Create(ctx, input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "type not match")
}

// --- M2: UpdateExpireDate must not override Suspended ---

// A Suspended subscription stays Suspended when its expire date changes; only
// non-Suspended statuses are derived from the new expire date.
func TestSubscriptionUpdateExpireDate_SuspendedPreserved(t *testing.T) {
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
		Code:                "SU_SUS",
		OrganizationId:      1,
		PricingPlanId:       plan.Id,
		PricingPlanPolicyId: policy.Id,
		Status:              constants.OrganizationSubscriptionStatusSuspended,
		ExpireDate:          &sql.NullTime{Time: time.Now().AddDate(0, 0, 5), Valid: true},
		OriginalAmount:      100,
		ActualAmount:        80,
	}
	require.NoError(t, db.Create(sub).Error)

	input := dto.OrganizationSubscriptionUpdateExpireDateRequestDto{
		Id:         sub.Id,
		ExpireDate: &sql.NullTime{Time: time.Now().AddDate(0, 0, 30), Valid: true},
	}
	require.NoError(t, subLogic.UpdateExpireDate(ctx, input))

	var result model.OrganizationSubscription
	require.NoError(t, db.First(&result, sub.Id).Error)
	assert.Equal(t, string(constants.OrganizationSubscriptionStatusSuspended), string(result.Status))
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

// --- syncSubscriptionEffect: refactored no-op branches ---

// When effective status is unchanged, syncSubscriptionEffect must short-circuit
// before touching any role/resource table. Asserts no error and nil result.
func TestSyncSubscriptionEffect_NoChange_NoOp(t *testing.T) {
	db := setupTestDB(t)
	repos := repo.NewAllRepos(db)
	ctx := context.Background()

	plan := &model.PricingPlan{
		Name:     "Active Plan",
		Type:     constants.UserTypeLender,
		PlanType: constants.PricingPlanTypeStandard,
		Status:   constants.CommonStatusActive,
	}
	require.NoError(t, db.Create(plan).Error)

	err := db.Transaction(func(tx *gorm.DB) error {
		result, err := syncSubscriptionEffect(ctx, repos, tx, subscriptionEffectSyncReq{
			OrganizationId: 1,
			Plan:           plan,
			OldEffective:   true,
			NewEffective:   true,
		})
		assert.NoError(t, err)
		assert.Nil(t, result)
		return nil
	})
	require.NoError(t, err)
}

// When the plan has no resources bound, syncSubscriptionEffect must short-circuit
// even if the effective boundary was crossed. Asserts no error and nil result.
func TestSyncSubscriptionEffect_NoResources_NoOp(t *testing.T) {
	db := setupTestDB(t)
	repos := repo.NewAllRepos(db)
	ctx := context.Background()

	plan := &model.PricingPlan{
		Name:     "Plan Without Resources",
		Type:     constants.UserTypeLender,
		PlanType: constants.PricingPlanTypeStandard,
		Status:   constants.CommonStatusActive,
	}
	require.NoError(t, db.Create(plan).Error)

	err := db.Transaction(func(tx *gorm.DB) error {
		result, err := syncSubscriptionEffect(ctx, repos, tx, subscriptionEffectSyncReq{
			OrganizationId: 1,
			Plan:           plan,
			OldEffective:   false,
			NewEffective:   true,
		})
		assert.NoError(t, err)
		assert.Nil(t, result)
		return nil
	})
	require.NoError(t, err)
}

// --- unbindEffectiveSubscriptions: non-effective filtering ---

// seedUnbindFixture seeds a plan+resource, a lender role for org=1 with the
// resource bound, and one subscription per given status. Returns the role.
func seedUnbindFixture(t *testing.T, db *gorm.DB, subStatuses []constants.OrganizationSubscriptionStatus) *model.Role {
	t.Helper()

	plan := &model.PricingPlan{
		Name:     "Active Plan",
		Type:     constants.UserTypeLender,
		PlanType: constants.PricingPlanTypeStandard,
		Status:   constants.CommonStatusActive,
	}
	require.NoError(t, db.Create(plan).Error)

	resource := &model.Resource{
		Name:   "Res1",
		Type:   constants.UserTypeLender,
		Status: constants.CommonStatusActive,
	}
	require.NoError(t, db.Create(resource).Error)

	role := &model.Role{
		OrganizationId: 1,
		Name:           "Lender Role",
		Type:           constants.UserTypeLender,
		Status:         constants.CommonStatusActive,
	}
	require.NoError(t, db.Create(role).Error)
	require.NoError(t, db.Exec("INSERT INTO role_and_resources (role_id, resource_id) VALUES (?, ?)", role.Id, resource.Id).Error)

	// bind resource to plan so GetOne(plan, Preload Resources) yields a non-empty set
	require.NoError(t, db.Exec("INSERT INTO pricing_plan_and_resources (pricing_plan_id, resource_id) VALUES (?, ?)", plan.Id, resource.Id).Error)

	for i, st := range subStatuses {
		sub := &model.OrganizationSubscription{
			Code:           fmt.Sprintf("SU_%d", i),
			OrganizationId: 1,
			PricingPlanId:  plan.Id,
			Status:         st,
			ExpireDate:     &sql.NullTime{Time: time.Now().AddDate(0, 0, 10), Valid: true},
			OriginalAmount: 100,
			ActualAmount:   80,
		}
		require.NoError(t, db.Create(sub).Error)
	}
	return role
}

// When every subscription is non-effective (Suspended/Expired), unbindEffectiveSubscriptions
// must skip them all: no role-resource row removed, no policy changes returned.
func TestUnbindEffectiveSubscriptions_AllNonEffective_NoOp(t *testing.T) {
	db := setupTestDB(t)
	repos := repo.NewAllRepos(db)
	ctx := context.Background()

	role := seedUnbindFixture(t, db, []constants.OrganizationSubscriptionStatus{
		constants.OrganizationSubscriptionStatusSuspended,
		constants.OrganizationSubscriptionStatusExpired,
	})

	var policies map[int64]*organizationPolicyChanges
	err := db.Transaction(func(tx *gorm.DB) error {
		var subs []*model.OrganizationSubscription
		require.NoError(t, tx.WithContext(ctx).Find(&subs).Error)

		var innerErr error
		policies, innerErr = unbindEffectiveSubscriptions(ctx, repos, tx, subs)
		return innerErr
	})
	require.NoError(t, err)

	assert.Empty(t, policies)

	var count int64
	require.NoError(t, db.Table("role_and_resources").Where("role_id = ?", role.Id).Count(&count).Error)
	assert.Equal(t, int64(1), count, "non-effective subs must not unbind resources")
}

// When the subscription set mixes effective and non-effective statuses, only
// effective ones trigger unbind: the role loses its resource and policy changes
// are emitted exactly once for the org (same-plan subs collapse to one unbind).
func TestUnbindEffectiveSubscriptions_Mixed_OnlyEffectiveUnbound(t *testing.T) {
	db := setupTestDB(t)
	repos := repo.NewAllRepos(db)
	ctx := context.Background()

	role := seedUnbindFixture(t, db, []constants.OrganizationSubscriptionStatus{
		constants.OrganizationSubscriptionStatusActive,
		constants.OrganizationSubscriptionStatusSuspended,
		constants.OrganizationSubscriptionStatusTrialing,
		constants.OrganizationSubscriptionStatusExpired,
	})

	var policies map[int64]*organizationPolicyChanges
	err := db.Transaction(func(tx *gorm.DB) error {
		var subs []*model.OrganizationSubscription
		require.NoError(t, tx.WithContext(ctx).Find(&subs).Error)

		var innerErr error
		policies, innerErr = unbindEffectiveSubscriptions(ctx, repos, tx, subs)
		return innerErr
	})
	require.NoError(t, err)

	require.Contains(t, policies, int64(1))
	assert.NotEmpty(t, policies[1].removed)

	var count int64
	require.NoError(t, db.Table("role_and_resources").Where("role_id = ?", role.Id).Count(&count).Error)
	assert.Equal(t, int64(0), count, "effective subs must unbind resources")
}
