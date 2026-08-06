package organizationemaillogic

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

// T063 US-17: SaveEmail persists an active organization's sender config.
func TestSaveEmail(t *testing.T) {
	env := newTestEnv(t)
	org := testutil.SeedOrganization(t, env.DB, &model.Organization{
		Name: "Org", Type: constants.UserTypeLender, Status: constants.OrganizationStatusActive,
	})

	_, err := NewSaveEmailLogic(context.Background(), env.SvcCtx).SaveEmail(&tenant.OrganizationEmailSaveEmailRequest{Data: &tenant.OrganizationEmailSaveEmailRequestData{
		UserType: string(constants.UserTypeLender), OrganizationId: org.Id,
		Email: "sender@test.local", Purpose: string(constants.EmailPurposeContact),
		AuthStrategy: string(constants.EmailAuthStrategyOrganization),
	}})
	require.NoError(t, err)

	var got model.OrganizationEmail
	require.NoError(t, env.DB.Where("organization_id = ?", org.Id).First(&got).Error)
	assert.Equal(t, "sender@test.local", got.Email)
}

// T063 US-17: GenerateAuthCodeUrl is explicitly deferred (SMTP-first scope).
func TestGenerateAuthCodeUrlDeferred(t *testing.T) {
	env := newTestEnv(t)
	_, err := NewGenerateAuthCodeUrlLogic(context.Background(), env.SvcCtx).GenerateAuthCodeUrl(&tenant.OrganizationEmailGenerateAuthCodeUrlRequest{Data: &tenant.OrganizationEmailGenerateAuthCodeUrlRequestData{
		UserType: string(constants.UserTypeLender), EmailProvider: "gmail",
	}})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "oauth not configured")
}
