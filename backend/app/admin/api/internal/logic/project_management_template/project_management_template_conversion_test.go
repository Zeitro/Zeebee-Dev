package project_management_template

import (
	"testing"

	"github.com/stretchr/testify/require"

	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
)

func TestProjectManagementTemplateScheduleRequestConversion(t *testing.T) {
	create, err := copierx.Trans[types.ProjectManagementTemplateCreateRequest, admin.ProjectManagementTemplateCreateRequestData](&types.ProjectManagementTemplateCreateRequest{
		Name: "Template", SortOrder: 2, ProjectType: "Planned", StartAt: "01/31/2026", EndAt: "02/28/2026", PeriodValue: 1, PeriodUnit: "Month",
	})
	require.NoError(t, err)
	require.Equal(t, "02/28/2026", create.EndAt)
	require.NotNil(t, create.SortOrder)

	update, err := copierx.Trans[types.ProjectManagementTemplateUpdateRequest, admin.ProjectManagementTemplateUpdateRequestData](&types.ProjectManagementTemplateUpdateRequest{
		Id: 1, SortOrder: 3, ProjectType: "Planned", StartAt: "01/31/2026", EndAt: "02/28/2026", PeriodValue: 1, PeriodUnit: "Month",
	})
	require.NoError(t, err)
	require.Equal(t, "02/28/2026", update.EndAt)
	require.NotNil(t, update.SortOrder)
}
