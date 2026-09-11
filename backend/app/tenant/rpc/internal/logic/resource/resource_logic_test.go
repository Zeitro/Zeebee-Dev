package resourcelogic

import (
	"context"
	"database/sql"
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

func newTestEnv(t *testing.T) *testutil.Env {
	t.Helper()
	env := testutil.NewEnv(t)
	logic.ResetBaseLogicForTest()
	return env
}

// T040 US-16: resource CRUD round-trips through the tenant RPC logic.
func TestResourceCRUD(t *testing.T) {
	env := newTestEnv(t)

	createResp, err := NewCreateLogic(context.Background(), env.SvcCtx).Create(&tenant.ResourceCreateRequest{Data: &tenant.ResourceCreateRequestData{
		UserType: string(constants.UserTypeAdmin), Type: string(constants.UserTypeAdmin),
		ResourceType: string(constants.ResourceTypeMenu), Name: "Reports",
		Status: string(constants.CommonStatusActive),
	}})
	require.NoError(t, err)
	require.NotNil(t, createResp.GetData())
	id := createResp.GetData().GetId()

	detailResp, err := NewDetailLogic(context.Background(), env.SvcCtx).Detail(&tenant.ResourceDetailRequest{Data: &tenant.ResourceDetailRequestData{
		UserType: string(constants.UserTypeAdmin), Type: string(constants.UserTypeAdmin), Id: id,
	}})
	require.NoError(t, err)
	assert.Equal(t, id, detailResp.GetData().GetId())

	treeResp, err := NewResourceTreeLogic(context.Background(), env.SvcCtx).ResourceTree(&tenant.ResourceTreeRequest{Data: &tenant.ResourceTreeRequestData{
		UserType: string(constants.UserTypeAdmin), Type: string(constants.UserTypeAdmin), Scope: string(constants.ResourceTreeScopeAll),
	}})
	require.NoError(t, err)
	assert.NotEmpty(t, treeResp.GetData())

	updateResp, err := NewUpdateLogic(context.Background(), env.SvcCtx).Update(&tenant.ResourceUpdateRequest{Data: &tenant.ResourceUpdateRequestData{
		UserType: string(constants.UserTypeAdmin), Type: string(constants.UserTypeAdmin), Id: id, Name: "Reports2",
	}})
	require.NoError(t, err)
	assert.Equal(t, "Reports2", updateResp.GetData().GetName())

	_, err = NewDeleteLogic(context.Background(), env.SvcCtx).Delete(&tenant.ResourceDeleteRequest{Data: &tenant.ResourceDeleteRequestData{
		UserType: string(constants.UserTypeAdmin), Type: string(constants.UserTypeAdmin), Id: id,
	}})
	require.NoError(t, err)
}

func TestResourceTree_RejectsCrossPlatformScope(t *testing.T) {
	logic := &ResourceTreeLogic{}

	_, err := logic.ResourceTree(&tenant.ResourceTreeRequest{Data: &tenant.ResourceTreeRequestData{
		UserType: string(constants.UserTypeLender),
		Type:     string(constants.UserTypeBorrower),
		Scope:    string(constants.ResourceTreeScopeAll),
	}})

	require.EqualError(t, err, "resource type must match caller platform")
}

func TestUserMenu_OwnerIsPlanBound(t *testing.T) {
	env := newTestEnv(t)
	org := testutil.SeedOrganization(t, env.DB, &model.Organization{Name: "Owner org", OwnerId: 1, Type: constants.UserTypeLender, Status: constants.OrganizationStatusActive})
	member := testutil.SeedOrganizationUser(t, env.DB, &model.OrganizationUser{OrganizationId: org.Id, UserId: 1, Status: constants.OrganizationUserStatusActive})
	resource := testutil.SeedResource(t, env.DB, &model.Resource{Type: constants.UserTypeLender, ResourceType: constants.ResourceTypeMenu, Name: "Plan menu", Status: constants.CommonStatusActive})
	plan := testutil.SeedPricingPlan(t, env.DB, &model.PricingPlan{Name: "Plan", Type: constants.UserTypeLender, Status: constants.CommonStatusActive})
	require.NoError(t, env.DB.Create(&model.PricingPlanAndResource{PricingPlanId: plan.Id, ResourceId: resource.Id}).Error)
	testutil.SeedSubscription(t, env.DB, &model.OrganizationSubscription{Code: "MENU", OrganizationId: org.Id, PricingPlanId: plan.Id, Status: constants.OrganizationSubscriptionStatusActive, ExpireDate: &sql.NullTime{Time: time.Now().AddDate(0, 0, 1), Valid: true}})

	menu := NewUserMenuLogic(context.Background(), env.SvcCtx)
	resp, err := menu.UserMenu(&tenant.UserMenuRequest{Data: &tenant.UserMenuRequestData{UserType: string(constants.UserTypeLender), OrganizationId: org.Id, OrganizationUserId: member.Id}})
	require.NoError(t, err)
	require.Len(t, resp.GetData(), 1)
	require.NoError(t, env.DB.Where("pricing_plan_id = ? AND resource_id = ?", plan.Id, resource.Id).Delete(&model.PricingPlanAndResource{}).Error)
	resp, err = menu.UserMenu(&tenant.UserMenuRequest{Data: &tenant.UserMenuRequestData{UserType: string(constants.UserTypeLender), OrganizationId: org.Id, OrganizationUserId: member.Id}})
	require.NoError(t, err)
	assert.Empty(t, resp.GetData())
}
