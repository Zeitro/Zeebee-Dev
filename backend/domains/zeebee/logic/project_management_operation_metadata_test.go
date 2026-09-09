package logic

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"
)

func TestProjectManagementOperationMetadata(t *testing.T) {
	metadata, err := marshalProjectManagementOperationMetadata(
		projectManagementOperationMetadata{ProjectID: 1, TaskID: 2, WorkflowID: 3, StageID: 4, TransitionID: 5},
		map[string]any{"assignee_id": int64(2), "assignee_name": "Zhang San"},
		nil,
	)
	require.NoError(t, err)

	var actual map[string]any
	require.NoError(t, json.Unmarshal(metadata, &actual))
	require.Equal(t, float64(1), actual["project_id"])
	require.Equal(t, float64(2), actual["task_id"])
	require.Equal(t, float64(3), actual["workflow_id"])
	require.Equal(t, float64(4), actual["stage_id"])
	require.Equal(t, float64(5), actual["transition_id"])
	require.Equal(t, map[string]any{"assignee_id": float64(2), "assignee_name": "Zhang San"}, actual["before"])
	require.NotContains(t, actual, "after")
}

func TestOperationValuesDiff(t *testing.T) {
	before, after := operationValuesDiff(
		map[string]any{"name": "before", "unchanged": []string{"a"}},
		map[string]any{"name": "after", "unchanged": []string{"a"}},
	)
	require.Equal(t, map[string]any{"name": "before"}, before)
	require.Equal(t, map[string]any{"name": "after"}, after)
}

func TestWorkflowOperationIdentities(t *testing.T) {
	require.Equal(t, map[string]any{"id": int64(1), "name": "Pending"}, workflowStageOperationIdentity(&model.ProjectManagementWorkflowStage{Base: model.Base{Id: 1}, Name: "Pending"}))
	require.Equal(t, map[string]any{"id": int64(2), "name": "Start"}, workflowTransitionOperationIdentity(&model.ProjectManagementWorkflowTransition{Base: model.Base{Id: 2}, Name: "Start"}))
}

func TestTaskOperationIdentities(t *testing.T) {
	require.Equal(t, map[string]any{"id": int64(1), "file_name": "contract.pdf"}, attachmentOperationIdentity(&model.File{Base: model.Base{Id: 1}, FileName: "contract.pdf"}))
	require.Equal(t, map[string]any{"id": int64(2), "name": "Prepare contract"}, dependencyOperationIdentity(&model.ProjectManagementTaskDependency{Base: model.Base{Id: 2}}, &model.ProjectManagementTask{Title: "Prepare contract"}))
}

func TestWorkflowTransitionOperationValuesKeepsAllowedActorTypesSlice(t *testing.T) {
	values := workflowTransitionOperationValues(&model.ProjectManagementWorkflowTransition{AllowedActorTypes: []constants.ProjectManagementMemberRole{constants.ProjectManagementMemberRoleProjectOwner, constants.ProjectManagementMemberRoleTaskExecutor}})
	require.Equal(t, []constants.ProjectManagementMemberRole{constants.ProjectManagementMemberRoleProjectOwner, constants.ProjectManagementMemberRoleTaskExecutor}, values["allowed_actor_types"])
}

func TestWorkflowOperationValuesUseColumnNames(t *testing.T) {
	stage := &model.ProjectManagementWorkflowStage{Base: model.Base{Id: 1}, Name: "Pending", Description: "Waiting", Phase: constants.ProjectManagementStagePhasePending, Color: "grey", SortOrder: 10}
	require.Equal(t, map[string]any{"name": "Pending", "description": "Waiting", "phase": constants.ProjectManagementStagePhasePending, "color": "grey", "sort_order": 10}, workflowStageOperationValues(stage))

	transition := &model.ProjectManagementWorkflowTransition{Base: model.Base{Id: 2}, Name: "Start", FromStageId: 1, ToStageId: 3}
	require.Equal(t, map[string]any{"name": "Start", "from_stage_id": int64(1), "to_stage_id": int64(3), "allowed_actor_types": []constants.ProjectManagementMemberRole(nil)}, workflowTransitionOperationValues(transition))
}
