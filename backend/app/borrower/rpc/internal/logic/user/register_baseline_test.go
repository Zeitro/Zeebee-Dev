package userlogic

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"
)

// T012 characterization baseline: borrower register new-user path. Only the
// synchronous identity contract is pinned — user + organization + membership
// created in the states RegisterWithOrganization produces for a self-service
// (!isManual) borrower signup.
//
// The post-register goroutines (verify email, draft preapproval loan) are
// fire-and-forget business logic, not the identity contract under test. The
// new-user path is exercised WITHOUT PreApprovalRequestId so the preapproval
// goroutine early-returns; the verify-email goroutine reads only the redis +
// no-op email client the harness already wires.

// Register with a fresh email creates an unverified borrower user owning a new
// (not-yet-activated) organization, linked by a membership row.
func TestBaselineRegister_NewUser(t *testing.T) {
	env := newBaselineEnv(t)
	// mixed-case on purpose: Register lowercases the email before use.
	email := "NewBorrower@TEST.local"

	l := NewRegisterLogic(context.Background(), env.SvcCtx)
	resp, err := l.Register(&borrower.UserRegisterRequest{Data: &borrower.UserRegisterRequestData{
		Email:            email,
		Password:         baselinePassword,
		FirstName:        "New",
		LastName:         "Borrower",
		OrganizationName: "New Borrower Org",
	}})
	require.NoError(t, err)
	require.NotNil(t, resp)

	lower := strings.ToLower(email)

	// user created, unverified, borrower-typed, email lowercased.
	var user model.User
	require.NoError(t, env.DB.Where("email = ?", lower).First(&user).Error)
	assert.Equal(t, lower, user.Email)
	assert.Equal(t, constants.UserTypeBorrower, user.Type)
	assert.Equal(t, constants.UserStatusUnverified, user.Status)

	// organization created, owned by the user, not yet activated (!isManual).
	var org model.Organization
	require.NoError(t, env.DB.Where("id = ?", user.DefaultOrganizationId).First(&org).Error)
	assert.Equal(t, "New Borrower Org", org.Name)
	assert.Equal(t, constants.OrganizationStatusNotActivated, org.Status)
	assert.Equal(t, user.Id, org.OwnerId)

	// exactly one membership links the user to their default org.
	var count int64
	require.NoError(t, env.DB.Model(&model.OrganizationUser{}).
		Where("organization_id = ? AND user_id = ?", org.Id, user.Id).Count(&count).Error)
	assert.EqualValues(t, 1, count)
}
