package authlogic

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"

	"zeebee-backend/app/tenant/rpc/internal/logic"
	"zeebee-backend/app/tenant/rpc/internal/logic/testutil"
	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/auth"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/proto"
)

const (
	authzMethod = "GET"
	authzPath   = "/api/v1/widget"
)

type authzFixture struct {
	userId    int64
	orgId     int64
	orgUserId int64
	roleId    int64
	resId     int64
}

func newAuthTestEnv(t *testing.T) *testutil.Env {
	t.Helper()
	env := testutil.NewEnv(t)
	logic.ResetBaseLogicForTest()
	return env
}

// seedAuthzGranted builds the full grant chain for one platform: active org,
// role bound to a protected api resource, org-user carrying the role, active
// plan covering the resource, active subscription, and the casbin policy row.
func seedAuthzGranted(t *testing.T, env *testutil.Env, ut constants.UserType, roleName string) authzFixture {
	t.Helper()
	user := testutil.SeedUser(t, env.DB, &model.User{
		Type:   ut,
		Email:  string(ut) + "@test.local",
		Status: constants.UserStatusActive,
	})
	org := testutil.SeedOrganization(t, env.DB, &model.Organization{
		Name:   "Org-" + string(ut),
		Status: constants.OrganizationStatusActive,
		Type:   ut,
	})
	role := testutil.SeedRole(t, env.DB, &model.Role{
		OrganizationId: org.Id,
		Name:           roleName,
		Type:           ut,
		Status:         constants.CommonStatusActive,
	})
	resource := testutil.SeedResource(t, env.DB, &model.Resource{
		Type:     ut,
		Status:   constants.CommonStatusActive,
		Name:     "Res" + string(ut),
		IsPublic: false,
	})
	api := &model.Api{
		Type:               ut,
		Status:             constants.CommonStatusActive,
		Name:               "Widget",
		Path:               authzPath,
		Method:             authzMethod,
		AuthRequired:       true,
		PermissionRequired: true,
		Resources:          []*model.Resource{resource},
	}
	require.NoError(t, env.DB.Create(api).Error) // fills resource_and_apis join
	testutil.SeedRoleResource(t, env.DB, role.Id, resource.Id)

	ou := &model.OrganizationUser{
		OrganizationId: org.Id,
		UserId:         user.Id,
		Status:         constants.OrganizationUserStatusActive,
		Roles:          []*model.Role{role},
	}
	require.NoError(t, env.DB.Create(ou).Error) // fills organization_user_and_roles join

	plan := testutil.SeedPricingPlan(t, env.DB, &model.PricingPlan{
		Name:   "Plan" + string(ut),
		Status: constants.CommonStatusActive,
		Type:   ut,
	})
	require.NoError(t, env.DB.Create(&model.PricingPlanAndResource{
		PricingPlanId: plan.Id,
		ResourceId:    resource.Id,
	}).Error)
	testutil.SeedSubscription(t, env.DB, &model.OrganizationSubscription{
		Code:           "SU-" + string(ut),
		OrganizationId: org.Id,
		PricingPlanId:  plan.Id,
		Status:         constants.OrganizationSubscriptionStatusActive,
		ExpireDate:     &sql.NullTime{Time: time.Now().AddDate(0, 0, 10), Valid: true},
	})

	require.NoError(t, env.DB.Exec(
		`INSERT INTO casbin_rule (ptype, v0, v1, v2, v3) VALUES ('p', ?, ?, ?, 'allow')`,
		fmt.Sprintf("Role:%d", role.Id), string(ut), fmt.Sprintf("Resource:%d", resource.Id),
	).Error)
	require.NoError(t, env.SvcCtx.Enforcer.Enforcer.LoadPolicy())

	return authzFixture{
		userId:    user.Id,
		orgId:     org.Id,
		orgUserId: ou.Id,
		roleId:    role.Id,
		resId:     resource.Id,
	}
}

// ctxWithToken signs a platform JWT for the fixture and carries it in gRPC
// metadata, as a backend would forward it to tenant.
func ctxWithToken(t *testing.T, env *testutil.Env, ut constants.UserType, f authzFixture) context.Context {
	t.Helper()
	secret := env.SvcCtx.Config.JwtConf.AccessSecret
	tok, err := utils.GenerateToken(utils.JwtUser{
		UserId:             f.userId,
		OrganizationId:     f.orgId,
		OrganizationUserId: f.orgUserId,
		UserType:           ut,
	}, secret, 3600)
	require.NoError(t, err)
	return metadata.NewIncomingContext(context.Background(), metadata.Pairs(utils.JwtTokenCtxKey, tok))
}

func authzCheckReq(ut constants.UserType) *proto.AuthCheckRequest {
	return &proto.AuthCheckRequest{Data: &proto.AuthCheckRequestData{
		Method: authzMethod,
		Path:   authzPath,
		Type:   string(ut),
	}}
}

// T019 US-4: a role+plan granted user is allowed live for lender and borrower.
func TestAuthzCheck_AllowForGrantedUser(t *testing.T) {
	for _, ut := range []constants.UserType{constants.UserTypeLender, constants.UserTypeBorrower} {
		t.Run(string(ut), func(t *testing.T) {
			env := newAuthTestEnv(t)
			f := seedAuthzGranted(t, env, ut, "Staff")
			ctx := ctxWithToken(t, env, ut, f)

			resp, err := NewAuthCheckLogic(ctx, env.SvcCtx).AuthCheck(authzCheckReq(ut))
			require.NoError(t, err)
			require.NotNil(t, resp.GetData())
			assert.True(t, resp.GetData().GetAuthenticated())
			assert.True(t, resp.GetData().GetAuthorizationRequired())
			assert.True(t, resp.GetData().GetAuthorized())
		})
	}
}

// T019 US-4: revoking the casbin policy flips the next request to denied
// (direction B: decision is live from the enforcer).
func TestAuthzCheck_RevokePolicyFlipsNextRequest(t *testing.T) {
	env := newAuthTestEnv(t)
	ut := constants.UserTypeLender
	f := seedAuthzGranted(t, env, ut, "Staff")
	ctx := ctxWithToken(t, env, ut, f)

	resp, err := NewAuthCheckLogic(ctx, env.SvcCtx).AuthCheck(authzCheckReq(ut))
	require.NoError(t, err)
	assert.True(t, resp.GetData().GetAuthorized())

	require.NoError(t, env.DB.Exec(`DELETE FROM casbin_rule WHERE v0 = ?`, fmt.Sprintf("Role:%d", f.roleId)).Error)
	require.NoError(t, env.SvcCtx.Enforcer.Enforcer.LoadPolicy())

	resp, err = NewAuthCheckLogic(ctx, env.SvcCtx).AuthCheck(authzCheckReq(ut))
	require.NoError(t, err)
	assert.True(t, resp.GetData().GetAuthenticated())
	assert.False(t, resp.GetData().GetAuthorized())
}

// T019 US-4: a role change published as a RoleReload event is reflected by the
// very next request once the subscriber refreshes the role cache.
func TestAuthzCheck_ImmediateEffectAfterRoleChange(t *testing.T) {
	env := newAuthTestEnv(t)
	ut := constants.UserTypeLender
	f := seedAuthzGranted(t, env, ut, "Staff")
	ctx := ctxWithToken(t, env, ut, f)

	resp, err := NewAuthCheckLogic(ctx, env.SvcCtx).AuthCheck(authzCheckReq(ut))
	require.NoError(t, err)
	assert.True(t, resp.GetData().GetAuthorized())

	// Revoke the org-user's role binding and policy, then publish the same
	// RoleReload event the assign/revoke write path emits.
	require.NoError(t, env.DB.Exec(`DELETE FROM organization_user_and_roles WHERE organization_user_id = ?`, f.orgUserId).Error)
	require.NoError(t, env.DB.Exec(`DELETE FROM casbin_rule WHERE v0 = ?`, fmt.Sprintf("Role:%d", f.roleId)).Error)
	msg := auth.Msg{
		Action:             auth.ActionRoleReload,
		Scope:              auth.ScopeOrganizationUser,
		Type:               ut,
		OrganizationId:     f.orgId,
		OrganizationUserId: f.orgUserId,
		RemovedPolicies: [][]string{{
			fmt.Sprintf("Role:%d", f.roleId), string(ut), fmt.Sprintf("Resource:%d", f.resId), "allow",
		}},
	}
	require.NoError(t, env.SvcCtx.RedisCli.GetClient().
		Publish(context.Background(), env.SvcCtx.Config.CasbinConf.Watcher.Channel, msg.String()).Err())

	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		resp, err = NewAuthCheckLogic(ctx, env.SvcCtx).AuthCheck(authzCheckReq(ut))
		require.NoError(t, err)
		if !resp.GetData().GetAuthorized() {
			break
		}
		time.Sleep(20 * time.Millisecond)
	}
	assert.False(t, resp.GetData().GetAuthorized(), "role revocation should flip the next decision")
}

// T019 US-4: an unauthenticated call is rejected with gRPC Unauthenticated
// (the caller middleware maps it to 401).
func TestAuthzCheck_UnauthenticatedRejected(t *testing.T) {
	env := newAuthTestEnv(t)
	ut := constants.UserTypeBorrower
	seedAuthzGranted(t, env, ut, "Staff")

	_, err := NewAuthCheckLogic(context.Background(), env.SvcCtx).AuthCheck(authzCheckReq(ut))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not login")
}
