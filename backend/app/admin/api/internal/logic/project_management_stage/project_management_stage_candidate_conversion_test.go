package project_management_stage

import (
	"testing"

	"github.com/stretchr/testify/require"

	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
)

func TestProjectManagementStageCandidateRequestConversion(t *testing.T) {
	data, err := copierx.Trans[types.ProjectManagementStageCandidateRequest, admin.ProjectManagementStageCandidateRequestData](&types.ProjectManagementStageCandidateRequest{TemplateId: 1, WorkflowId: 2, WorkflowType: "Bound"})
	require.NoError(t, err)
	require.Equal(t, int64(1), data.TemplateId)
	require.Equal(t, int64(2), data.WorkflowId)
	require.Equal(t, "Bound", data.WorkflowType)
}
