package apilogic

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

// T040 US-16: api registry CRUD round-trips through the tenant RPC logic.
func TestApiCRUD(t *testing.T) {
	env := newTestEnv(t)

	createResp, err := NewCreateLogic(context.Background(), env.SvcCtx).Create(&tenant.ApiCreateRequest{Data: &tenant.ApiCreateRequestData{
		UserType: string(constants.UserTypeAdmin), Type: string(constants.UserTypeAdmin),
		Name: "Widget", Path: "/api/v1/widget", Method: "GET",
		Status: string(constants.CommonStatusActive), AuthRequired: true, PermissionRequired: true,
	}})
	require.NoError(t, err)
	require.NotNil(t, createResp.GetData())
	id := createResp.GetData().GetId()

	detailResp, err := NewDetailLogic(context.Background(), env.SvcCtx).Detail(&tenant.ApiDetailRequest{Data: &tenant.ApiDetailRequestData{
		UserType: string(constants.UserTypeAdmin), Type: string(constants.UserTypeAdmin), Id: id,
	}})
	require.NoError(t, err)
	assert.Equal(t, id, detailResp.GetData().GetId())

	treeResp, err := NewApiTreeLogic(context.Background(), env.SvcCtx).ApiTree(&tenant.ApiTreeRequest{Data: &tenant.ApiTreeRequestData{
		UserType: string(constants.UserTypeAdmin), Type: string(constants.UserTypeAdmin),
	}})
	require.NoError(t, err)
	assert.NotEmpty(t, treeResp.GetData())

	updateResp, err := NewUpdateLogic(context.Background(), env.SvcCtx).Update(&tenant.ApiUpdateRequest{Data: &tenant.ApiUpdateRequestData{
		UserType: string(constants.UserTypeAdmin), Type: string(constants.UserTypeAdmin), Id: id, Name: "Widget2",
	}})
	require.NoError(t, err)
	assert.Equal(t, "Widget2", updateResp.GetData().GetName())

	_, err = NewDeleteLogic(context.Background(), env.SvcCtx).Delete(&tenant.ApiDeleteRequest{Data: &tenant.ApiDeleteRequestData{
		UserType: string(constants.UserTypeAdmin), Type: string(constants.UserTypeAdmin), Id: id,
	}})
	require.NoError(t, err)
}

func TestApiRejectsCrossPlatformTarget(t *testing.T) {
	env := newTestEnv(t)

	_, err := NewApiTreeLogic(context.Background(), env.SvcCtx).ApiTree(&tenant.ApiTreeRequest{Data: &tenant.ApiTreeRequestData{
		UserType: string(constants.UserTypeLender), Type: string(constants.UserTypeBorrower),
	}})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cannot manage another platform api")
}
