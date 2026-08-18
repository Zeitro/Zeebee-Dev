package projectmanagementlogic

import (
	"testing"

	"zeebee-backend/app/borrower/rpc/pb/borrower"
)

func TestTaskTransitionRequestCarriesServerControlledTaskID(t *testing.T) {
	request := &borrower.ProjectManagementTaskTransitionRequest{Data: &borrower.ProjectManagementTaskTransitionRequestData{TaskId: 12, TargetStageId: 34}}
	if request.Data.TaskId != 12 || request.Data.TargetStageId != 34 {
		t.Fatalf("transition request lost task context: %#v", request.Data)
	}
}
