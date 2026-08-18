package project_management

import (
	"testing"

	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
)

func TestWorkflowConvertersKeepDefinitionSnapshots(t *testing.T) {
	request := &types.ProjectManagementWorkflowCreateRequest{
		Name: "Template", Enabled: true,
		Stages:      []*types.ProjectManagementWorkflowStageEntity{{Name: "Pending", Phase: "Pending", SortOrder: 10}},
		Transitions: []*types.ProjectManagementWorkflowTransitionEntity{{FromStageId: 1, ToStageId: 2, AllowedActorTypes: []string{"ProjectOwner", "TaskExecutor"}}},
	}
	data, err := copierx.Trans[types.ProjectManagementWorkflowCreateRequest, borrower.ProjectManagementWorkflowCreateRequestData](request)
	if err != nil {
		t.Fatal(err)
	}
	if len(data.Stages) != 1 || data.Stages[0].Phase != "Pending" || len(data.Transitions) != 1 || data.Transitions[0].AllowedActorTypes[1] != "TaskExecutor" {
		t.Fatalf("workflow request lost definition: %#v", data)
	}

	response, err := copierx.Trans[borrower.ProjectManagementWorkflowEntity, types.ProjectManagementWorkflowEntity](&borrower.ProjectManagementWorkflowEntity{
		Id: 1, SourceTemplateId: 9, SourceTemplateName: "Template",
		Stages:      []*borrower.ProjectManagementWorkflowStageEntity{{Id: 2, Name: "Pending", Phase: "Pending", SortOrder: 10}},
		Transitions: []*borrower.ProjectManagementWorkflowTransitionEntity{{Id: 3, FromStageId: 2, ToStageId: 4, AllowedActorTypes: []string{"ProjectAdmin", "TaskCreator"}}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if response.SourceTemplateId != 9 || response.SourceTemplateName != "Template" || len(response.Stages) != 1 || response.Stages[0].Phase != "Pending" || len(response.Transitions) != 1 || response.Transitions[0].AllowedActorTypes[1] != "TaskCreator" {
		t.Fatalf("workflow response lost definition: %#v", response)
	}
}
