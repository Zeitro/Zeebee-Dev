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
	baseLogic := NewBaseLogic(repos, nil, nil, dto.JwtConfig{})
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

func TestSubscriptionSwitchStatus_ExpiredSetsExpireDateNow(t *testing.T) {
	subLogic, db := setupSubscriptionLogic(t)
	plan := &model.PricingPlan{Name: "Active Plan", Type: constants.UserTypeLender, Status: constants.CommonStatusActive}
	require.NoError(t, db.Create(plan).Error)
	sub := &model.OrganizationSubscription{
		Code:           "SU_EXPIRED_NOW",
		OrganizationId: 1,
		PricingPlanId:  plan.Id,
		Status:         constants.OrganizationSubscriptionStatusActive,
		ExpireDate:     &sql.NullTime{Time: time.Now().AddDate(0, 1, 0), Valid: true},
	}
	require.NoError(t, db.Create(sub).Error)

	before := time.Now()
	require.NoError(t, subLogic.SwitchStatus(context.Background(), sub.Id, constants.OrganizationSubscriptionStatusExpired))

	var got model.OrganizationSubscription
	require.NoError(t, db.First(&got, sub.Id).Error)
	require.NotNil(t, got.ExpireDate)
	assert.True(t, got.ExpireDate.Valid)
	assert.Equal(t, constants.OrganizationSubscriptionStatusExpired, got.Status)
	assert.WithinDuration(t, before, got.ExpireDate.Time, time.Second)
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

func TestSubscriptionUpdateExpireDate_ExpiredSnapshotsAndPrunes(t *testing.T) {
	subLogic, db := setupSubscriptionLogic(t)
	plan := &model.PricingPlan{Name: "Active Plan", Type: constants.UserTypeLender, Status: constants.CommonStatusActive}
	resource := &model.Resource{Name: "Plan resource", Type: constants.UserTypeLender, Status: constants.CommonStatusActive}
	role := &model.Role{OrganizationId: 1, Type: constants.UserTypeLender, Name: "Manager", Status: constants.CommonStatusActive}
	require.NoError(t, db.Create(plan).Error)
	require.NoError(t, db.Create(resource).Error)
	require.NoError(t, db.Create(role).Error)
	require.NoError(t, db.Exec("INSERT INTO pricing_plan_and_resources (pricing_plan_id, resource_id) VALUES (?, ?)", plan.Id, resource.Id).Error)
	require.NoError(t, db.Create(&model.RoleAndResource{RoleId: role.Id, ResourceId: resource.Id, DataScope: constants.DataScopeAll}).Error)
	sub := &model.OrganizationSubscription{Code: "SU_EXPIRE_SNAPSHOT", OrganizationId: 1, PricingPlanId: plan.Id, Status: constants.OrganizationSubscriptionStatusActive, ExpireDate: &sql.NullTime{Time: time.Now().AddDate(0, 0, 1), Valid: true}}
	require.NoError(t, db.Create(sub).Error)

	require.NoError(t, subLogic.UpdateExpireDate(context.Background(), dto.OrganizationSubscriptionUpdateExpireDateRequestDto{Id: sub.Id, ExpireDate: &sql.NullTime{Time: time.Now().Add(-time.Second), Valid: true}}))

	var count int64
	require.NoError(t, db.Model(&model.RoleAndResource{}).Where("role_id = ?", role.Id).Count(&count).Error)
	assert.Zero(t, count)
	var snapshot model.RoleResourceSnapshot
	require.NoError(t, db.Where("organization_id = ? AND pricing_plan_id = ? AND role_id = ?", 1, plan.Id, role.Id).First(&snapshot).Error)
	assert.Equal(t, []model.RoleResourceSnapshotResource{{ResourceId: resource.Id, DataScope: constants.DataScopeAll}}, snapshot.Resources)
}

func TestSuspendEffectiveSubscription_MissingPlanRollsBack(t *testing.T) {
	subLogic, db := setupSubscriptionLogic(t)
	sub := &model.OrganizationSubscription{Code: "SU_MISSING_PLAN", OrganizationId: 1, PricingPlanId: 999, Status: constants.OrganizationSubscriptionStatusActive}
	require.NoError(t, db.Create(sub).Error)

	err := subLogic.SuspendEffectiveByOrganizationId(context.Background(), 1, constants.UserTypeLender)
	require.ErrorContains(t, err, "pricing plan")

	var got model.OrganizationSubscription
	require.NoError(t, db.First(&got, sub.Id).Error)
	assert.Equal(t, constants.OrganizationSubscriptionStatusActive, got.Status)
}

func TestHandleExpiredSubscription_SamePlanDoesNotReplaySnapshot(t *testing.T) {
	subLogic, db := setupSubscriptionLogic(t)
	plan := &model.PricingPlan{Name: "Default Plan", Type: constants.UserTypeLender, Status: constants.CommonStatusActive, IsDefault: true}
	policy := &model.PricingPlanPolicy{PricingPlanId: 0, Type: constants.PricingPolicyTypeMonthly}
	resource := &model.Resource{Name: "Historical resource", Type: constants.UserTypeLender, Status: constants.CommonStatusActive}
	role := &model.Role{OrganizationId: 1, Type: constants.UserTypeLender, Name: "Manager", Status: constants.CommonStatusActive}
	require.NoError(t, db.Create(plan).Error)
	policy.PricingPlanId = plan.Id
	require.NoError(t, db.Create(policy).Error)
	require.NoError(t, db.Create(resource).Error)
	require.NoError(t, db.Create(role).Error)
	plan.Resources = []*model.Resource{resource}
	require.NoError(t, db.Create(&model.RoleResourceSnapshot{OrganizationId: 1, PricingPlanId: plan.Id, RoleId: role.Id, Resources: []model.RoleResourceSnapshotResource{{ResourceId: resource.Id, DataScope: constants.DataScopeAll}}}).Error)
	sub := &model.OrganizationSubscription{Code: "SU_SAME_PLAN_CRON", OrganizationId: 1, PricingPlanId: plan.Id, Status: constants.OrganizationSubscriptionStatusActive}
	require.NoError(t, db.Create(sub).Error)

	require.NoError(t, db.Transaction(func(tx *gorm.DB) error {
		_, err := subLogic.handleExpiredSubscription(context.Background(), tx, sub, expiredSubscriptionRecovery{DefaultPlan: plan, DefaultPolicy: policy, Now: time.Now()})
		return err
	}))

	var count int64
	require.NoError(t, db.Model(&model.RoleAndResource{}).Where("role_id = ?", role.Id).Count(&count).Error)
	assert.Zero(t, count, "cron renewal must not replay a historical permission removed during the current plan")
}

func TestSubscriptionReactivation_DisabledPlanRejected(t *testing.T) {
	subLogic, db := setupSubscriptionLogic(t)
	ctx := context.Background()
	plan := &model.PricingPlan{Name: "Disabled Plan", Type: constants.UserTypeLender, Status: constants.CommonStatusDisabled}
	require.NoError(t, db.Create(plan).Error)

	suspended := &model.OrganizationSubscription{Code: "SU_DISABLED_SWITCH", OrganizationId: 1, PricingPlanId: plan.Id, Status: constants.OrganizationSubscriptionStatusSuspended}
	require.NoError(t, db.Create(suspended).Error)
	err := subLogic.SwitchStatus(ctx, suspended.Id, constants.OrganizationSubscriptionStatusActive)
	require.ErrorContains(t, err, "disabled")
	var got model.OrganizationSubscription
	require.NoError(t, db.First(&got, suspended.Id).Error)
	assert.Equal(t, constants.OrganizationSubscriptionStatusSuspended, got.Status)

	expired := &model.OrganizationSubscription{Code: "SU_DISABLED_DATE", OrganizationId: 1, PricingPlanId: plan.Id, Status: constants.OrganizationSubscriptionStatusExpired}
	require.NoError(t, db.Create(expired).Error)
	err = subLogic.UpdateExpireDate(ctx, dto.OrganizationSubscriptionUpdateExpireDateRequestDto{Id: expired.Id, ExpireDate: &sql.NullTime{Time: time.Now().AddDate(0, 0, 1), Valid: true}})
	require.ErrorContains(t, err, "disabled")
	var expiredGot model.OrganizationSubscription
	require.NoError(t, db.First(&expiredGot, expired.Id).Error)
	assert.Equal(t, constants.OrganizationSubscriptionStatusExpired, expiredGot.Status)
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
// before touching any role/resource table.
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
		return syncSubscriptionEffect(ctx, repos, tx, 1, plan, true, true)
	})
	require.NoError(t, err)
}

// When the plan has no resources bound, syncSubscriptionEffect must short-circuit
// even if the effective boundary was crossed.
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
		return syncSubscriptionEffect(ctx, repos, tx, 1, plan, false, true)
	})
	require.NoError(t, err)
}

func TestRoleResourceSnapshotRepo_SaveListAndCleanup(t *testing.T) {
	db := setupTestDB(t)
	repos := repo.NewAllRepos(db)
	ctx := context.Background()
	require.NoError(t, db.Transaction(func(tx *gorm.DB) error {
		if err := repos.RoleResourceSnapshotRepo.SaveWithTx(ctx, tx, &model.RoleResourceSnapshot{OrganizationId: 1, PricingPlanId: 2, RoleId: 3, Resources: []model.RoleResourceSnapshotResource{{ResourceId: 4, DataScope: constants.DataScopeAll}}}); err != nil {
			return err
		}
		return repos.RoleResourceSnapshotRepo.SaveWithTx(ctx, tx, &model.RoleResourceSnapshot{OrganizationId: 1, PricingPlanId: 2, RoleId: 3, Resources: []model.RoleResourceSnapshotResource{}})
	}))

	var snapshots []*model.RoleResourceSnapshot
	require.NoError(t, db.Transaction(func(tx *gorm.DB) error {
		var err error
		snapshots, err = repos.RoleResourceSnapshotRepo.FindByOrganizationAndPlanWithTx(ctx, tx, 1, 2)
		return err
	}))
	require.Len(t, snapshots, 1)
	assert.Empty(t, snapshots[0].Resources)
	require.NoError(t, db.Transaction(func(tx *gorm.DB) error {
		return repos.RoleResourceSnapshotRepo.DeleteByRoleIdWithTx(ctx, tx, 3)
	}))
	var count int64
	require.NoError(t, db.Model(&model.RoleResourceSnapshot{}).Count(&count).Error)
	assert.Zero(t, count)
	require.NoError(t, db.Create(&model.RoleResourceSnapshot{OrganizationId: 1, PricingPlanId: 2, RoleId: 4, Resources: []model.RoleResourceSnapshotResource{{ResourceId: 5, DataScope: constants.DataScopeAll}}}).Error)
	require.NoError(t, db.Transaction(func(tx *gorm.DB) error {
		return repos.RoleResourceSnapshotRepo.DeleteByPricingPlanIdWithTx(ctx, tx, 2)
	}))
	require.NoError(t, db.Model(&model.RoleResourceSnapshot{}).Count(&count).Error)
	assert.Zero(t, count)
}

func TestReconcileRoleResourcesForPlanChange_RestoresSnapshotWithinPlan(t *testing.T) {
	db := setupTestDB(t)
	repos := repo.NewAllRepos(db)
	ctx := context.Background()
	role := &model.Role{OrganizationId: 1, Type: constants.UserTypeLender, Name: "Manager", Status: constants.CommonStatusActive}
	resource := &model.Resource{Type: constants.UserTypeLender, Name: "Plan resource", Status: constants.CommonStatusActive}
	require.NoError(t, db.Create(role).Error)
	require.NoError(t, db.Create(resource).Error)
	require.NoError(t, db.Exec("INSERT INTO role_and_resources (role_id, resource_id) VALUES (?, ?)", role.Id, resource.Id).Error)
	plan := &model.PricingPlan{Base: model.Base{Id: 9}, Type: constants.UserTypeLender, Resources: []*model.Resource{resource}}
	require.NoError(t, db.Transaction(func(tx *gorm.DB) error {
		return reconcileRoleResourcesForPlanChange(ctx, repos, tx, 1, plan, nil)
	}))
	var snapshot model.RoleResourceSnapshot
	require.NoError(t, db.Where("organization_id = ? AND pricing_plan_id = ? AND role_id = ?", 1, 9, role.Id).First(&snapshot).Error)
	assert.Equal(t, []model.RoleResourceSnapshotResource{{ResourceId: resource.Id, DataScope: constants.DataScopeAll}}, snapshot.Resources)
	disabled := &model.Resource{Type: constants.UserTypeLender, Name: "Disabled", Status: constants.CommonStatusDisabled}
	require.NoError(t, db.Create(disabled).Error)
	snapshot.Resources = []model.RoleResourceSnapshotResource{{ResourceId: resource.Id, DataScope: constants.DataScopeAll}, {ResourceId: disabled.Id, DataScope: constants.DataScopeAll}}
	require.NoError(t, db.Save(&snapshot).Error)
	var count int64
	require.NoError(t, db.Model(&model.RoleAndResource{}).Where("role_id = ?", role.Id).Count(&count).Error)
	assert.Zero(t, count)
	targetPlan := &model.PricingPlan{Base: model.Base{Id: 9}, Type: constants.UserTypeLender, Resources: []*model.Resource{resource, disabled}}
	require.NoError(t, db.Transaction(func(tx *gorm.DB) error {
		return reconcileRoleResourcesForPlanChange(ctx, repos, tx, 1, nil, targetPlan)
	}))
	require.NoError(t, db.Model(&model.RoleAndResource{}).Where("role_id = ?", role.Id).Count(&count).Error)
	assert.Equal(t, int64(1), count)
}

func TestReconcileRoleResourcesForPlanChange_RestoreDataScopeAllWins(t *testing.T) {
	db := setupTestDB(t)
	repos := repo.NewAllRepos(db)
	ctx := context.Background()
	resource := &model.Resource{Type: constants.UserTypeLender, Name: "Plan resource", Status: constants.CommonStatusActive}
	roleWithCurrentAll := &model.Role{OrganizationId: 1, Type: constants.UserTypeLender, Name: "Current All", Status: constants.CommonStatusActive}
	roleWithSnapshotAll := &model.Role{OrganizationId: 1, Type: constants.UserTypeLender, Name: "Snapshot All", Status: constants.CommonStatusActive}
	require.NoError(t, db.Create(resource).Error)
	require.NoError(t, db.Create(roleWithCurrentAll).Error)
	require.NoError(t, db.Create(roleWithSnapshotAll).Error)
	require.NoError(t, db.Create(&model.RoleAndResource{RoleId: roleWithCurrentAll.Id, ResourceId: resource.Id, DataScope: constants.DataScopeAll}).Error)
	require.NoError(t, db.Create(&model.RoleAndResource{RoleId: roleWithSnapshotAll.Id, ResourceId: resource.Id, DataScope: constants.DataScopeOrganization}).Error)
	require.NoError(t, db.Create(&model.RoleResourceSnapshot{OrganizationId: 1, PricingPlanId: 9, RoleId: roleWithCurrentAll.Id, Resources: []model.RoleResourceSnapshotResource{{ResourceId: resource.Id, DataScope: constants.DataScopeOwn}}}).Error)
	require.NoError(t, db.Create(&model.RoleResourceSnapshot{OrganizationId: 1, PricingPlanId: 9, RoleId: roleWithSnapshotAll.Id, Resources: []model.RoleResourceSnapshotResource{{ResourceId: resource.Id, DataScope: constants.DataScopeAll}}}).Error)

	plan := &model.PricingPlan{Base: model.Base{Id: 9}, Type: constants.UserTypeLender, Resources: []*model.Resource{resource}}
	require.NoError(t, db.Transaction(func(tx *gorm.DB) error {
		return reconcileRoleResourcesForPlanChange(ctx, repos, tx, 1, nil, plan)
	}))

	for _, role := range []*model.Role{roleWithCurrentAll, roleWithSnapshotAll} {
		var relation model.RoleAndResource
		require.NoError(t, db.Where("role_id = ? AND resource_id = ?", role.Id, resource.Id).First(&relation).Error)
		assert.Equal(t, constants.DataScopeAll, relation.DataScope)
	}
}
