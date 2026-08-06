package resourcelogic

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"zeebee-backend/app/tenant/rpc/internal/logic"
	"zeebee-backend/app/tenant/rpc/internal/logic/testutil"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/domains/zeebee/constants"
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
