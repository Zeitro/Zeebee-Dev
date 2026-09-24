package project_management_workflow

import (
	"testing"

	"github.com/stretchr/testify/require"

	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
)

func TestWorkflowCreateConversionKeepsEntityScope(t *testing.T) {
	req := &types.ProjectManagementWorkflowCreateRequest{
		Name:       "Project workflow",
		EntityType: "Project",
		EntityId:   42,
	}
	data, err := copierx.Trans[types.ProjectManagementWorkflowCreateRequest, borrower.ProjectManagementWorkflowCreateRequestData](req)
	require.NoError(t, err)
	require.Equal(t, "Project workflow", data.Name)
	require.Equal(t, "Project", data.EntityType)
	require.Equal(t, int64(42), data.EntityId)
}
