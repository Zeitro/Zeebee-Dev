package projectmanagementlogic

import (
	"testing"
	"time"

	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"
)

func TestWorkflowConvertersKeepDefinitionSnapshots(t *testing.T) {
	now := time.Date(2026, 8, 19, 12, 0, 0, 0, time.UTC)
	updatedAt := now.Add(time.Hour)
	data, err := workflowDto("Template", "", true, true,
		[]*borrower.ProjectManagementWorkflowStageEntity{{Id: 2, Name: "Pending", Phase: "Pending", SortOrder: 10}},
		[]*borrower.ProjectManagementWorkflowTransitionEntity{{Id: 3, FromStageId: 2, ToStageId: 4, AllowedActorTypes: []string{"ProjectAdmin", "TaskCreator"}}},
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(data.Stages) != 1 || data.Stages[0].Phase != constants.ProjectManagementStagePhasePending || len(data.Transitions) != 1 || data.Transitions[0].AllowedActorTypes[1] != constants.ProjectManagementMemberRoleTaskCreator {
		t.Fatalf("workflow definition lost during request conversion: %#v", data)
	}

	entity, err := workflowEntity(&model.ProjectManagementWorkflow{
		Base: model.Base{Id: 1, CreatedAt: now, UpdatedAt: updatedAt}, Name: "Project workflow", SourceTemplateId: 9, SourceTemplateName: "Template",
		Stages:      []*model.ProjectManagementWorkflowStage{{Base: model.Base{Id: 2, CreatedAt: now, UpdatedAt: updatedAt}, Name: "Pending", Phase: constants.ProjectManagementStagePhasePending, SortOrder: 10}},
		Transitions: []*model.ProjectManagementWorkflowTransition{{Base: model.Base{Id: 3, CreatedAt: now, UpdatedAt: updatedAt}, FromStageId: 2, ToStageId: 4, AllowedActorTypes: []constants.ProjectManagementMemberRole{constants.ProjectManagementMemberRoleProjectAdmin, constants.ProjectManagementMemberRoleTaskCreator}}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if entity.SourceTemplateId != 9 || entity.SourceTemplateName != "Template" || entity.CreatedAt != now.Format(time.RFC3339) || entity.UpdatedAt != updatedAt.Format(time.RFC3339) || len(entity.Stages) != 1 || entity.Stages[0].Phase != "Pending" || entity.Stages[0].CreatedAt != now.Format(time.RFC3339) || len(entity.Transitions) != 1 || entity.Transitions[0].AllowedActorTypes[1] != "TaskCreator" || entity.Transitions[0].UpdatedAt != updatedAt.Format(time.RFC3339) {
		t.Fatalf("workflow definition lost during response conversion: %#v", entity)
	}
}
