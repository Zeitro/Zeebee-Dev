package project_management

import (
	"testing"

	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
)

func TestTaskTransitionRequestConverter(t *testing.T) {
	request := &types.ProjectManagementTaskTransitionRequest{Id: 12, TargetStageId: 34}
	data, err := copierx.Trans[types.ProjectManagementTaskTransitionRequest, borrower.ProjectManagementTaskTransitionRequestData](request)
	if err != nil {
		t.Fatal(err)
	}
	data.TaskId = request.Id
	if data.TaskId != 12 || data.TargetStageId != 34 {
		t.Fatalf("transition conversion lost data: %#v", data)
	}
}
