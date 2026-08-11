package organizationuserlogic

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"zeebee-backend/app/tenant/config"
	_ "zeebee-backend/app/tenant/rpc/internal/email_tpl"
	"zeebee-backend/app/tenant/rpc/internal/logic"
	"zeebee-backend/app/tenant/rpc/internal/logic/testutil"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/stringx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"
)

// Invitation marker values mirror the unexported domain vars in
// domains/zeebee/logic/organization_user_logic.go (sent="1", accepted="2").
// They are duplicated here only to seed redis into the state the production
// invite path produces; if the domain values change, these must follow.
const (
	inviteMarkerSent     = "1"
	inviteMarkerAccepted = "2"
)

func newTestEnv(t *testing.T) *testutil.Env {
	t.Helper()
	env := testutil.NewEnv(t)
	logic.ResetBaseLogicForTest()
	return env
}

// seedActiveOrganization builds an active organization owned by an active user
// who is itself bound as an active member. The owner binding is what the
// disable/delete guards check (organization.OwnerId == member.UserId).
func seedActiveOrganization(t *testing.T, env *testutil.Env, userType constants.UserType) (
	*model.Organization, *model.OrganizationUser,
) {
	t.Helper()
	owner := testutil.SeedUser(t, env.DB, &model.User{
		Type:   userType,
		Email:  "owner-" + string(userType) + "@test.local",
		Status: constants.UserStatusActive,
	})
	organization := testutil.SeedOrganization(t, env.DB, &model.Organization{
		Name:    "Org-" + string(userType),
		Status:  constants.OrganizationStatusActive,
		OwnerId: owner.Id,
		Type:    userType,
	})
	ownerMember := testutil.SeedOrganizationUser(t, env.DB, &model.OrganizationUser{
		OrganizationId: organization.Id,
		UserId:         owner.Id,
		Status:         constants.OrganizationUserStatusActive,
	})
	return organization, ownerMember
}

// seedMember adds a second, non-owner active member to an organization.
func seedMember(t *testing.T, env *testutil.Env, organization *model.Organization, email string) (
	*model.User, *model.OrganizationUser,
) {
	t.Helper()
	user := testutil.SeedUser(t, env.DB, &model.User{
		Type:   organization.Type,
		Email:  email,
		Status: constants.UserStatusActive,
	})
	member := testutil.SeedOrganizationUser(t, env.DB, &model.OrganizationUser{
		OrganizationId: organization.Id,
		UserId:         user.Id,
		Status:         constants.OrganizationUserStatusActive,
	})
	return user, member
}

// seedInviteCode writes the invitation code and "sent" marker to redis exactly
// as the production invite path does, so AcceptInvitation / CheckInvitationCode
// can validate against it. Uses the prefixed SvcCtx.RedisCli so keys match.
func seedInviteCode(t *testing.T, env *testutil.Env, userType constants.UserType, email string, organizationId int64, code string) {
	t.Helper()
	ctx := context.Background()
	codeKey := stringx.Format(constants.UserInviteCode, userType, email, organizationId)
	markerKey := stringx.Format(constants.UserInviteCodeMarker, userType, email, organizationId)
	require.NoError(t, env.SvcCtx.RedisCli.SetExpire(ctx, codeKey, code, 3600))
	require.NoError(t, env.SvcCtx.RedisCli.SetExpire(ctx, markerKey, inviteMarkerSent, 7200))
}

// inviteEmailConfig is a dummy config; the invite email send runs async and
// no-ops in tests (nil SMTP client -> recovered panic), but the code/marker
// are persisted to redis before the send, which is what the accept path needs.
func inviteEmailConfig() *tenant.OrganizationUserInviteEmailConfig {
	return &tenant.OrganizationUserInviteEmailConfig{
		HostEmail:         "noreply@test.local",
		TplName:           "invite",
		ExpireSeconds:     3600,
		SendPeriodSeconds: 60,
		Url:               "http://test.local/accept",
	}
}

// T026 US-7: inviting a new email creates an unverified user bound to the
// organization as a not-yet-activated member.
func TestInvite_CreatesPendingMember(t *testing.T) {
	env := newTestEnv(t)
	organization, _ := seedActiveOrganization(t, env, constants.UserTypeLender)

	l := NewInviteOrganizationUserLogic(context.Background(), env.SvcCtx)
	resp, err := l.InviteOrganizationUser(&tenant.OrganizationUserInviteRequest{Data: &tenant.OrganizationUserInviteRequestData{
		UserType:       string(constants.UserTypeLender),
		OrganizationId: organization.Id,
		Email:          "invitee@test.local",
		Operator:       "Owner",
		EmailConfig:    inviteEmailConfig(),
	}})
	require.NoError(t, err)
	require.NotNil(t, resp.GetData())
	assert.Greater(t, resp.GetData().GetId(), int64(0))

	var user model.User
	require.NoError(t, env.DB.Where("email = ?", "invitee@test.local").First(&user).Error)
	assert.Equal(t, constants.UserStatusUnverified, user.Status)

	var member model.OrganizationUser
	require.NoError(t, env.DB.Where("organization_id = ? AND user_id = ?", organization.Id, user.Id).First(&member).Error)
	assert.Equal(t, constants.OrganizationUserStatusNotActivated, member.Status)
}

// T026 US-7: inviting an email that is already a member is rejected.
func TestInvite_DuplicateMemberRejected(t *testing.T) {
	env := newTestEnv(t)
	organization, _ := seedActiveOrganization(t, env, constants.UserTypeLender)

	l := NewInviteOrganizationUserLogic(context.Background(), env.SvcCtx)
	req := &tenant.OrganizationUserInviteRequest{Data: &tenant.OrganizationUserInviteRequestData{
		UserType:       string(constants.UserTypeLender),
		OrganizationId: organization.Id,
		Email:          "dup@test.local",
		EmailConfig:    inviteEmailConfig(),
	}}
	// first invite creates the pending member
	_, err := l.InviteOrganizationUser(req)
	require.NoError(t, err)
	// second invite to the same email is rejected
	_, err = l.InviteOrganizationUser(req)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "already in the organization")
}

// T026 US-7: accepting a valid invitation activates the member and the user.
func TestAcceptInvitation_ActivatesMember(t *testing.T) {
	env := newTestEnv(t)
	organization, _ := seedActiveOrganization(t, env, constants.UserTypeLender)

	invitee := testutil.SeedUser(t, env.DB, &model.User{
		Type:   constants.UserTypeLender,
		Email:  "invitee-accept@test.local",
		Status: constants.UserStatusUnverified,
	})
	testutil.SeedOrganizationUser(t, env.DB, &model.OrganizationUser{
		OrganizationId: organization.Id,
		UserId:         invitee.Id,
		Status:         constants.OrganizationUserStatusNotActivated,
	})
	code := "1234567890"
	seedInviteCode(t, env, constants.UserTypeLender, invitee.Email, organization.Id, code)

	resp, err := NewAcceptInvitationLogic(context.Background(), env.SvcCtx).AcceptInvitation(&tenant.OrganizationUserAcceptInvitationRequest{Data: &tenant.OrganizationUserAcceptInvitationRequestData{
		UserType:       string(constants.UserTypeLender),
		OrganizationId: organization.Id,
		Code:           code,
		Email:          invitee.Email,
		FirstName:      "Invitee",
		LastName:       "Test",
	}})
	require.NoError(t, err)
	require.NotNil(t, resp.GetData())

	var member model.OrganizationUser
	require.NoError(t, env.DB.Where("organization_id = ? AND user_id = ?", organization.Id, invitee.Id).First(&member).Error)
	assert.Equal(t, constants.OrganizationUserStatusActive, member.Status)
	assert.Equal(t, "Invitee", member.FirstName)

	var user model.User
	require.NoError(t, env.DB.First(&user, invitee.Id).Error)
	assert.Equal(t, constants.UserStatusActive, user.Status)
}

// T026 US-7: CheckInvitationCode distinguishes valid / expired / accepted /
// never-invited from the redis marker state.
func TestCheckInvitationCode_States(t *testing.T) {
	const (
		email       = "check@test.local"
		orgID int64 = 101
	)
	ut := constants.UserTypeLender
	ctx := context.Background()

	t.Run("valid code", func(t *testing.T) {
		env := newTestEnv(t)
		seedInviteCode(t, env, ut, email, orgID, "1111111111")
		resp, err := NewCheckInvitationCodeLogic(ctx, env.SvcCtx).CheckInvitationCode(&tenant.OrganizationUserCheckInvitationCodeRequest{Data: &tenant.OrganizationUserCheckInvitationCodeRequestData{
			UserType: string(ut), OrganizationId: orgID, Code: "1111111111", Email: email,
		}})
		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.False(t, resp.GetIsExpired())
	})

	t.Run("expired (code gone, sent marker)", func(t *testing.T) {
		env := newTestEnv(t)
		markerKey := stringx.Format(constants.UserInviteCodeMarker, ut, email, orgID)
		require.NoError(t, env.SvcCtx.RedisCli.SetExpire(ctx, markerKey, inviteMarkerSent, 7200))
		resp, err := NewCheckInvitationCodeLogic(ctx, env.SvcCtx).CheckInvitationCode(&tenant.OrganizationUserCheckInvitationCodeRequest{Data: &tenant.OrganizationUserCheckInvitationCodeRequestData{
			UserType: string(ut), OrganizationId: orgID, Code: "9999999999", Email: email,
		}})
		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.True(t, resp.GetIsExpired())
	})

	t.Run("already accepted", func(t *testing.T) {
		env := newTestEnv(t)
		markerKey := stringx.Format(constants.UserInviteCodeMarker, ut, email, orgID)
		require.NoError(t, env.SvcCtx.RedisCli.SetExpire(ctx, markerKey, inviteMarkerAccepted, 7200))
		_, err := NewCheckInvitationCodeLogic(ctx, env.SvcCtx).CheckInvitationCode(&tenant.OrganizationUserCheckInvitationCodeRequest{Data: &tenant.OrganizationUserCheckInvitationCodeRequestData{
			UserType: string(ut), OrganizationId: orgID, Code: "9999999999", Email: email,
		}})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "has been accepted")
	})

	t.Run("never invited", func(t *testing.T) {
		env := newTestEnv(t)
		_, err := NewCheckInvitationCodeLogic(ctx, env.SvcCtx).CheckInvitationCode(&tenant.OrganizationUserCheckInvitationCodeRequest{Data: &tenant.OrganizationUserCheckInvitationCodeRequestData{
			UserType: string(ut), OrganizationId: orgID, Code: "9999999999", Email: email,
		}})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid invitation")
	})
}

// T026 US-7: AssignRoles binds the requested roles onto a member.
func TestAssignRoles_RebindsMemberRoles(t *testing.T) {
	env := newTestEnv(t)
	organization, _ := seedActiveOrganization(t, env, constants.UserTypeLender)
	_, member := seedMember(t, env, organization, "assign@test.local")
	roleA := testutil.SeedRole(t, env.DB, &model.Role{
		OrganizationId: organization.Id, Name: "A", Type: constants.UserTypeLender, Status: constants.CommonStatusActive,
	})
	roleB := testutil.SeedRole(t, env.DB, &model.Role{
		OrganizationId: organization.Id, Name: "B", Type: constants.UserTypeLender, Status: constants.CommonStatusActive,
	})

	_, err := NewAssignRolesLogic(context.Background(), env.SvcCtx).AssignRoles(&tenant.OrganizationUserAssignRolesRequest{Data: &tenant.OrganizationUserAssignRolesRequestData{
		UserType: string(constants.UserTypeLender), OrganizationId: organization.Id, Id: member.Id, RoleIds: []int64{roleA.Id, roleB.Id},
	}})
	require.NoError(t, err)

	var count int64
	require.NoError(t, env.DB.Table("organization_user_and_roles").
		Where("organization_user_id = ?", member.Id).Count(&count).Error)
	assert.Equal(t, int64(2), count)
}

// T026 US-7: SwitchStatus disables a non-owner member; the owner is protected.
func TestSwitchStatus_DisablesMemberAndProtectsOwner(t *testing.T) {
	env := newTestEnv(t)
	organization, ownerMember := seedActiveOrganization(t, env, constants.UserTypeLender)
	_, member := seedMember(t, env, organization, "disable@test.local")

	_, err := NewSwitchStatusLogic(context.Background(), env.SvcCtx).SwitchStatus(&tenant.OrganizationUserSwitchStatusRequest{Data: &tenant.OrganizationUserSwitchStatusRequestData{
		UserType: string(constants.UserTypeLender), OrganizationId: organization.Id, Id: member.Id, Status: string(constants.OrganizationUserStatusDisabled),
	}})
	require.NoError(t, err)

	var got model.OrganizationUser
	require.NoError(t, env.DB.First(&got, member.Id).Error)
	assert.Equal(t, constants.OrganizationUserStatusDisabled, got.Status)

	_, err = NewSwitchStatusLogic(context.Background(), env.SvcCtx).SwitchStatus(&tenant.OrganizationUserSwitchStatusRequest{Data: &tenant.OrganizationUserSwitchStatusRequestData{
		UserType: string(constants.UserTypeLender), OrganizationId: organization.Id, Id: ownerMember.Id, Status: string(constants.OrganizationUserStatusDisabled),
	}})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cannot disable owner")
}

// T026 US-7: Delete removes a non-owner member.
func TestDelete_RemovesMember(t *testing.T) {
	env := newTestEnv(t)
	organization, _ := seedActiveOrganization(t, env, constants.UserTypeLender)
	_, member := seedMember(t, env, organization, "delete@test.local")

	_, err := NewDeleteLogic(context.Background(), env.SvcCtx).Delete(&tenant.OrganizationUserDeleteRequest{Data: &tenant.OrganizationUserDeleteRequestData{
		UserType: string(constants.UserTypeLender), OrganizationId: organization.Id, Id: member.Id,
	}})
	require.NoError(t, err)

	var count int64
	require.NoError(t, env.DB.Model(&model.OrganizationUser{}).Where("id = ?", member.Id).Count(&count).Error)
	assert.Zero(t, count)
}

// T028 US-7: Page lists all members of an organization (any status).
func TestPage_ListsMembers(t *testing.T) {
	env := newTestEnv(t)
	organization, _ := seedActiveOrganization(t, env, constants.UserTypeLender)
	seedMember(t, env, organization, "page@test.local")

	resp, err := NewPageLogic(context.Background(), env.SvcCtx).Page(&tenant.OrganizationUserPageRequest{Data: &tenant.OrganizationUserPageRequestData{
		Page: 1, Limit: 10, UserType: string(constants.UserTypeLender),
		Filter: &tenant.OrganizationUserPageRequestFilter{OrganizationId: organization.Id},
	}})
	require.NoError(t, err)
	require.NotNil(t, resp.GetData())
	assert.Equal(t, int64(2), resp.GetData().GetTotal())
	assert.Len(t, resp.GetData().GetList(), 2)
}

// T069 US-7: AllOrganizationUsers returns active members only.
func TestAllOrganizationUsers_OnlyActive(t *testing.T) {
	env := newTestEnv(t)
	organization, _ := seedActiveOrganization(t, env, constants.UserTypeLender) // owner = active
	seedMember(t, env, organization, "all-active@test.local")                   // active
	// a not-yet-activated member is excluded from the active-only list
	pending := testutil.SeedUser(t, env.DB, &model.User{
		Type: constants.UserTypeLender, Email: "all-pending@test.local", Status: constants.UserStatusUnverified,
	})
	testutil.SeedOrganizationUser(t, env.DB, &model.OrganizationUser{
		OrganizationId: organization.Id, UserId: pending.Id, Status: constants.OrganizationUserStatusNotActivated,
	})

	resp, err := NewAllOrganizationUsersLogic(context.Background(), env.SvcCtx).AllOrganizationUsers(&tenant.OrganizationUserAllOrganizationUsersRequest{Data: &tenant.OrganizationUserAllOrganizationUsersRequestData{
		UserType: string(constants.UserTypeLender), OrganizationId: organization.Id,
	}})
	require.NoError(t, err)
	assert.Len(t, resp.GetData(), 2) // owner + one active member; pending excluded
}

// T069 US-7: Detail returns a member by id (org-scoped).
func TestDetail_ReturnsMember(t *testing.T) {
	env := newTestEnv(t)
	organization, _ := seedActiveOrganization(t, env, constants.UserTypeLender)
	_, member := seedMember(t, env, organization, "detail@test.local")

	resp, err := NewDetailLogic(context.Background(), env.SvcCtx).Detail(&tenant.OrganizationUserDetailRequest{Data: &tenant.OrganizationUserDetailRequestData{
		UserType: string(constants.UserTypeLender), OrganizationId: organization.Id, Id: member.Id,
	}})
	require.NoError(t, err)
	require.NotNil(t, resp.GetData())
	assert.Equal(t, member.Id, resp.GetData().GetId())
}

// T069 US-7: CreateAdminOrganizationUser provisions an active member with
// credentials and roles.
func TestCreateAdminOrganizationUser_CreatesActiveMember(t *testing.T) {
	env := newTestEnv(t)
	organization, _ := seedActiveOrganization(t, env, constants.UserTypeLender)
	role := testutil.SeedRole(t, env.DB, &model.Role{
		OrganizationId: organization.Id, Name: "Admin", Type: constants.UserTypeLender, Status: constants.CommonStatusActive,
	})

	resp, err := NewCreateAdminOrganizationUserLogic(context.Background(), env.SvcCtx).CreateAdminOrganizationUser(&tenant.OrganizationUserCreateAdminRequest{Data: &tenant.OrganizationUserCreateAdminRequestData{
		UserType: string(constants.UserTypeLender), OrganizationId: organization.Id,
		Email: "admin-create@test.local", FirstName: "Ad", LastName: "Min",
		Password: "Passw0rd2026", RoleIds: []int64{role.Id},
	}})
	require.NoError(t, err)
	require.NotNil(t, resp.GetData())

	var user model.User
	require.NoError(t, env.DB.Where("email = ?", "admin-create@test.local").First(&user).Error)
	assert.Equal(t, constants.UserStatusActive, user.Status)
	assert.NotEmpty(t, user.PasswordHash)

	var member model.OrganizationUser
	require.NoError(t, env.DB.Where("organization_id = ? AND user_id = ?", organization.Id, user.Id).First(&member).Error)
	assert.Equal(t, constants.OrganizationUserStatusActive, member.Status)

	var roleCount int64
	require.NoError(t, env.DB.Table("organization_user_and_roles").Where("organization_user_id = ?", member.Id).Count(&roleCount).Error)
	assert.Equal(t, int64(1), roleCount)
	_ = role
}

// T028 US-7: UpdateOrganizationUser updates a member's profile fields and roles.
func TestUpdateOrganizationUser_UpdatesFieldsAndRoles(t *testing.T) {
	env := newTestEnv(t)
	organization, _ := seedActiveOrganization(t, env, constants.UserTypeLender)
	_, member := seedMember(t, env, organization, "update@test.local")
	role := testutil.SeedRole(t, env.DB, &model.Role{
		OrganizationId: organization.Id, Name: "R", Type: constants.UserTypeLender, Status: constants.CommonStatusActive,
	})

	_, err := NewUpdateOrganizationUserLogic(context.Background(), env.SvcCtx).UpdateOrganizationUser(&tenant.OrganizationUserUpdateRequest{Data: &tenant.OrganizationUserUpdateRequestData{
		UserType: string(constants.UserTypeLender), OrganizationId: organization.Id, Id: member.Id,
		FirstName: "Updated", LastName: "Member", Remark: "rm", RoleIds: []int64{role.Id},
	}})
	require.NoError(t, err)

	var got model.OrganizationUser
	require.NoError(t, env.DB.First(&got, member.Id).Error)
	assert.Equal(t, "Updated", got.FirstName)
	assert.Equal(t, "rm", got.Remark)

	var roleCount int64
	require.NoError(t, env.DB.Table("organization_user_and_roles").Where("organization_user_id = ?", member.Id).Count(&roleCount).Error)
	assert.Equal(t, int64(1), roleCount)
}

// T069 US-7: UpdateProfile updates the requesting member's own name/phone.
func TestUpdateProfile_UpdatesOwnName(t *testing.T) {
	env := newTestEnv(t)
	organization, _ := seedActiveOrganization(t, env, constants.UserTypeLender)
	user, _ := seedMember(t, env, organization, "profile@test.local")

	_, err := NewUpdateProfileLogic(context.Background(), env.SvcCtx).UpdateProfile(&tenant.OrganizationUserUpdateProfileRequest{Data: &tenant.OrganizationUserUpdateProfileRequestData{
		UserType: string(constants.UserTypeLender), OrganizationId: organization.Id, UserId: user.Id,
		FirstName: "Self", LastName: "Updated",
	}})
	require.NoError(t, err)

	var got model.OrganizationUser
	require.NoError(t, env.DB.Where("organization_id = ? AND user_id = ?", organization.Id, user.Id).First(&got).Error)
	assert.Equal(t, "Self", got.FirstName)
	assert.Equal(t, "Updated", got.LastName)
}

// T069 US-7: SendInviteEmail rejects an already-active member (only pending
// members can be re-invited). The email send itself needs SMTP and is not
// exercised here; the guard is the testable, deterministic behavior.
func TestSendInviteEmail_RejectsActiveMember(t *testing.T) {
	env := newTestEnv(t)
	organization, _ := seedActiveOrganization(t, env, constants.UserTypeLender)
	_, member := seedMember(t, env, organization, "resend@test.local")

	_, err := NewSendInviteEmailLogic(context.Background(), env.SvcCtx).SendInviteEmail(&tenant.OrganizationUserSendInviteEmailRequest{Data: &tenant.OrganizationUserSendInviteEmailRequestData{
		UserType: string(constants.UserTypeLender), OrganizationId: organization.Id, Id: member.Id,
		Operator: "Owner", EmailConfig: inviteEmailConfig(),
	}})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not activated")
}

func TestSendInviteEmail_UsesPlatformEmailClient(t *testing.T) {
	env := newTestEnv(t)
	env.SvcCtx.Config.PlatformConf.Borrower.EmailConf = config.EmailConf{
		Host: "127.0.0.1", Port: 1, Username: "noreply@test.local",
	}
	organization, _ := seedActiveOrganization(t, env, constants.UserTypeBorrower)
	user := testutil.SeedUser(t, env.DB, &model.User{
		Type:   constants.UserTypeBorrower,
		Email:  "pending@test.local",
		Status: constants.UserStatusUnverified,
	})
	member := testutil.SeedOrganizationUser(t, env.DB, &model.OrganizationUser{
		OrganizationId: organization.Id,
		UserId:         user.Id,
		Status:         constants.OrganizationUserStatusNotActivated,
	})

	_, err := NewSendInviteEmailLogic(context.Background(), env.SvcCtx).SendInviteEmail(&tenant.OrganizationUserSendInviteEmailRequest{Data: &tenant.OrganizationUserSendInviteEmailRequestData{
		UserType: string(constants.UserTypeBorrower), OrganizationId: organization.Id, Id: member.Id,
		Operator: "Owner", EmailConfig: inviteEmailConfig(),
	}})
	require.Error(t, err)
	assert.NotContains(t, err.Error(), "email client is not initialized")
}

func TestSendInviteEmail_DefaultsToInviteTemplate(t *testing.T) {
	env := newTestEnv(t)
	env.SvcCtx.Config.PlatformConf.Borrower.EmailConf = config.EmailConf{
		Host: "127.0.0.1", Port: 1, Username: "noreply@test.local",
	}
	organization, _ := seedActiveOrganization(t, env, constants.UserTypeBorrower)
	user := testutil.SeedUser(t, env.DB, &model.User{
		Type:   constants.UserTypeBorrower,
		Email:  "pending-template@test.local",
		Status: constants.UserStatusUnverified,
	})
	member := testutil.SeedOrganizationUser(t, env.DB, &model.OrganizationUser{
		OrganizationId: organization.Id,
		UserId:         user.Id,
		Status:         constants.OrganizationUserStatusNotActivated,
	})

	_, err := NewSendInviteEmailLogic(context.Background(), env.SvcCtx).SendInviteEmail(&tenant.OrganizationUserSendInviteEmailRequest{Data: &tenant.OrganizationUserSendInviteEmailRequestData{
		UserType: string(constants.UserTypeBorrower), OrganizationId: organization.Id, Id: member.Id,
		Operator: "Owner", EmailConfig: &tenant.OrganizationUserInviteEmailConfig{
			HostEmail: "noreply@test.local", ExpireSeconds: 3600, SendPeriodSeconds: 60,
			Url: "http://test.local/invite",
		},
	}})
	require.Error(t, err)
	assert.NotContains(t, err.Error(), "Email template not found")
}

func TestToInviteEmailConfig_UsesPlatformDefaults(t *testing.T) {
	logo := "https://cdn.zeebee.com/logo.png"
	address := "Zeebee | 1 Market Street | San Francisco, CA, 94105"
	profile := config.PlatformProfile{
		EmailConf:                    config.EmailConf{Username: "noreply@test.local"},
		InviteUserUrl:                "http://test.local/invite",
		InviteUserExpireSeconds:      3600,
		SendInviteEmailPeriodSeconds: 60,
		LogoUrl:                      logo,
		FooterAddress:                address,
	}

	got := toInviteEmailConfig(nil, profile)
	assert.Equal(t, "noreply@test.local", got.HostEmail)
	assert.Equal(t, 3600, got.ExpireSeconds)
	assert.Equal(t, 60, got.SendPeriodSeconds)
	assert.Equal(t, "http://test.local/invite", got.Url)
	assert.Equal(t, logo, got.Logo)
	assert.Equal(t, address, got.Address)

	got = toInviteEmailConfig(&tenant.OrganizationUserInviteEmailConfig{
		HostEmail:         "client@test.local",
		ExpireSeconds:     7200,
		SendPeriodSeconds: 120,
		Url:               "http://client.test/invite",
		Logo:              "https://client.example/logo.png",
		Address:           "Client address",
	}, profile)
	assert.Equal(t, "client@test.local", got.HostEmail)
	assert.Equal(t, 7200, got.ExpireSeconds)
	assert.Equal(t, 120, got.SendPeriodSeconds)
	assert.Equal(t, "http://client.test/invite", got.Url)
	assert.Equal(t, "https://client.example/logo.png", got.Logo)
	assert.Equal(t, "Client address", got.Address)
}
