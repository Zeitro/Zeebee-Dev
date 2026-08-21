package logic

import (
	"context"
	"fmt"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
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

func TestValidateProjectManagementWorkflowDetails(t *testing.T) {
	if err := validateProjectManagementWorkflowDetails(dto.ProjectManagementWorkflowDto{Name: "Template"}); err != nil {
		t.Fatalf("valid workflow details rejected: %v", err)
	}
	if err := validateProjectManagementWorkflowDetails(dto.ProjectManagementWorkflowDto{Name: " "}); err == nil {
		t.Fatal("blank workflow name accepted")
	}
}

func TestProjectManagementWorkflowCreateCopiesDefaultAndUpdatePreservesDefinition(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	for _, statement := range []string{
		"CREATE TABLE project_management_workflows (id integer primary key, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, kind text not null, organization_id integer not null, project_id integer not null default 0, name text not null default '', description text not null default '', is_default boolean not null default false, source_template_id integer not null default 0, source_template_name text not null default '')",
		"CREATE TABLE project_management_workflow_stages (id integer primary key, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, workflow_id integer not null, name text not null default '', phase text not null, sort_order integer not null)",
		"CREATE TABLE project_management_workflow_transitions (id integer primary key, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, workflow_id integer not null, from_stage_id integer not null, to_stage_id integer not null, allowed_actor_types json not null default '[]')",
		"CREATE TABLE operations (id integer primary key autoincrement, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, entity_id integer not null, entity_type text not null, operator_id integer not null, operator_name text not null default '', operator_type text not null, operation text not null, note text not null default '', metadata json not null default '{}', organization_id integer not null)",
	} {
		if err := db.Exec(statement).Error; err != nil {
			t.Fatal(err)
		}
	}
	for _, statement := range []string{
		"INSERT INTO project_management_workflows (id, kind, organization_id, name, is_default) VALUES (1, 'OrganizationTemplate', 1, 'Default', true)",
		"INSERT INTO project_management_workflow_stages (id, workflow_id, name, phase, sort_order) VALUES (11, 1, 'Pending', 'Pending', 1), (12, 1, 'Finished', 'Finished', 2)",
		"INSERT INTO project_management_workflow_transitions (id, workflow_id, from_stage_id, to_stage_id, allowed_actor_types) VALUES (13, 1, 11, 12, '[\"ProjectOwner\"]')",
	} {
		if err := db.Exec(statement).Error; err != nil {
			t.Fatal(err)
		}
	}
	logic := NewProjectManagementWorkflowLogic(NewBaseLogic(repo.NewAllRepos(db), nil, nil, dto.JwtConfig{}))
	actor := utils.JwtUser{OrganizationId: 1, OrganizationUserId: 7, Username: "Admin"}
	created, err := logic.Create(context.Background(), actor, dto.ProjectManagementWorkflowDto{Name: "Copy", Description: "Copied from default"})
	if err != nil {
		t.Fatalf("create workflow: %v", err)
	}
	if created.IsDefault || len(created.Stages) != 2 || created.Stages[0].Name != "Pending" || created.Stages[1].Phase != constants.ProjectManagementStagePhaseFinished || len(created.Transitions) != 1 || len(created.Transitions[0].AllowedActorTypes) != 1 || created.Transitions[0].AllowedActorTypes[0] != constants.ProjectManagementMemberRoleProjectOwner {
		t.Fatalf("unexpected copied workflow: %#v", created)
	}
	if created.Stages[0].Id == 11 || created.Stages[1].Id == 12 || created.Transitions[0].FromStageId == 11 || created.Transitions[0].ToStageId == 12 {
		t.Fatalf("default workflow definition was not cloned: %#v", created)
	}

	updated, err := logic.Update(context.Background(), actor, created.Id, dto.ProjectManagementWorkflowDto{
		Name: "Renamed", Description: "Only details change", IsDefault: true,
		Stages: []dto.ProjectManagementWorkflowStageDto{{Id: 99, Name: "Ignored", Phase: constants.ProjectManagementStagePhasePending}},
	})
	if err != nil {
		t.Fatalf("update workflow: %v", err)
	}
	if updated.Name != "Renamed" || updated.Description != "Only details change" || updated.IsDefault || len(updated.Stages) != 2 || len(updated.Transitions) != 1 {
		t.Fatalf("workflow update changed definition: %#v", updated)
	}

	if err := logic.Delete(context.Background(), actor, 1); err == nil {
		t.Fatal("default workflow was deleted")
	}
	if err := logic.SetDefault(context.Background(), actor, created.Id); err != nil {
		t.Fatalf("set default workflow: %v", err)
	}
	defaultWorkflow, err := logic.Detail(context.Background(), actor.OrganizationId, created.Id)
	if err != nil || !defaultWorkflow.IsDefault || defaultWorkflow.Id != created.Id {
		t.Fatalf("unexpected default workflow: %#v, %v", defaultWorkflow, err)
	}
	if err := logic.Delete(context.Background(), actor, created.Id); err == nil {
		t.Fatal("new default workflow was deleted")
	}
	if err := logic.Delete(context.Background(), actor, 1); err != nil {
		t.Fatalf("delete previous default workflow: %v", err)
	}
}
