package organizationemaillogic

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"

	"zeebee-backend/app/tenant/config"
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

func TestGenerateAuthCodeUrl(t *testing.T) {
	env := newTestEnv(t)
	env.SvcCtx.Config.PlatformConf.Lender.EmailConf.Gmail = config.GmailConf{
		ClientId: "test-client", ClientSecret: "test-secret", CallbackUrl: "http://test.local/gmail",
	}
	org := testutil.SeedOrganization(t, env.DB, &model.Organization{
		Name: "Org", Type: constants.UserTypeLender, Status: constants.OrganizationStatusActive,
	})
	_, err := NewSaveEmailLogic(context.Background(), env.SvcCtx).SaveEmail(&tenant.OrganizationEmailSaveEmailRequest{Data: &tenant.OrganizationEmailSaveEmailRequestData{
		UserType: string(constants.UserTypeLender), OrganizationId: org.Id,
		Email: "sender@test.local", Purpose: string(constants.EmailPurposeContact),
		AuthStrategy: string(constants.EmailAuthStrategyOrganization),
	}})
	require.NoError(t, err)

	resp, err := NewGenerateAuthCodeUrlLogic(context.Background(), env.SvcCtx).GenerateAuthCodeUrl(&tenant.OrganizationEmailGenerateAuthCodeUrlRequest{Data: &tenant.OrganizationEmailGenerateAuthCodeUrlRequestData{
		UserType: string(constants.UserTypeLender), OrganizationId: org.Id,
		EmailProvider: "gmail", Purpose: string(constants.EmailPurposeContact),
	}})
	require.NoError(t, err)
	assert.Contains(t, resp.Data.Url, "client_id=test-client")
}

func TestGetAuthorization(t *testing.T) {
	env := newTestEnv(t)
	org := testutil.SeedOrganization(t, env.DB, &model.Organization{Name: "Org", Type: constants.UserTypeBorrower, Status: constants.OrganizationStatusActive})
	_, err := NewSaveEmailLogic(context.Background(), env.SvcCtx).SaveEmail(&tenant.OrganizationEmailSaveEmailRequest{Data: &tenant.OrganizationEmailSaveEmailRequestData{
		UserType: string(constants.UserTypeBorrower), OrganizationId: org.Id, Email: "sender@test.local", Purpose: string(constants.EmailPurposeContact), AuthStrategy: string(constants.EmailAuthStrategyOrganization),
	}})
	require.NoError(t, err)
	var organizationEmail model.OrganizationEmail
	require.NoError(t, env.DB.Where("organization_id = ?", org.Id).First(&organizationEmail).Error)
	require.NoError(t, env.DB.Create(&model.EmailAuthorization{
		OrganizationEmailId: organizationEmail.Id, Email: "sender@test.local", EmailProvider: constants.EmailProviderGmail, IsAuthorized: true, Purpose: constants.EmailPurposeContact,
		Token: &model.EmailToken{Token: &oauth2.Token{AccessToken: "access", RefreshToken: "refresh", Expiry: time.Now().UTC()}},
	}).Error)

	resp, err := NewGetAuthorizationLogic(context.Background(), env.SvcCtx).GetAuthorization(&tenant.OrganizationEmailGetAuthorizationRequest{Data: &tenant.OrganizationEmailGetAuthorizationRequestData{
		OrganizationId: org.Id, Purpose: string(constants.EmailPurposeContact),
	}})
	require.NoError(t, err)
	assert.Equal(t, "access", resp.Data.AccessToken)
	assert.True(t, resp.Data.IsAuthorized)
}
