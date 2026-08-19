package logic

import (
	"context"
	"fmt"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/repo"
)

func TestValidateProjectManagementWorkflow(t *testing.T) {
	valid := dto.ProjectManagementWorkflowDto{
		Name: "Template",
		Stages: []dto.ProjectManagementWorkflowStageDto{
			{Id: 1, Name: "Pending", Phase: constants.ProjectManagementStagePhasePending, SortOrder: 1},
			{Id: 2, Name: "Finished", Phase: constants.ProjectManagementStagePhaseFinished, SortOrder: 2},
		},
		Transitions: []dto.ProjectManagementWorkflowTransitionDto{{FromStageId: 1, ToStageId: 2}},
	}
	if err := ValidateProjectManagementWorkflow(valid); err != nil {
		t.Fatalf("valid workflow rejected: %v", err)
	}

	valid.Stages = valid.Stages[1:]
	if err := ValidateProjectManagementWorkflow(valid); err == nil {
		t.Fatal("workflow without pending stage accepted")
	}

	valid.Stages = []dto.ProjectManagementWorkflowStageDto{
		{Id: 1, Name: "Pending", Phase: constants.ProjectManagementStagePhasePending, SortOrder: 1},
		{Id: 2, Name: "Also pending", Phase: constants.ProjectManagementStagePhasePending, SortOrder: 2},
	}
	if err := ValidateProjectManagementWorkflow(valid); err == nil {
		t.Fatal("workflow with multiple pending stages accepted")
	}

	valid.Stages = []dto.ProjectManagementWorkflowStageDto{
		{Id: 1, Name: "Pending", Phase: constants.ProjectManagementStagePhasePending, SortOrder: 1},
		{Id: 2, Name: "Finished", Phase: constants.ProjectManagementStagePhaseFinished, SortOrder: 2},
	}
	valid.Transitions = []dto.ProjectManagementWorkflowTransitionDto{{FromStageId: 1, ToStageId: 3}}
	if err := ValidateProjectManagementWorkflow(valid); err == nil {
		t.Fatal("workflow transition to an unknown stage accepted")
	}
}

func TestNormalizeProjectManagementTemplateStageReferences(t *testing.T) {
	workflow := dto.ProjectManagementWorkflowDto{
		Stages: []dto.ProjectManagementWorkflowStageDto{
			{Name: "Pending", Phase: constants.ProjectManagementStagePhasePending, SortOrder: 1},
			{Id: -1, Name: "Finished", Phase: constants.ProjectManagementStagePhaseFinished, SortOrder: 2},
			{Name: "Archived", Phase: constants.ProjectManagementStagePhaseFinished, SortOrder: 3},
		},
	}

	normalizeProjectManagementTemplateStageReferences(&workflow)

	if workflow.Stages[0].Id != -2 || workflow.Stages[1].Id != -1 || workflow.Stages[2].Id != -3 {
		t.Fatalf("unexpected stage references: %#v", workflow.Stages)
	}
	workflow.Name = "Template"
	if err := ValidateProjectManagementWorkflow(workflow); err != nil {
		t.Fatalf("normalized workflow rejected: %v", err)
	}
}

func TestProjectManagementTemplateModels(t *testing.T) {
	input := dto.ProjectManagementWorkflowDto{
		Name: "Template",
		Stages: []dto.ProjectManagementWorkflowStageDto{
			{Id: -1, Name: "Pending", Phase: constants.ProjectManagementStagePhasePending, SortOrder: 1},
			{Id: -2, Name: "Finished", Phase: constants.ProjectManagementStagePhaseFinished, SortOrder: 2},
		},
		Transitions: []dto.ProjectManagementWorkflowTransitionDto{{
			FromStageId: -1, ToStageId: -2, AllowedActorTypes: []constants.ProjectManagementMemberRole{constants.ProjectManagementMemberRoleProjectOwner},
		}},
	}
	workflow, stages, transitions, err := projectManagementTemplateModels(input)
	if err != nil {
		t.Fatalf("models: %v", err)
	}
	if workflow.Name != input.Name || len(stages) != 2 || len(transitions) != 1 || transitions[0].AllowedActorTypes[0] != "ProjectOwner" {
		t.Fatalf("unexpected models: %#v %#v %#v", workflow, stages, transitions)
	}
}

func TestValidateProjectManagementTemplateWorkflow(t *testing.T) {
	workflow := dto.ProjectManagementWorkflowDto{
		Name:      "Template",
		Enabled:   true,
		IsDefault: true,
		Stages: []dto.ProjectManagementWorkflowStageDto{
			{Id: 1, Name: "Pending", Phase: constants.ProjectManagementStagePhasePending, SortOrder: 1},
		},
	}
	if err := validateProjectManagementTemplateWorkflow(workflow); err != nil {
		t.Fatalf("valid template rejected: %v", err)
	}
	workflow.Enabled = false
	if err := validateProjectManagementTemplateWorkflow(workflow); err == nil {
		t.Fatal("disabled default template accepted")
	}
}

func TestProjectManagementWorkflowDefaultRules(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	for _, statement := range []string{
		"CREATE TABLE project_management_workflows (id integer primary key, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, kind text not null, organization_id integer not null, project_id integer not null default 0, name text not null default '', description text not null default '', enabled boolean not null default true, is_default boolean not null default false, source_template_id integer not null default 0, source_template_name text not null default '')",
		"CREATE TABLE project_management_workflow_stages (id integer primary key, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, workflow_id integer not null, name text not null default '', phase text not null, sort_order integer not null)",
		"CREATE TABLE project_management_workflow_transitions (id integer primary key, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, workflow_id integer not null, from_stage_id integer not null, to_stage_id integer not null, allowed_actor_types json not null default '[]')",
		"CREATE TABLE operations (id integer primary key autoincrement, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, entity_id integer not null, entity_type text not null, operator_id integer not null, operator_name text not null default '', operator_type text not null, operation text not null, note text not null default '', metadata json not null default '{}', organization_id integer not null)",
	} {
		if err := db.Exec(statement).Error; err != nil {
			t.Fatal(err)
		}
	}
	for _, statement := range []string{
		"INSERT INTO project_management_workflows (id, kind, organization_id, name, enabled, is_default) VALUES (1, 'OrganizationTemplate', 1, 'Default', true, true), (2, 'OrganizationTemplate', 1, 'Other', true, false)",
		"INSERT INTO project_management_workflow_stages (id, workflow_id, name, phase, sort_order) VALUES (11, 1, 'Pending', 'Pending', 1), (21, 2, 'Pending', 'Pending', 1)",
	} {
		if err := db.Exec(statement).Error; err != nil {
			t.Fatal(err)
		}
	}
	logic := NewProjectManagementLogic(NewBaseLogic(repo.NewAllRepos(db), nil, nil, dto.JwtConfig{}))
	actor := &dto.ProjectManagementActor{OrganizationId: 1, OrganizationUserId: 7, Name: "Admin"}
	defaultInput := dto.ProjectManagementWorkflowDto{Name: "Default", Enabled: true, Stages: []dto.ProjectManagementWorkflowStageDto{{Id: 11, Name: "Pending", Phase: constants.ProjectManagementStagePhasePending, SortOrder: 1}}}
	if _, err := logic.UpdateWorkflow(context.Background(), actor, 1, defaultInput); err == nil {
		t.Fatal("default workflow was unset without a replacement")
	}
	if err := logic.DeleteWorkflow(context.Background(), actor, 1); err == nil {
		t.Fatal("default workflow was deleted")
	}
	replacement := dto.ProjectManagementWorkflowDto{Name: "Other", Enabled: true, IsDefault: true, Stages: []dto.ProjectManagementWorkflowStageDto{{Id: 21, Name: "Pending", Phase: constants.ProjectManagementStagePhasePending, SortOrder: 1}}}
	if _, err := logic.UpdateWorkflow(context.Background(), actor, 2, replacement); err != nil {
		t.Fatalf("default workflow replacement rejected: %v", err)
	}
	var defaultCount int64
	if err := db.Model(&model.ProjectManagementWorkflow{}).Where("organization_id = ? AND kind = ? AND is_default = ?", 1, constants.ProjectManagementWorkflowKindOrganizationTemplate, true).Count(&defaultCount).Error; err != nil {
		t.Fatal(err)
	}
	if defaultCount != 1 {
		t.Fatalf("expected one default workflow, got %d", defaultCount)
	}
}
