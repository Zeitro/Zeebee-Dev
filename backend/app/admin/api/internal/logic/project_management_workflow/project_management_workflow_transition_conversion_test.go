package project_management_workflow

import (
	"testing"

	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
)

func TestProjectManagementWorkflowUpdateRequestPreservesTransitionID(t *testing.T) {
	req := types.ProjectManagementWorkflowUpdateRequest{
		Transitions: []types.ProjectManagementWorkflowTransitionInput{{
			Id:          99,
			Name:        "start",
			FromStageId: 60,
			ToStageId:   69,
		}},
	}

	data, err := copierx.Trans[types.ProjectManagementWorkflowUpdateRequest, admin.ProjectManagementWorkflowUpdateRequestData](&req)
	if err != nil {
		t.Fatal(err)
	}
	if len(data.Transitions) != 1 {
		t.Fatalf("transition count = %d, want 1", len(data.Transitions))
	}
	if got := data.Transitions[0].Id; got != 99 {
		t.Fatalf("transition ID = %d, want 99", got)
	}
}
