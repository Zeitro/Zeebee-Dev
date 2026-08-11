package organizationsubscriptionlogic

import (
	"context"
	"database/sql"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"zeebee-backend/app/tenant/rpc/internal/logic"
	"zeebee-backend/app/tenant/rpc/internal/logic/testutil"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"
)

var subCodeSeq int64

func newTestEnv(t *testing.T) *testutil.Env {
	t.Helper()
	env := testutil.NewEnv(t)
	logic.ResetBaseLogicForTest()
	return env
}

func seedOrg(t *testing.T, env *testutil.Env, userType constants.UserType) *model.Organization {
	t.Helper()
	return testutil.SeedOrganization(t, env.DB, &model.Organization{
		Name: "Org-" + string(userType), Status: constants.OrganizationStatusActive, Type: userType,
	})
}

func seedPlan(t *testing.T, env *testutil.Env, userType constants.UserType, name string) *model.PricingPlan {
	t.Helper()
	return testutil.SeedPricingPlan(t, env.DB, &model.PricingPlan{
		Name: name, Status: constants.CommonStatusActive, Type: userType,
	})
}

// seedPolicy binds a monthly policy to a plan (Type left zero -> gorm default
// "Monthly"; calculateExpireDate treats the default as one month).
func seedPolicy(t *testing.T, env *testutil.Env, plan *model.PricingPlan) *model.PricingPlanPolicy {
	t.Helper()
	policy := &model.PricingPlanPolicy{
		PricingPlanId:  plan.Id,
		OriginalAmount: 100,
		ActualAmount:   50,
	}
	require.NoError(t, env.DB.Create(policy).Error)
	return policy
}

func seedSubscription(t *testing.T, env *testutil.Env, org *model.Organization, plan *model.PricingPlan, status constants.OrganizationSubscriptionStatus) *model.OrganizationSubscription {
	t.Helper()
	return testutil.SeedSubscription(t, env.DB, &model.OrganizationSubscription{
		Code:           fmt.Sprintf("SUB-%d", atomic.AddInt64(&subCodeSeq, 1)),
		OrganizationId: org.Id,
		PricingPlanId:  plan.Id,
		Status:         status,
		ExpireDate:     &sql.NullTime{Time: time.Now().AddDate(0, 1, 0), Valid: true},
	})
}

// T037 US-10: Create activates a subscription bound to a plan + policy.
func TestCreate_ActiveSubscription(t *testing.T) {
	env := newTestEnv(t)
	org := seedOrg(t, env, constants.UserTypeLender)
	plan := seedPlan(t, env, constants.UserTypeLender, "Plan")
	policy := seedPolicy(t, env, plan)

	resp, err := NewCreateLogic(context.Background(), env.SvcCtx).Create(&tenant.OrganizationSubscriptionCreateRequest{Data: &tenant.OrganizationSubscriptionCreateRequestData{
		UserType: string(constants.UserTypeLender), OrganizationId: org.Id,
		PricingPlanId: plan.Id, PricingPlanPolicyId: policy.Id,
	}})
	require.NoError(t, err)
	require.NotNil(t, resp.GetData())

	var sub model.OrganizationSubscription
	require.NoError(t, env.DB.First(&sub, resp.GetData().GetId()).Error)
	assert.Equal(t, constants.OrganizationSubscriptionStatusActive, sub.Status)
	assert.Equal(t, org.Id, sub.OrganizationId)
	assert.Equal(t, plan.Id, sub.PricingPlanId)
}

// T037 US-10: Create rejects a plan whose type does not match the request.
func TestCreate_PlanTypeMismatch(t *testing.T) {
	env := newTestEnv(t)
	org := seedOrg(t, env, constants.UserTypeLender)
	plan := seedPlan(t, env, constants.UserTypeBorrower, "BorrowerPlan") // wrong type

	_, err := NewCreateLogic(context.Background(), env.SvcCtx).Create(&tenant.OrganizationSubscriptionCreateRequest{Data: &tenant.OrganizationSubscriptionCreateRequestData{
		UserType: string(constants.UserTypeLender), OrganizationId: org.Id, PricingPlanId: plan.Id,
	}})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "type not match")
}

// T037 US-10: Detail returns a subscription by id.
func TestDetail_ReturnsSubscription(t *testing.T) {
	env := newTestEnv(t)
	org := seedOrg(t, env, constants.UserTypeLender)
	plan := seedPlan(t, env, constants.UserTypeLender, "Plan")
	sub := seedSubscription(t, env, org, plan, constants.OrganizationSubscriptionStatusActive)

	resp, err := NewDetailLogic(context.Background(), env.SvcCtx).Detail(&tenant.OrganizationSubscriptionDetailRequest{Data: &tenant.OrganizationSubscriptionDetailRequestData{
		UserType: string(constants.UserTypeLender), Id: sub.Id, OrganizationId: org.Id,
	}})
	require.NoError(t, err)
	require.NotNil(t, resp.GetData())
	assert.Equal(t, sub.Id, resp.GetData().GetId())
}

func TestDetail_RejectsAnotherOrganization(t *testing.T) {
	env := newTestEnv(t)
	org := seedOrg(t, env, constants.UserTypeLender)
	otherOrg := testutil.SeedOrganization(t, env.DB, &model.Organization{Name: "Other", Status: constants.OrganizationStatusActive, Type: constants.UserTypeLender})
	plan := seedPlan(t, env, constants.UserTypeLender, "Plan")
	sub := seedSubscription(t, env, org, plan, constants.OrganizationSubscriptionStatusActive)

	_, err := NewDetailLogic(context.Background(), env.SvcCtx).Detail(&tenant.OrganizationSubscriptionDetailRequest{Data: &tenant.OrganizationSubscriptionDetailRequestData{
		UserType: string(constants.UserTypeLender), Id: sub.Id, OrganizationId: otherOrg.Id,
	}})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not belong")
}

// T037 US-10: SwitchStatus suspends a subscription.
func TestSwitchStatus_Suspends(t *testing.T) {
	env := newTestEnv(t)
	org := seedOrg(t, env, constants.UserTypeLender)
	plan := seedPlan(t, env, constants.UserTypeLender, "Plan")
	sub := seedSubscription(t, env, org, plan, constants.OrganizationSubscriptionStatusActive)

	_, err := NewSwitchStatusLogic(context.Background(), env.SvcCtx).SwitchStatus(&tenant.OrganizationSubscriptionSwitchStatusRequest{Data: &tenant.OrganizationSubscriptionSwitchStatusRequestData{
		UserType: string(constants.UserTypeLender), Id: sub.Id, Status: string(constants.OrganizationSubscriptionStatusSuspended),
	}})
	require.NoError(t, err)

	var got model.OrganizationSubscription
	require.NoError(t, env.DB.First(&got, sub.Id).Error)
	assert.Equal(t, constants.OrganizationSubscriptionStatusSuspended, got.Status)
}

// T037 US-10: Page lists subscriptions (filter/sorter optional).
func TestPage_ListsSubscriptions(t *testing.T) {
	env := newTestEnv(t)
	org := seedOrg(t, env, constants.UserTypeLender)
	plan := seedPlan(t, env, constants.UserTypeLender, "Plan")
	seedSubscription(t, env, org, plan, constants.OrganizationSubscriptionStatusActive)
	seedSubscription(t, env, org, plan, constants.OrganizationSubscriptionStatusActive)
	otherOrg := testutil.SeedOrganization(t, env.DB, &model.Organization{Name: "Other", Status: constants.OrganizationStatusActive, Type: constants.UserTypeLender})
	seedSubscription(t, env, otherOrg, plan, constants.OrganizationSubscriptionStatusActive)

	resp, err := NewPageLogic(context.Background(), env.SvcCtx).Page(&tenant.OrganizationSubscriptionPageRequest{Data: &tenant.OrganizationSubscriptionPageRequestData{
		Page: 1, Limit: 10, UserType: string(constants.UserTypeLender), OrganizationId: org.Id,
		// Filter and Sorter intentionally omitted.
	}})
	require.NoError(t, err)
	require.NotNil(t, resp.GetData())
	assert.Equal(t, int64(2), resp.GetData().GetTotal())
}

// T037 US-10: UpdateExpireDate overrides the subscription expire date.
func TestUpdateExpireDate_ChangesExpire(t *testing.T) {
	env := newTestEnv(t)
	org := seedOrg(t, env, constants.UserTypeLender)
	plan := seedPlan(t, env, constants.UserTypeLender, "Plan")
	sub := seedSubscription(t, env, org, plan, constants.OrganizationSubscriptionStatusActive)

	_, err := NewUpdateExpireDateLogic(context.Background(), env.SvcCtx).UpdateExpireDate(&tenant.OrganizationSubscriptionUpdateExpireDateRequest{Data: &tenant.OrganizationSubscriptionUpdateExpireDateRequestData{
		UserType: string(constants.UserTypeLender), Id: sub.Id, ExpireDate: "2027-01-01T00:00:00Z",
	}})
	require.NoError(t, err)

	var got model.OrganizationSubscription
	require.NoError(t, env.DB.First(&got, sub.Id).Error)
	require.NotNil(t, got.ExpireDate)
	assert.True(t, got.ExpireDate.Valid)
	assert.Equal(t, "2027-01-01 00:00:00 +0000 UTC", got.ExpireDate.Time.UTC().String())
}

// T037 US-10: SuspendEffectiveByOrganizationId suspends an org's effective subs.
func TestSuspendEffectiveByOrganizationId_SuspendsActive(t *testing.T) {
	env := newTestEnv(t)
	org := seedOrg(t, env, constants.UserTypeLender)
	plan := seedPlan(t, env, constants.UserTypeLender, "Plan")
	sub := seedSubscription(t, env, org, plan, constants.OrganizationSubscriptionStatusActive)

	_, err := NewSuspendEffectiveByOrganizationIdLogic(context.Background(), env.SvcCtx).SuspendEffectiveByOrganizationId(&tenant.OrganizationSubscriptionSuspendEffectiveByOrganizationIdRequest{Data: &tenant.OrganizationSubscriptionSuspendEffectiveByOrganizationIdRequestData{
		UserType: string(constants.UserTypeLender), OrganizationId: org.Id,
	}})
	require.NoError(t, err)

	var got model.OrganizationSubscription
	require.NoError(t, env.DB.First(&got, sub.Id).Error)
	assert.Equal(t, constants.OrganizationSubscriptionStatusSuspended, got.Status)
}
