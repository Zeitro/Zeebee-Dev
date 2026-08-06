package userlogic

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"
)

// T048 US-15: UserInfo returns the profile scoped to the user's default org.
func TestUserInfo_ReturnsProfile(t *testing.T) {
	env := newTestEnv(t)
	userType := constants.UserTypeLender
	userId, orgId := seedLoginUser(t, env, userType, "profile@test.local", testPassword)

	resp, err := NewUserInfoLogic(ctxWithUserClaim(t, env, userType, userId, orgId), env.SvcCtx).UserInfo(&tenant.UserInfoRequest{Data: &tenant.UserInfoRequestData{
		UserId: userId,
	}})
	require.NoError(t, err)
	require.NotNil(t, resp.GetData())
	assert.Equal(t, "profile@test.local", resp.GetData().GetEmail())
	assert.Equal(t, orgId, resp.GetData().GetOrganizationId())
	assert.NotEmpty(t, resp.GetData().GetOrganizations())
}

// T048 US-15: EditProfile updates the member's editable name fields
// (first/mid/last). The tenant profile contract is name-only; phone is preserved.
func TestEditProfile_UpdatesNames(t *testing.T) {
	env := newTestEnv(t)
	userType := constants.UserTypeLender
	userId, orgId := seedLoginUser(t, env, userType, "edit@test.local", testPassword)

	resp, err := NewEditProfileLogic(ctxWithUserClaim(t, env, userType, userId, orgId), env.SvcCtx).EditProfile(&tenant.EditProfileRequest{Data: &tenant.EditProfileRequestData{
		UserId:    userId,
		FirstName: "New", MidName: "M", LastName: "Name",
	}})
	require.NoError(t, err)
	require.NotNil(t, resp.GetData())

	var orgUser model.OrganizationUser
	require.NoError(t, env.DB.Where("user_id = ? AND organization_id = ?", userId, orgId).First(&orgUser).Error)
	assert.Equal(t, "New", orgUser.FirstName)
	assert.Equal(t, "M", orgUser.MidName)
	assert.Equal(t, "Name", orgUser.LastName)
}

// T048 US-15: UserInfo rejects an unknown user.
func TestUserInfo_UnknownUser(t *testing.T) {
	env := newTestEnv(t)
	_, err := NewUserInfoLogic(ctxWithUserClaim(t, env, constants.UserTypeLender, 1, 1), env.SvcCtx).UserInfo(&tenant.UserInfoRequest{Data: &tenant.UserInfoRequestData{
		UserId: 999999,
	}})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}
