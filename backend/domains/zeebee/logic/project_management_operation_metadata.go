package logic

import (
	"encoding/json"
	"reflect"
)

type projectManagementOperationMetadata struct {
	ProjectID                int64          `json:"project_id,omitempty"`
	TaskID                   int64          `json:"task_id,omitempty"`
	TaskTypeID               int64          `json:"task_type_id,omitempty"`
	WorkflowID               int64          `json:"workflow_id,omitempty"`
	StageID                  int64          `json:"stage_id,omitempty"`
	TargetStageID            int64          `json:"target_stage_id,omitempty"`
	TransitionID             int64          `json:"transition_id,omitempty"`
	FileID                   int64          `json:"file_id,omitempty"`
	DependencyID             int64          `json:"dependency_id,omitempty"`
	OrganizationUserID       int64          `json:"organization_user_id,omitempty"`
	OrganizationUserFullName string         `json:"organization_user_full_name,omitempty"`
	Before                   map[string]any `json:"before,omitempty"`
	After                    map[string]any `json:"after,omitempty"`
}

func marshalProjectManagementOperationMetadata(metadata projectManagementOperationMetadata, before, after map[string]any) ([]byte, error) {
	if len(before) != 0 {
		metadata.Before = before
	}
	if len(after) != 0 {
		metadata.After = after
	}
	return json.Marshal(metadata)
}

func operationValuesDiff(before, after map[string]any) (map[string]any, map[string]any) {
	changedBefore, changedAfter := make(map[string]any), make(map[string]any)
	// Operation value pairs are built by the same helper and therefore share keys.
	for key, beforeValue := range before {
		if afterValue := after[key]; !reflect.DeepEqual(beforeValue, afterValue) {
			changedBefore[key], changedAfter[key] = beforeValue, afterValue
		}
	}
	return changedBefore, changedAfter
}
