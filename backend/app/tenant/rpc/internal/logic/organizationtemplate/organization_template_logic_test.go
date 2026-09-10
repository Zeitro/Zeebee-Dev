package organizationtemplatelogic

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

// T067 US-19: template detail / detail-by-purpose / page / variable group edit /
// reset round-trip through the tenant RPC logic.
func TestTemplateLifecycle(t *testing.T) {
	env := newTestEnv(t)
	org := testutil.SeedOrganization(t, env.DB, &model.Organization{
		Name: "Org", Type: constants.UserTypeLender, Status: constants.OrganizationStatusActive,
	})
	base := &model.BaseTemplate{
		Name:           "Base",
		Type:           constants.UserTypeLender,
		Category:       constants.TemplateCategoryEmail,
		Purpose:        constants.OrganizationPurposePreapprovalLetter,
		Content:        "Hello {{name}}",
		Status:         constants.CommonStatusActive,
		Scope:          constants.BaseTemplateScopeSystem,
		IsConfigurable: true,
	}
	require.NoError(t, env.DB.Create(base).Error)
	tpl := &model.OrganizationTemplate{
		OrganizationId: org.Id,
		BaseTemplateId: base.Id,
		Name:           "Org Tpl",
		Type:           constants.UserTypeLender,
		Category:       constants.TemplateCategoryEmail,
		Purpose:        constants.OrganizationPurposePreapprovalLetter,
		Content:        "Hello {{name}}",
		Status:         constants.CommonStatusActive,
		HasConfigured:  false,
	}
	require.NoError(t, env.DB.Create(tpl).Error)

	detailResp, err := NewDetailLogic(context.Background(), env.SvcCtx).Detail(&tenant.OrganizationTemplateDetailRequest{Data: &tenant.OrganizationTemplateDetailRequestData{
		UserType: string(constants.UserTypeLender), Id: tpl.Id, OrganizationId: org.Id,
	}})
	require.NoError(t, err)
	assert.Equal(t, tpl.Id, detailResp.GetData().GetId())

	_, err = NewDetailLogic(context.Background(), env.SvcCtx).Detail(&tenant.OrganizationTemplateDetailRequest{Data: &tenant.OrganizationTemplateDetailRequestData{
		UserType: string(constants.UserTypeLender), Id: tpl.Id, OrganizationId: org.Id + 1,
	}})
	require.Error(t, err)

	byPurposeResp, err := NewDetailByPurposeLogic(context.Background(), env.SvcCtx).DetailByPurpose(&tenant.OrganizationTemplateDetailByPurposeRequest{Data: &tenant.OrganizationTemplateDetailByPurposeRequestData{
		UserType: string(constants.UserTypeLender), OrganizationId: org.Id,
		Purpose: constants.OrganizationPurposePreapprovalLetter,
	}})
	require.NoError(t, err)
	assert.Equal(t, tpl.Id, byPurposeResp.GetData().GetId())

	// Page filters by organization (required); filter/sorter may be omitted.
	pageResp, err := NewPageLogic(context.Background(), env.SvcCtx).Page(&tenant.OrganizationTemplatePageRequest{Data: &tenant.OrganizationTemplatePageRequestData{
		UserType: string(constants.UserTypeLender), Page: 1, Limit: 10,
		Filter: &tenant.OrganizationTemplatePageRequestFilter{OrganizationId: org.Id},
	}})
	require.NoError(t, err)
	assert.Equal(t, int64(1), pageResp.GetData().GetTotal())

	updateResp, err := NewUpdateVariableGroupLogic(context.Background(), env.SvcCtx).UpdateVariableGroup(&tenant.OrganizationTemplateUpdateVariableGroupRequest{Data: &tenant.OrganizationTemplateUpdateVariableGroupRequestData{
		UserType: string(constants.UserTypeLender), Id: tpl.Id, OrganizationId: org.Id,
	}})
	require.NoError(t, err)
	assert.Equal(t, tpl.Id, updateResp.GetData().GetId())

	resetResp, err := NewResetLogic(context.Background(), env.SvcCtx).Reset(&tenant.OrganizationTemplateResetRequest{Data: &tenant.OrganizationTemplateResetRequestData{
		UserType: string(constants.UserTypeLender), Id: tpl.Id, OrganizationId: org.Id,
	}})
	require.NoError(t, err)
	assert.Equal(t, tpl.Id, resetResp.GetData().GetId())

	_, err = NewResetByBaseTemplateIdLogic(context.Background(), env.SvcCtx).ResetByBaseTemplateId(&tenant.OrganizationTemplateResetByBaseTemplateIdRequest{Data: &tenant.OrganizationTemplateResetByBaseTemplateIdRequestData{
		UserType: string(constants.UserTypeLender), BaseTemplateId: base.Id,
	}})
	require.NoError(t, err)
}
