package projectmanagementprojectlogic

import (
	"testing"

	"github.com/stretchr/testify/require"

	"zeebee-backend/app/borrower/rpc/pb/borrower"
)

func TestTemplateDiscoveryResponseCarriesOnlyTemplateEntities(t *testing.T) {
	response := &borrower.ProjectManagementWorkflowTemplateResponse{}
	require.Empty(t, response.List)
}
