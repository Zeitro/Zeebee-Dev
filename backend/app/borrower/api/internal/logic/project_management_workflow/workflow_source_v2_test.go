package project_management_workflow

import (
	"testing"

	"github.com/stretchr/testify/require"

	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
)

func TestWorkflowSourceCreateConversionKeepsDefinition(t *testing.T) {
	req := &types.ProjectManagementWorkflowCreateRequest{
		Name:   "Default",
		Stages: []*types.ProjectManagementWorkflowStageInput{{Id: 1, Name: "Pending", Phase: "Pending", SortOrder: 1}},
	}
	data, err := copierx.Trans[types.ProjectManagementWorkflowCreateRequest, borrower.ProjectManagementWorkflowCreateRequestData](req)
	require.NoError(t, err)
	require.Equal(t, "Default", data.Name)
	require.Len(t, data.Stages, 1)
	require.Equal(t, int64(1), data.Stages[0].Id)
}
