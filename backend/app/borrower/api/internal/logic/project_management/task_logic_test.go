package project_management

import (
	"testing"

	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
)

func TestTaskRequestConverters(t *testing.T) {
	request := &types.ProjectManagementTaskCreateRequest{Id: 9, ParentTaskId: 8, Title: "Child", AssigneeIds: []int64{2, 3}}
	data, err := copierx.Trans[types.ProjectManagementTaskCreateRequest, borrower.ProjectManagementTaskCreateRequestData](request)
	if err != nil {
		t.Fatal(err)
	}
	data.ProjectId = request.Id
	if data.ProjectId != 9 || data.ParentTaskId != 8 || data.Title != "Child" || len(data.AssigneeIds) != 2 {
		t.Fatalf("task request conversion lost data: %#v", data)
	}
}
