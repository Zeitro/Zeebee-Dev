package projectmanagementtasklogic

import (
	"testing"

	"github.com/stretchr/testify/require"

	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/data/dto"
)

func TestTaskUpdateCarriesWorkflowID(t *testing.T) {
	var input dto.ProjectManagementTaskUpdateDto
	require.NoError(t, copierx.Assign(&borrower.ProjectManagementTaskUpdateRequestData{Id: 9, WorkflowId: 33, Title: "Task"}, &input))
	require.EqualValues(t, 33, input.WorkflowId)
}
