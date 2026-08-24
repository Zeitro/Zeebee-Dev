package project_management_workflow

import (
	"testing"

	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
)

func TestWorkflowUpdateRequestConversion(t *testing.T) {
	req := &types.ProjectManagementWorkflowUpdateRequest{Stages: []*types.ProjectManagementWorkflowUpdateStage{{
		Id: 1, Name: "Pending", Phase: "Pending", Color: "grey", SortOrder: 10,
		Transitions: []*types.ProjectManagementWorkflowUpdateTransition{{Name: "Next", FromStageId: 1, ToStageId: 2, AllowedActorTypes: []string{"ProjectOwner"}}},
	}}}
	data, err := copierx.Trans[types.ProjectManagementWorkflowUpdateRequest, borrower.ProjectManagementWorkflowUpdateRequestData](req)
	if err != nil {
		t.Fatal(err)
	}
	if len(data.Stages) != 1 || len(data.Stages[0].Transitions) != 1 || data.Stages[0].Transitions[0].FromStageId != 1 || data.Stages[0].Transitions[0].ToStageId != 2 {
		t.Fatalf("unexpected converted workflow: %#v", data)
	}
}

func TestStageSwapOrderRequestConversion(t *testing.T) {
	req := &types.ProjectManagementWorkflowStageSwapOrderRequest{Id: 1, FromStageId: 2, ToStageId: 3}
	data, err := copierx.Trans[types.ProjectManagementWorkflowStageSwapOrderRequest, borrower.ProjectManagementWorkflowStageSwapOrderRequestData](req)
	if err != nil {
		t.Fatal(err)
	}
	if data.Id != 1 || data.FromStageId != 2 || data.ToStageId != 3 {
		t.Fatalf("unexpected converted stage swap request: %#v", data)
	}
}
