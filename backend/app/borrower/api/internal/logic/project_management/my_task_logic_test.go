package project_management

import (
	"testing"

	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
)

func TestMyTaskFilterCarriesProjectAndOverdue(t *testing.T) {
	overdue := true
	request := &types.ProjectManagementMyTaskPageRequest{
		Page:   1,
		Limit:  20,
		Filter: types.ProjectManagementMyTaskPageFilter{ProjectId: 2, StageId: 3, TaskTypeId: 4, Keyword: "roof", Overdue: &overdue},
	}
	data, err := copierx.Trans[types.ProjectManagementMyTaskPageRequest, borrower.ProjectManagementMyTaskPageRequestData](request)
	if err != nil {
		t.Fatal(err)
	}
	if data.Filter == nil || data.Filter.ProjectId != 2 || data.Filter.StageId != 3 || data.Filter.TaskTypeId != 4 || data.Filter.Keyword != "roof" || data.Filter.Overdue == nil || !*data.Filter.Overdue {
		t.Fatalf("my task filters lost: %#v", data.Filter)
	}
}

func TestOperationTimelineResponseKeepsTaskContext(t *testing.T) {
	data := &borrower.ProjectManagementOperationPageResponseData{
		Page:      1,
		Limit:     20,
		Total:     1,
		TotalPage: 1,
		List: []*borrower.ProjectManagementOperationEntity{{
			Id: 1, ProjectId: 2, TaskId: 3, OperatorId: 4, OperatorName: "Operator", Operation: "task updated", CreatedAt: "2026-08-18T12:00:00Z",
		}},
	}
	response, err := copierx.Trans[borrower.ProjectManagementOperationPageResponseData, types.ProjectManagementOperationPageResponse](data)
	if err != nil {
		t.Fatal(err)
	}
	if response.Total != 1 || len(response.List) != 1 || response.List[0].ProjectId != 2 || response.List[0].TaskId != 3 || response.List[0].OperatorName != "Operator" {
		t.Fatalf("operation timeline context lost: %#v", response)
	}
}
