package organizationsettinglogic

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

// T065 US-18: DetailByPurpose returns the org's setting; SwitchStatus toggles
// and rejects another org's setting.
func TestDetailAndSwitchStatus(t *testing.T) {
	env := newTestEnv(t)
	org := testutil.SeedOrganization(t, env.DB, &model.Organization{
		Name: "Org", Type: constants.UserTypeLender, Status: constants.OrganizationStatusActive,
	})
	setting := &model.OrganizationSetting{
		OrganizationId: org.Id,
		Purpose:        constants.OrganizationPurposePreapprovalLetter,
		Enabled:        false,
		Setting:        map[string]interface{}{},
	}
	require.NoError(t, env.DB.Create(setting).Error)

	detailResp, err := NewDetailByPurposeLogic(context.Background(), env.SvcCtx).DetailByPurpose(&tenant.OrganizationSettingDetailByPurposeRequest{Data: &tenant.OrganizationSettingDetailByPurposeRequestData{
		UserType: string(constants.UserTypeLender), OrganizationId: org.Id,
		Purpose: constants.OrganizationPurposePreapprovalLetter,
	}})
	require.NoError(t, err)
	assert.Equal(t, setting.Id, detailResp.GetData().GetId())

	_, err = NewSwitchStatusLogic(context.Background(), env.SvcCtx).SwitchStatus(&tenant.OrganizationSettingSwitchStatusRequest{Data: &tenant.OrganizationSettingSwitchStatusRequestData{
		UserType: string(constants.UserTypeLender), Id: setting.Id, OrganizationId: org.Id, Enabled: false,
	}})
	require.NoError(t, err)

	// another org cannot switch this org's setting
	_, err = NewSwitchStatusLogic(context.Background(), env.SvcCtx).SwitchStatus(&tenant.OrganizationSettingSwitchStatusRequest{Data: &tenant.OrganizationSettingSwitchStatusRequestData{
		UserType: string(constants.UserTypeLender), Id: setting.Id, OrganizationId: org.Id + 999, Enabled: false,
	}})
	require.Error(t, err)
}
