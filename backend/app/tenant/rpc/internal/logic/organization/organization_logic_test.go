package organizationlogic

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

func seedOrg(t *testing.T, env *testutil.Env, userType constants.UserType, name string, status constants.OrganizationStatus) *model.Organization {
	t.Helper()
	return testutil.SeedOrganization(t, env.DB, &model.Organization{
		Name:   name,
		Status: status,
		Type:   userType,
	})
}

// T034 US-9: Page lists organizations of a platform type.
func TestPage_ListsOrganizations(t *testing.T) {
	env := newTestEnv(t)
	seedOrg(t, env, constants.UserTypeLender, "Alpha", constants.OrganizationStatusActive)
	seedOrg(t, env, constants.UserTypeLender, "Beta", constants.OrganizationStatusActive)

	resp, err := NewPageLogic(context.Background(), env.SvcCtx).Page(&tenant.OrganizationPageRequest{Data: &tenant.OrganizationPageRequestData{
		Page: 1, Limit: 10, UserType: string(constants.UserTypeLender),
		// Filter and Sorter intentionally omitted: a page request need not carry
		// search or sort settings.
	}})
	require.NoError(t, err)
	require.NotNil(t, resp.GetData())
	assert.Equal(t, int64(2), resp.GetData().GetTotal())
	assert.Len(t, resp.GetData().GetList(), 2)
}

// T034 US-9: Detail returns an organization by id.
func TestDetail_ReturnsOrg(t *testing.T) {
	env := newTestEnv(t)
	organization := seedOrg(t, env, constants.UserTypeLender, "Detail", constants.OrganizationStatusActive)

	resp, err := NewDetailLogic(context.Background(), env.SvcCtx).Detail(&tenant.OrganizationDetailRequest{Data: &tenant.OrganizationDetailRequestData{
		UserType: string(constants.UserTypeLender), Id: organization.Id,
	}})
	require.NoError(t, err)
	require.NotNil(t, resp.GetData())
	assert.Equal(t, organization.Id, resp.GetData().GetId())
	assert.Equal(t, "Detail", resp.GetData().GetName())
}

// T034 US-9: DetailWithType rejects a request whose type does not match the org.
func TestDetailWithType_WrongTypeNotFound(t *testing.T) {
	env := newTestEnv(t)
	organization := seedOrg(t, env, constants.UserTypeLender, "Typed", constants.OrganizationStatusActive)

	_, err := NewDetailWithTypeLogic(context.Background(), env.SvcCtx).DetailWithType(&tenant.OrganizationDetailWithTypeRequest{Data: &tenant.OrganizationDetailWithTypeRequestData{
		UserType: string(constants.UserTypeBorrower), Id: organization.Id,
	}})
	require.Error(t, err)
}

// T034 US-9: All returns only organizations of the requested type.
func TestAll_ReturnsOnlyByType(t *testing.T) {
	env := newTestEnv(t)
	seedOrg(t, env, constants.UserTypeLender, "L", constants.OrganizationStatusActive)
	seedOrg(t, env, constants.UserTypeBorrower, "B", constants.OrganizationStatusActive)

	resp, err := NewAllLogic(context.Background(), env.SvcCtx).All(&tenant.OrganizationAllRequest{Data: &tenant.OrganizationAllRequestData{
		UserType: string(constants.UserTypeLender),
	}})
	require.NoError(t, err)
	assert.Len(t, resp.GetData(), 1)
}

// T034 US-9: UpdateOrder rejects a negative order. (The positive path shifts
// sibling orgs via raw SQL on the reserved column name "order", which SQLite
// rejects but PostgreSQL accepts; the order-set behavior is not asserted here.)
func TestUpdateOrder_RejectsNegative(t *testing.T) {
	env := newTestEnv(t)
	organization := seedOrg(t, env, constants.UserTypeLender, "Order", constants.OrganizationStatusActive)

	_, err := NewUpdateOrderLogic(context.Background(), env.SvcCtx).UpdateOrder(&tenant.OrganizationUpdateOrderRequest{Data: &tenant.OrganizationUpdateOrderRequestData{
		UserType: string(constants.UserTypeLender), Id: organization.Id, Order: -1,
	}})
	require.Error(t, err)
	assert.Contains(t, err.Error(), ">= 0")
}

// T034 US-9: SwitchStatus disables an active organization.
func TestSwitchStatus_DisablesOrg(t *testing.T) {
	env := newTestEnv(t)
	organization := seedOrg(t, env, constants.UserTypeLender, "Toggle", constants.OrganizationStatusActive)

	_, err := NewSwitchStatusLogic(context.Background(), env.SvcCtx).SwitchStatus(&tenant.OrganizationSwitchStatusRequest{Data: &tenant.OrganizationSwitchStatusRequestData{
		UserType: string(constants.UserTypeLender), Id: organization.Id, Status: string(constants.OrganizationStatusDisabled),
	}})
	require.NoError(t, err)

	var got model.Organization
	require.NoError(t, env.DB.First(&got, organization.Id).Error)
	assert.Equal(t, constants.OrganizationStatusDisabled, got.Status)
}

// T034 US-9: Delete removes a non-self, non-admin organization.
func TestDelete_RemovesOrg(t *testing.T) {
	env := newTestEnv(t)
	organization := seedOrg(t, env, constants.UserTypeLender, "Deletable", constants.OrganizationStatusActive)

	_, err := NewDeleteLogic(context.Background(), env.SvcCtx).Delete(&tenant.OrganizationDeleteRequest{Data: &tenant.OrganizationDeleteRequestData{
		UserType: string(constants.UserTypeLender), Id: organization.Id,
	}})
	require.NoError(t, err)

	var count int64
	require.NoError(t, env.DB.Model(&model.Organization{}).Where("id = ?", organization.Id).Count(&count).Error)
	assert.Zero(t, count)
}

// T034 US-9: Delete rejects an admin-type organization.
func TestDelete_RejectsAdminOrg(t *testing.T) {
	env := newTestEnv(t)
	organization := seedOrg(t, env, constants.UserTypeAdmin, "AdminOrg", constants.OrganizationStatusActive)

	_, err := NewDeleteLogic(context.Background(), env.SvcCtx).Delete(&tenant.OrganizationDeleteRequest{Data: &tenant.OrganizationDeleteRequestData{
		UserType: string(constants.UserTypeAdmin), Id: organization.Id,
	}})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Cannot delete")
}
