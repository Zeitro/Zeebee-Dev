package repo

import (
	"context"
	"fmt"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"
)

func TestReplaceProjectDefinitionPreservesReferencedStages(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	for _, statement := range []string{
		"CREATE TABLE project_management_workflow_stages (id integer primary key autoincrement, created_at datetime, updated_at datetime, is_del integer not null default 0, workflow_id integer not null, name text not null, phase text not null, color text not null default '', sort_order integer not null)",
		"CREATE TABLE project_management_workflow_transitions (id integer primary key autoincrement, created_at datetime, updated_at datetime, is_del integer not null default 0, workflow_id integer not null, name text not null default '', from_stage_id integer not null, to_stage_id integer not null, allowed_actor_types json not null default '[]')",
		"CREATE TABLE project_management_tasks (id integer primary key autoincrement, created_at datetime, updated_at datetime, is_del integer not null default 0, project_id integer not null default 0, organization_id integer not null default 0, parent_task_id integer not null default 0, stage_id integer not null, creator_id integer not null default 0, creator_name text not null default '', title text not null default '', description text not null default '', priority text not null default '', start_at datetime, end_at datetime, task_type text not null default '{}')",
	} {
		if err := db.Exec(statement).Error; err != nil {
			t.Fatal(err)
		}
	}
	repo := NewProjectManagementWorkflowRepo(db)
	kept := &model.ProjectManagementWorkflowStage{WorkflowId: 1, Name: "Pending", Phase: constants.ProjectManagementStagePhasePending, SortOrder: 10}
	template := &model.ProjectManagementWorkflowStage{WorkflowId: 2, Name: "Template", Phase: constants.ProjectManagementStagePhasePending, SortOrder: 10}
	if err := db.Create(kept).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Create(template).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Transaction(func(tx *gorm.DB) error {
		return repo.ReplaceProjectDefinitionWithTx(context.Background(), tx, 1, []*model.ProjectManagementWorkflowStage{
			{Base: model.Base{Id: kept.Id}, Name: "Not started", Phase: constants.ProjectManagementStagePhasePending, SortOrder: 10},
			{Base: model.Base{Id: -1}, Name: "In progress", Phase: constants.ProjectManagementStagePhaseInProgress, SortOrder: 20},
		}, nil)
	}); err != nil {
		t.Fatalf("replace project workflow: %v", err)
	}
	var got model.ProjectManagementWorkflowStage
	if err := db.First(&got, kept.Id).Error; err != nil || got.Name != "Not started" {
		t.Fatalf("referenced stage was not preserved: %#v, %v", got, err)
	}
	got = model.ProjectManagementWorkflowStage{}
	if err := db.First(&got, template.Id).Error; err != nil || got.Name != "Template" {
		t.Fatalf("template stage was modified: %#v, %v", got, err)
	}
	if err := db.Create(&model.ProjectManagementTask{StageId: kept.Id}).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Transaction(func(tx *gorm.DB) error {
		return repo.ReplaceProjectDefinitionWithTx(context.Background(), tx, 1, []*model.ProjectManagementWorkflowStage{
			{Base: model.Base{Id: -1}, Name: "New pending", Phase: constants.ProjectManagementStagePhasePending, SortOrder: 10},
		}, nil)
	}); err == nil {
		t.Fatal("replacement deleted a task-referenced stage")
	}
	got = model.ProjectManagementWorkflowStage{}
	if err := db.First(&got, kept.Id).Error; err != nil {
		t.Fatalf("referenced stage deleted after rejected replacement: %v", err)
	}
}

func TestFindTemplateStagePreloadsTransitions(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	for _, statement := range []string{
		"CREATE TABLE project_management_workflows (id integer primary key, is_del integer not null default 0, organization_id integer not null, kind text not null)",
		"CREATE TABLE project_management_workflow_stages (id integer primary key, is_del integer not null default 0, workflow_id integer not null, name text not null, phase text not null, color text not null default '', sort_order integer not null)",
		"CREATE TABLE project_management_workflow_transitions (id integer primary key, is_del integer not null default 0, workflow_id integer not null, name text not null default '', from_stage_id integer not null, to_stage_id integer not null, allowed_actor_types json not null default '[]')",
		"INSERT INTO project_management_workflows (id, organization_id, kind) VALUES (1, 4, 'OrganizationTemplate')",
		"INSERT INTO project_management_workflow_stages (id, workflow_id, name, phase, sort_order) VALUES (10, 1, 'Pending', 'Pending', 1), (11, 1, 'In progress', 'InProgress', 2), (12, 1, 'Finished', 'Finished', 3)",
		"INSERT INTO project_management_workflow_transitions (id, workflow_id, name, from_stage_id, to_stage_id) VALUES (2, 1, 'Pause', 10, 11), (1, 1, 'Finish', 10, 12), (3, 1, 'Finish work', 11, 12)",
	} {
		if err := db.Exec(statement).Error; err != nil {
			t.Fatal(err)
		}
	}

	stage, err := NewProjectManagementWorkflowRepo(db).FindTemplateStage(context.Background(), 10, 1, 4)
	if err != nil {
		t.Fatal(err)
	}
	if len(stage.Transitions) != 2 || stage.Transitions[0].Id != 1 || stage.Transitions[1].Id != 2 {
		t.Fatalf("unexpected transitions: %#v", stage.Transitions)
	}
	transitions, err := NewProjectManagementWorkflowRepo(db).ListTransitionsByFromStage(context.Background(), 10)
	if err != nil || len(transitions) != 2 || transitions[0].FromStage.Name != "Pending" || transitions[0].ToStage.Name != "Finished" {
		t.Fatalf("unexpected transition list: %#v, %v", transitions, err)
	}
	transition, err := NewProjectManagementWorkflowRepo(db).FindTemplateTransition(context.Background(), 2, 1, 4)
	if err != nil || transition.FromStage.Name != "Pending" || transition.ToStage.Name != "In progress" {
		t.Fatalf("unexpected transition detail: %#v, %v", transition, err)
	}
}

func TestUpdateTemplateTransitionSerializesAllowedActorTypes(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Exec("CREATE TABLE project_management_workflow_transitions (id integer primary key, created_at datetime, updated_at datetime, is_del integer not null default 0, workflow_id integer not null, name text not null default '', from_stage_id integer not null, to_stage_id integer not null, allowed_actor_types json not null default '[]')").Error; err != nil {
		t.Fatal(err)
	}

	repo := NewProjectManagementWorkflowRepo(db)
	transition := &model.ProjectManagementWorkflowTransition{Base: model.Base{Id: 1}, WorkflowId: 1, Name: "Restart", FromStageId: 1, ToStageId: 2, AllowedActorTypes: []constants.ProjectManagementMemberRole{constants.ProjectManagementMemberRoleProjectOwner}}
	if err := db.Create(transition).Error; err != nil {
		t.Fatal(err)
	}
	transition.FromStageId = 3
	transition.AllowedActorTypes = []constants.ProjectManagementMemberRole{constants.ProjectManagementMemberRoleProjectOwner, constants.ProjectManagementMemberRoleTaskExecutor}
	if err := db.Transaction(func(tx *gorm.DB) error {
		return repo.UpdateTemplateTransitionWithTx(context.Background(), tx, transition)
	}); err != nil {
		t.Fatalf("update transition allowed actor types: %v", err)
	}

	var got model.ProjectManagementWorkflowTransition
	if err := db.First(&got, transition.Id).Error; err != nil {
		t.Fatal(err)
	}
	if got.FromStageId != 3 {
		t.Fatalf("unexpected from stage id: %d", got.FromStageId)
	}
	if len(got.AllowedActorTypes) != 2 || got.AllowedActorTypes[0] != constants.ProjectManagementMemberRoleProjectOwner || got.AllowedActorTypes[1] != constants.ProjectManagementMemberRoleTaskExecutor {
		t.Fatalf("unexpected allowed actor types: %#v", got.AllowedActorTypes)
	}
}
