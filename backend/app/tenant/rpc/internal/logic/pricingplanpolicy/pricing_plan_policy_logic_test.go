package pricingplanpolicylogic

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

// T040 US-16: pricing plan policy CRUD + duplicate-type guard.
func TestPricingPlanPolicyCRUD(t *testing.T) {
	env := newTestEnv(t)
	plan := testutil.SeedPricingPlan(t, env.DB, &model.PricingPlan{
		Name: "Pro", Type: constants.UserTypeLender, Status: constants.CommonStatusActive,
	})

	createResp, err := NewCreateLogic(context.Background(), env.SvcCtx).Create(&tenant.PricingPlanPolicyCreateRequest{Data: &tenant.PricingPlanPolicyCreateRequestData{
		UserType: string(constants.UserTypeLender), PricingPlanId: plan.Id,
		OriginalAmount: 100, ActualAmount: 80, Type: string(constants.PricingPolicyTypeMonthly),
	}})
	require.NoError(t, err)
	require.NotNil(t, createResp.GetData())
	id := createResp.GetData().GetId()

	// same (plan, type) is rejected
	_, err = NewCreateLogic(context.Background(), env.SvcCtx).Create(&tenant.PricingPlanPolicyCreateRequest{Data: &tenant.PricingPlanPolicyCreateRequestData{
		UserType: string(constants.UserTypeLender), PricingPlanId: plan.Id,
		OriginalAmount: 100, ActualAmount: 80, Type: string(constants.PricingPolicyTypeMonthly),
	}})
	require.Error(t, err)

	updateResp, err := NewUpdateLogic(context.Background(), env.SvcCtx).Update(&tenant.PricingPlanPolicyUpdateRequest{Data: &tenant.PricingPlanPolicyUpdateRequestData{
		UserType: string(constants.UserTypeLender), Id: id, ActualAmount: 70,
	}})
	require.NoError(t, err)
	assert.Equal(t, 70.0, updateResp.GetData().GetActualAmount())

	_, err = NewDeleteLogic(context.Background(), env.SvcCtx).Delete(&tenant.PricingPlanPolicyDeleteRequest{Data: &tenant.PricingPlanPolicyDeleteRequestData{
		UserType: string(constants.UserTypeLender), Id: id,
	}})
	require.NoError(t, err)
}
