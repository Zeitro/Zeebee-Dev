package logic

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

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

func TestProjectManagementTemplateOperationMetadata(t *testing.T) {
	copyTasks := true
	metadata, err := marshalProjectManagementOperationMetadata(
		projectManagementOperationMetadata{ProjectID: 1, TemplateID: 2, TemplateType: "Organization", CopyTasks: &copyTasks, SkippedAssignees: []map[string]any{{"organization_user_id": int64(9), "reason": "not_member"}}},
		nil, nil,
	)
	require.NoError(t, err)
	var actual map[string]any
	require.NoError(t, json.Unmarshal(metadata, &actual))
	require.Equal(t, float64(2), actual["template_id"])
	require.Equal(t, "Organization", actual["template_type"])
	require.Equal(t, true, actual["copy_tasks"])
	require.Equal(t, []any{map[string]any{"organization_user_id": float64(9), "reason": "not_member"}}, actual["skipped_assignees"])
}

func TestOperationValuesDiff(t *testing.T) {
	before, after := operationValuesDiff(
		map[string]any{"name": "before", "unchanged": []string{"a"}},
		map[string]any{"name": "after", "unchanged": []string{"a"}},
	)
	require.Equal(t, map[string]any{"name": "before"}, before)
	require.Equal(t, map[string]any{"name": "after"}, after)
}

func TestTaskOperationIdentities(t *testing.T) {
	require.Equal(t, map[string]any{"id": int64(1), "file_name": "contract.pdf"}, attachmentOperationIdentity(&model.File{Base: model.Base{Id: 1}, FileName: "contract.pdf"}))
	require.Equal(t, map[string]any{"id": int64(2), "name": "Prepare contract"}, dependencyOperationIdentity(&model.ProjectManagementTaskDependency{Base: model.Base{Id: 2}}, &model.ProjectManagementTask{Title: "Prepare contract"}))
}
