package projectmanagementprojectlogic

import (
	"testing"

	"github.com/stretchr/testify/require"

	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
)

func TestProjectWorkflowMaintenanceRequestKeepsScopedIDsAndTransitions(t *testing.T) {
	req := &borrower.ProjectManagementProjectWorkflowStageUpdateRequestData{
		ProjectId:  42,
		WorkflowId: 84,
		Stages: []*borrower.ProjectManagementWorkflowUpdateStage{{
			Id: 100, Name: "Pending", Phase: string(constants.ProjectManagementStagePhasePending), SortOrder: 1,
			Transitions: []*borrower.ProjectManagementWorkflowUpdateTransition{{FromStageId: 100, ToStageId: 101, Name: "Start"}},
		}},
	}
	stages, err := copierx.Trans[[]*borrower.ProjectManagementWorkflowUpdateStage, []dto.ProjectManagementWorkflowUpdateStageDto](&req.Stages)
	require.NoError(t, err)
	require.EqualValues(t, 42, req.ProjectId)
	require.EqualValues(t, 84, req.WorkflowId)
	require.Len(t, *stages, 1)
	require.EqualValues(t, 100, (*stages)[0].Id)
	require.EqualValues(t, 101, (*stages)[0].Transitions[0].ToStageId)
}
