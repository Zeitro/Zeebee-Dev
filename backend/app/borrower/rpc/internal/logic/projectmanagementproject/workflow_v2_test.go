package projectmanagementprojectlogic

import (
	"testing"

	"github.com/stretchr/testify/require"

	"zeebee-backend/app/borrower/rpc/pb/borrower"
)

func TestProjectWorkflowPageRequestCarriesProjectFilter(t *testing.T) {
	req := &borrower.ProjectManagementProjectWorkflowPageRequest{Data: &borrower.ProjectManagementProjectWorkflowPageRequestData{
		Page: 2, Limit: 20, Filter: &borrower.ProjectManagementProjectWorkflowPageFilter{ProjectId: 42},
	}}
	require.EqualValues(t, 42, req.Data.Filter.ProjectId)
}
