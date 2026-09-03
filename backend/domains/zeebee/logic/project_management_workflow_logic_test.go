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
		Transitions: []dto.ProjectManagementWorkflowTransitionDto{{Name: "Finish", FromStageId: 1, ToStageId: 2}},
	}
	if err := ValidateProjectManagementWorkflow(valid); err != nil {
		t.Fatalf("valid workflow rejected: %v", err)
	}

	valid.Stages = valid.Stages[1:]
	if err := ValidateProjectManagementWorkflow(valid); err == nil {
		t.Fatal("workflow without pending stage accepted")
	}

	valid.Stages = []dto.ProjectManagementWorkflowStageDto{{Id: 1, Name: "Pending", Phase: constants.ProjectManagementStagePhasePending, SortOrder: 1}}
	if err := ValidateProjectManagementWorkflow(valid); err == nil {
		t.Fatal("workflow without finished stage accepted")
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
	valid.Transitions = []dto.ProjectManagementWorkflowTransitionDto{{Name: "Invalid", FromStageId: 1, ToStageId: 3}}
	if err := ValidateProjectManagementWorkflow(valid); err == nil {
		t.Fatal("workflow transition to an unknown stage accepted")
	}

	valid.Transitions = []dto.ProjectManagementWorkflowTransitionDto{
		{Name: "Finish", FromStageId: 1, ToStageId: 2},
		{Name: "Finish again", FromStageId: 1, ToStageId: 2},
	}
	if err := ValidateProjectManagementWorkflow(valid); err == nil {
		t.Fatal("duplicate workflow transition accepted")
	}
}

func TestValidateProjectManagementWorkflowTransitionMessages(t *testing.T) {
	tests := []struct {
		name  string
		input dto.ProjectManagementWorkflowTransitionDto
		want  string
	}{
		{"missing stage ID", dto.ProjectManagementWorkflowTransitionDto{Name: "Move"}, "workflow transition source and target stage IDs are required"},
		{"same stage", dto.ProjectManagementWorkflowTransitionDto{Name: "Move", FromStageId: 1, ToStageId: 1}, "workflow transition source and target stages must differ (both are 1)"},
		{"invalid actor", dto.ProjectManagementWorkflowTransitionDto{Name: "Move", FromStageId: 1, ToStageId: 2, AllowedActorTypes: []constants.ProjectManagementMemberRole{"invalid"}}, "workflow transition allowed actor type \"invalid\" is invalid"},
		{"duplicate actor", dto.ProjectManagementWorkflowTransitionDto{Name: "Move", FromStageId: 1, ToStageId: 2, AllowedActorTypes: []constants.ProjectManagementMemberRole{constants.ProjectManagementMemberRoleProjectAdmin, constants.ProjectManagementMemberRoleProjectAdmin}}, "workflow transition allowed actor type \"ProjectAdmin\" is duplicated"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateProjectManagementWorkflowTransition(tt.input)
			if err == nil || err.Error() != tt.want {
				t.Fatalf("got %v, want %q", err, tt.want)
			}
		})
	}
}

func TestProjectManagementWorkflowUpdateDefinition(t *testing.T) {
	definition, err := projectManagementWorkflowUpdateDefinition(dto.ProjectManagementWorkflowUpdateDto{
		Name: "Template",
		Stages: []dto.ProjectManagementWorkflowUpdateStageDto{
			{Id: 1, Name: "Pending", Phase: constants.ProjectManagementStagePhasePending, SortOrder: 1, Transitions: []dto.ProjectManagementWorkflowUpdateTransitionDto{{Id: 10, Name: "Finish", FromStageId: 1, ToStageId: 2}}},
			{Name: "Finished", Phase: constants.ProjectManagementStagePhaseFinished, SortOrder: 2},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(definition.Stages) != 2 || definition.Stages[1].Id == 0 || len(definition.Transitions) != 1 || definition.Transitions[0].Id != 10 || definition.Transitions[0].FromStageId != 1 || definition.Transitions[0].ToStageId != 2 {
		t.Fatalf("unexpected workflow definition: %#v", definition)
	}
	_, err = projectManagementWorkflowUpdateDefinition(dto.ProjectManagementWorkflowUpdateDto{Stages: []dto.ProjectManagementWorkflowUpdateStageDto{{
		Id: 1, Transitions: []dto.ProjectManagementWorkflowUpdateTransitionDto{{FromStageId: 2, ToStageId: 2}},
	}}})
	if err == nil {
		t.Fatal("transition with a mismatched source stage accepted")
	}
	_, err = projectManagementWorkflowUpdateDefinition(dto.ProjectManagementWorkflowUpdateDto{Stages: []dto.ProjectManagementWorkflowUpdateStageDto{{
		Transitions: []dto.ProjectManagementWorkflowUpdateTransitionDto{{FromStageId: 1, ToStageId: 2}},
	}}})
	if err == nil {
		t.Fatal("new stage with transitions accepted")
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
			Name: "Finish", FromStageId: -1, ToStageId: -2, AllowedActorTypes: []constants.ProjectManagementMemberRole{constants.ProjectManagementMemberRoleProjectOwner},
		}},
	}
	workflow, stages, transitions, err := projectManagementTemplateModels(input)
	if err != nil {
		t.Fatalf("models: %v", err)
	}
	if workflow.Name != input.Name || len(stages) != 2 || len(transitions) != 1 || transitions[0].Name != "Finish" || transitions[0].AllowedActorTypes[0] != "ProjectOwner" {
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

func TestProjectManagementWorkflowDeleteRejectsProjectWorkflow(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err = db.Exec("CREATE TABLE project_management_workflows (id integer primary key, is_del integer not null default 0, organization_id integer not null, kind text not null)").Error; err != nil {
		t.Fatal(err)
	}
	for _, statement := range []string{
		"INSERT INTO project_management_workflows (id, organization_id, kind) VALUES (1, 1, 'OrganizationTemplate')",
		"INSERT INTO project_management_workflows (id, organization_id, kind) VALUES (2, 1, 'Project')",
	} {
		if err = db.Exec(statement).Error; err != nil {
			t.Fatal(err)
		}
	}

	logic := NewProjectManagementWorkflowLogic(NewBaseLogic(repo.NewAllRepos(db), nil, nil, dto.JwtConfig{}))
	if err = logic.Delete(context.Background(), utils.JwtUser{OrganizationId: 1}, 2); err == nil {
		t.Fatal("project workflow deletion accepted")
	}
	var workflow model.ProjectManagementWorkflow
	if err = db.First(&workflow, 2).Error; err != nil {
		t.Fatalf("project workflow deleted: %v", err)
	}
}

func TestProjectManagementWorkflowEnsureDefaultSetsStageColors(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	for _, statement := range []string{
		"CREATE TABLE project_management_workflows (id integer primary key, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, kind text not null, organization_id integer not null, project_id integer not null default 0, name text not null default '', description text not null default '', is_default boolean not null default false, source_template_id integer not null default 0, source_template_name text not null default '')",
		"CREATE TABLE project_management_workflow_stages (id integer primary key, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, workflow_id integer not null, name text not null default '', description text not null default '', phase text not null, color text not null default '', sort_order integer not null)",
		"CREATE TABLE project_management_workflow_transitions (id integer primary key, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, workflow_id integer not null, name text not null default '', from_stage_id integer not null, to_stage_id integer not null, allowed_actor_types json not null default '[]')",
	} {
		if err = db.Exec(statement).Error; err != nil {
			t.Fatal(err)
		}
	}

	logic := NewProjectManagementWorkflowLogic(NewBaseLogic(repo.NewAllRepos(db), nil, nil, dto.JwtConfig{}))
	if err = logic.ensureDefault(context.Background(), 1); err != nil {
		t.Fatalf("ensure default workflow: %v", err)
	}
	var stages []model.ProjectManagementWorkflowStage
	if err = db.Order("sort_order").Find(&stages).Error; err != nil {
		t.Fatalf("list default stages: %v", err)
	}
	if len(stages) != 3 || stages[0].Color != "#3b82f6" || stages[1].Color != "#f59e0b" || stages[2].Color != "#22c55e" {
		t.Fatalf("unexpected default stage colors: %#v", stages)
	}
	var transitions []model.ProjectManagementWorkflowTransition
	if err = db.Find(&transitions).Error; err != nil {
		t.Fatalf("list default transitions: %v", err)
	}
	if len(transitions) != len(stages)*(len(stages)-1) {
		t.Fatalf("default transitions = %d, want %d", len(transitions), len(stages)*(len(stages)-1))
	}
	edges := make(map[[2]int64]struct{}, len(transitions))
	for _, transition := range transitions {
		edges[[2]int64{transition.FromStageId, transition.ToStageId}] = struct{}{}
	}
	for _, from := range stages {
		for _, to := range stages {
			if from.Id != to.Id {
				if _, ok := edges[[2]int64{from.Id, to.Id}]; !ok {
					t.Fatalf("missing default transition %d -> %d", from.Id, to.Id)
				}
			}
		}
	}
}

func TestProjectManagementWorkflowCreateCopiesDefaultAndUpdateReplacesDefinition(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	for _, statement := range []string{
		"CREATE TABLE project_management_workflows (id integer primary key, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, kind text not null, organization_id integer not null, project_id integer not null default 0, name text not null default '', description text not null default '', is_default boolean not null default false, source_template_id integer not null default 0, source_template_name text not null default '')",
		"CREATE TABLE project_management_workflow_stages (id integer primary key, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, workflow_id integer not null, name text not null default '', description text not null default '', phase text not null, color text not null default '', sort_order integer not null)",
		"CREATE TABLE project_management_workflow_transitions (id integer primary key, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, workflow_id integer not null, name text not null default '', from_stage_id integer not null, to_stage_id integer not null, allowed_actor_types json not null default '[]')",
		"CREATE TABLE operations (id integer primary key autoincrement, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, entity_id integer not null, entity_type text not null, operator_id integer not null, operator_name text not null default '', operator_type text not null, operation text not null, note text not null default '', metadata json not null default '{}', organization_id integer not null)",
	} {
		if err := db.Exec(statement).Error; err != nil {
			t.Fatal(err)
		}
	}
	for _, statement := range []string{
		"INSERT INTO project_management_workflows (id, kind, organization_id, name, is_default) VALUES (1, 'OrganizationTemplate', 1, 'Default', true)",
		"INSERT INTO project_management_workflow_stages (id, workflow_id, name, phase, color, sort_order) VALUES (11, 1, 'Pending', 'Pending', '#3b82f6', 1), (12, 1, 'Finished', 'Finished', '#22c55e', 2)",
		"INSERT INTO project_management_workflow_transitions (id, workflow_id, name, from_stage_id, to_stage_id, allowed_actor_types) VALUES (13, 1, 'Finish', 11, 12, '[\"ProjectOwner\"]')",
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
	if created.IsDefault || len(created.Stages) != 2 || created.Stages[0].Name != "Pending" || created.Stages[0].Color != "#3b82f6" || created.Stages[1].Phase != constants.ProjectManagementStagePhaseFinished || len(created.Stages[0].Transitions) != 1 || len(created.Stages[0].Transitions[0].AllowedActorTypes) != 1 || created.Stages[0].Transitions[0].AllowedActorTypes[0] != constants.ProjectManagementMemberRoleProjectOwner {
		t.Fatalf("unexpected copied workflow: %#v", created)
	}
	if created.Stages[0].Id == 11 || created.Stages[1].Id == 12 || created.Stages[0].Transitions[0].FromStageId == 11 || created.Stages[0].Transitions[0].ToStageId == 12 {
		t.Fatalf("default workflow definition was not cloned: %#v", created)
	}

	updated, err := logic.Update(context.Background(), actor, created.Id, dto.ProjectManagementWorkflowUpdateDto{
		Name:        "Renamed",
		Description: "Updated definition",
		Stages: []dto.ProjectManagementWorkflowUpdateStageDto{
			{Id: created.Stages[0].Id, Name: "Queued", Description: "Waiting for review", Phase: constants.ProjectManagementStagePhasePending, SortOrder: 1, Transitions: []dto.ProjectManagementWorkflowUpdateTransitionDto{{
				Id: created.Stages[0].Transitions[0].Id, Name: "Complete", FromStageId: created.Stages[0].Id, ToStageId: created.Stages[1].Id,
			}}},
			{Id: created.Stages[1].Id, Name: "Completed", Phase: constants.ProjectManagementStagePhaseFinished, SortOrder: 2},
			{Name: "Review", Phase: constants.ProjectManagementStagePhaseInProgress},
		},
	})
	if err != nil {
		t.Fatalf("update workflow: %v", err)
	}
	if updated.Name != "Renamed" || updated.Description != "Updated definition" || updated.IsDefault || len(updated.Stages) != 3 || updated.Stages[0].Name != "Queued" || updated.Stages[0].Description != "Waiting for review" || len(updated.Stages[0].Transitions) != 1 || updated.Stages[0].Transitions[0].Name != "Complete" || updated.Stages[2].SortOrder != 3 {
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

func TestWorkflowStageAndTransitionRequireMatchingWorkflow(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	for _, statement := range []string{
		"CREATE TABLE project_management_workflows (id integer primary key, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, kind text not null, organization_id integer not null, project_id integer not null default 0, name text not null default '', description text not null default '', is_default boolean not null default false, source_template_id integer not null default 0, source_template_name text not null default '')",
		"CREATE TABLE project_management_workflow_stages (id integer primary key, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, workflow_id integer not null, name text not null default '', description text not null default '', phase text not null, color text not null default '', sort_order integer not null)",
		"CREATE TABLE project_management_workflow_transitions (id integer primary key, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, workflow_id integer not null, name text not null default '', from_stage_id integer not null, to_stage_id integer not null, allowed_actor_types json not null default '[]')",
		"CREATE TABLE project_management_tasks (id integer primary key, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, stage_id integer not null)",
		"CREATE TABLE operations (id integer primary key autoincrement, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, entity_id integer not null, entity_type text not null, operator_id integer not null, operator_name text not null default '', operator_type text not null, operation text not null, note text not null default '', metadata json not null default '{}', organization_id integer not null)",
	} {
		if err := db.Exec(statement).Error; err != nil {
			t.Fatal(err)
		}
	}
	for _, statement := range []string{
		"INSERT INTO project_management_workflows (id, kind, organization_id, name) VALUES (1, 'OrganizationTemplate', 1, 'One'), (2, 'OrganizationTemplate', 1, 'Two')",
		"INSERT INTO project_management_workflow_stages (id, workflow_id, name, phase, sort_order) VALUES (11, 1, 'Pending', 'Pending', 1), (12, 1, 'Done', 'Finished', 2), (21, 2, 'Pending', 'Pending', 1)",
		"INSERT INTO project_management_workflow_transitions (id, workflow_id, name, from_stage_id, to_stage_id, allowed_actor_types) VALUES (31, 1, 'Finish', 11, 12, '[]')",
	} {
		if err := db.Exec(statement).Error; err != nil {
			t.Fatal(err)
		}
	}
	workflowLogic := NewProjectManagementWorkflowLogic(NewBaseLogic(repo.NewAllRepos(db), nil, nil, dto.JwtConfig{}))
	ctx := context.Background()
	stage, err := workflowLogic.StageDetail(ctx, 1, 1, 11)
	if err != nil || stage.Id != 11 {
		t.Fatalf("stage detail: %#v, %v", stage, err)
	}
	transitions, err := workflowLogic.TransitionAll(ctx, 1, 1, 11)
	if err != nil || len(transitions) != 1 || transitions[0].Id != 31 {
		t.Fatalf("transition all: %#v, %v", transitions, err)
	}
	if _, err := workflowLogic.StageDetail(ctx, 1, 2, 11); err == nil {
		t.Fatal("stage accepted a mismatched workflow")
	}
	if _, err := workflowLogic.TransitionAll(ctx, 1, 2, 11); err == nil {
		t.Fatal("transition list accepted a mismatched workflow")
	}
	if _, err := workflowLogic.TransitionDetail(ctx, 1, 2, 31); err == nil {
		t.Fatal("transition accepted a mismatched workflow")
	}
	actor := utils.JwtUser{OrganizationId: 1, OrganizationUserId: 1, Username: "tester", UserType: constants.UserTypeBorrower}
	if _, err := workflowLogic.StageCreate(ctx, actor, dto.ProjectManagementWorkflowStageDto{Name: "Working", Phase: constants.ProjectManagementStagePhaseInProgress, SortOrder: 3}, 1); err != nil {
		t.Fatalf("create in-progress stage: %v", err)
	}
	if _, err := workflowLogic.StageCreate(ctx, actor, dto.ProjectManagementWorkflowStageDto{Name: "Another pending", Phase: constants.ProjectManagementStagePhasePending, SortOrder: 4}, 1); err == nil {
		t.Fatal("create accepted a second pending stage")
	}
	if _, err := workflowLogic.StageUpdate(ctx, actor, 1, 12, dto.ProjectManagementWorkflowStageDto{Name: "Done", Phase: constants.ProjectManagementStagePhaseInProgress, SortOrder: 2}); err == nil {
		t.Fatal("update accepted removal of the only finished stage")
	}
	if err := workflowLogic.StageDelete(ctx, actor, 1, 12); err == nil {
		t.Fatal("delete accepted removal of the only finished stage")
	}
	if _, err := workflowLogic.StageCreate(ctx, actor, dto.ProjectManagementWorkflowStageDto{Name: "Closed", Phase: constants.ProjectManagementStagePhaseFinished, SortOrder: 4}, 1); err != nil {
		t.Fatalf("create finished stage: %v", err)
	}
	if err := db.Exec("INSERT INTO project_management_tasks (id, stage_id) VALUES (41, 12)").Error; err != nil {
		t.Fatal(err)
	}
	if err := workflowLogic.StageDelete(ctx, actor, 1, 12); err == nil {
		t.Fatal("delete accepted a stage with active tasks")
	}
}
