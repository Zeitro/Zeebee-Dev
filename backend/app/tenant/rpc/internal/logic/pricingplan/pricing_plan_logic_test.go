package pricingplanlogic

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

// T037 US-10: pricing plan CRUD round-trips through the tenant RPC logic.
func TestPricingPlanCRUD(t *testing.T) {
	env := newTestEnv(t)

	createResp, err := NewCreateLogic(context.Background(), env.SvcCtx).Create(&tenant.PricingPlanCreateRequest{Data: &tenant.PricingPlanCreateRequestData{
		UserType: string(constants.UserTypeLender),
		Name:     "Pro", Type: string(constants.UserTypeLender),
	}})
	require.NoError(t, err)
	require.NotNil(t, createResp.GetData())
	id := createResp.GetData().GetId()
	assert.Equal(t, "Pro", createResp.GetData().GetName())

	updateResp, err := NewUpdateLogic(context.Background(), env.SvcCtx).Update(&tenant.PricingPlanUpdateRequest{Data: &tenant.PricingPlanUpdateRequestData{
		UserType: string(constants.UserTypeLender), Id: id, Subtitle: "better",
	}})
	require.NoError(t, err)
	assert.Equal(t, "better", updateResp.GetData().GetSubtitle())

	detailResp, err := NewDetailLogic(context.Background(), env.SvcCtx).Detail(&tenant.PricingPlanDetailRequest{Data: &tenant.PricingPlanDetailRequestData{
		UserType: string(constants.UserTypeLender), Id: id,
	}})
	require.NoError(t, err)
	assert.Equal(t, id, detailResp.GetData().GetId())

	pageResp, err := NewPageLogic(context.Background(), env.SvcCtx).Page(&tenant.PricingPlanPageRequest{Data: &tenant.PricingPlanPageRequestData{
		UserType: string(constants.UserTypeLender), Page: 1, Limit: 10,
	}})
	require.NoError(t, err)
	assert.Equal(t, int64(1), pageResp.GetData().GetTotal())

	_, err = NewSwitchStatusLogic(context.Background(), env.SvcCtx).SwitchStatus(&tenant.PricingPlanSwitchStatusRequest{Data: &tenant.PricingPlanSwitchStatusRequestData{
		UserType: string(constants.UserTypeLender), Id: id, Status: string(constants.CommonStatusDisabled),
	}})
	require.NoError(t, err)

	_, err = NewDeleteLogic(context.Background(), env.SvcCtx).Delete(&tenant.PricingPlanDeleteRequest{Data: &tenant.PricingPlanDeleteRequestData{
		UserType: string(constants.UserTypeLender), Id: id,
	}})
	require.NoError(t, err)
}

// T037 US-10: BindResources + ResourceList round-trip.
func TestBindResourcesAndResourceList(t *testing.T) {
	env := newTestEnv(t)
	plan := testutil.SeedPricingPlan(t, env.DB, &model.PricingPlan{
		Name: "Pro", Type: constants.UserTypeLender, Status: constants.CommonStatusActive,
	})
	resource := testutil.SeedResource(t, env.DB, &model.Resource{
		Type: constants.UserTypeLender, Status: constants.CommonStatusActive, Name: "R1",
	})

	_, err := NewBindResourcesLogic(context.Background(), env.SvcCtx).BindResources(&tenant.PricingPlanBindResourcesRequest{Data: &tenant.PricingPlanBindResourcesRequestData{
		UserType: string(constants.UserTypeLender), Id: plan.Id, ResourceIds: []int64{resource.Id},
	}})
	require.NoError(t, err)

	listResp, err := NewResourceListLogic(context.Background(), env.SvcCtx).ResourceList(&tenant.PricingPlanResourceListRequest{Data: &tenant.PricingPlanResourceListRequestData{
		UserType: string(constants.UserTypeLender), Id: plan.Id,
	}})
	require.NoError(t, err)
	assert.Len(t, listResp.GetData(), 1)
}
