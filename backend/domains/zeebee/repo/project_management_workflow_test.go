package repo

import (
	"testing"

	"gorm.io/gorm"

	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"
)

func TestProjectManagementWorkflowPreloadsStagesAndTransitions(t *testing.T) {
	db := newRepoTestDB(t)
	execTestSQL(t, db, []string{
		"CREATE TABLE project_management_workflows (id integer primary key autoincrement, created_at datetime, updated_at datetime, is_del integer not null default 0, kind text not null default '', organization_id integer not null default 0, project_id integer not null default 0, name text not null default '', description text not null default '', is_default boolean not null default false, source_template_id integer not null default 0, source_template_name text not null default '', entity_type text not null default 'Project', entity_id integer not null default 0, source_workflow_id integer not null default 0, source_workflow_name text not null default '')",
		"CREATE TABLE project_management_stages (id integer primary key autoincrement, created_at datetime, updated_at datetime, is_del integer not null default 0, entity_type text not null, entity_id integer not null default 0, name text not null default '', description text not null default '', phase text not null, color text not null default '', sort_order integer not null default 0)",
		"CREATE TABLE project_management_workflow_stages (id integer primary key autoincrement, is_del integer not null default 0, workflow_id integer not null, stage_id integer not null)",
		"CREATE TABLE project_management_workflow_transitions (id integer primary key autoincrement, created_at datetime, updated_at datetime, is_del integer not null default 0, workflow_id integer not null default 0, name text not null default '', from_stage_id integer not null default 0, to_stage_id integer not null default 0, allowed_actor_types json not null default '[]')",
	})
	workflow := &model.ProjectManagementWorkflow{EntityType: constants.ProjectManagementOwnerTypeProject, EntityId: 1, Name: "Workflow"}
	stages := []*model.ProjectManagementStage{
		{EntityType: constants.ProjectManagementOwnerTypeProject, EntityId: 1, Name: "Finished", Phase: constants.ProjectManagementStagePhaseFinished, SortOrder: 2},
		{EntityType: constants.ProjectManagementOwnerTypeProject, EntityId: 1, Name: "Pending", Phase: constants.ProjectManagementStagePhasePending, SortOrder: 1},
	}
	if err := db.Create(workflow).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Create(stages).Error; err != nil {
		t.Fatal(err)
	}
	links := []*model.ProjectManagementWorkflowStage{{WorkflowId: workflow.Id, StageId: stages[0].Id}, {WorkflowId: workflow.Id, StageId: stages[1].Id}}
	if err := db.Create(links).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&model.ProjectManagementWorkflowStageTransition{WorkflowId: workflow.Id, Name: "Finish", FromStageId: stages[1].Id, ToStageId: stages[0].Id}).Error; err != nil {
		t.Fatal(err)
	}

	got := new(model.ProjectManagementWorkflow)
	if err := db.Preload("Stages", func(db *gorm.DB) *gorm.DB { return db.Order("sort_order ASC") }).Preload("Transitions", preloadWorkflowTransitions).First(got, workflow.Id).Error; err != nil {
		t.Fatal(err)
	}
	if len(got.Stages) != 2 || got.Stages[0].Id != stages[1].Id || len(got.Transitions) != 1 {
		t.Fatalf("unexpected workflow relationships: %#v", got)
	}
	if got.Transitions[0].FromStageId != stages[1].Id || got.Transitions[0].ToStageId != stages[0].Id {
		t.Fatalf("transition stage ids = %#v", got.Transitions[0])
	}
}

func TestProjectManagementWorkflowRepoPendingStages(t *testing.T) {
	db := newRepoTestDB(t)
	execTestSQL(t, db, []string{
		"CREATE TABLE project_management_stages (id integer primary key, is_del integer not null default 0, phase text not null)",
		"CREATE TABLE project_management_workflow_stages (id integer primary key, is_del integer not null default 0, workflow_id integer not null, stage_id integer not null)",
		"INSERT INTO project_management_stages (id, phase) VALUES (11, 'Pending'), (12, 'Pending')",
		"INSERT INTO project_management_workflow_stages (id, workflow_id, stage_id) VALUES (1, 100, 11), (2, 200, 12)",
	})
	stages, err := NewProjectManagementWorkflowRepo(db).ListPendingOwnerStagesWithTx(t.Context(), db, []int64{100, 200})
	if err != nil || stages[100] == nil || stages[100].Id != 11 || stages[200] == nil || stages[200].Id != 12 {
		t.Fatalf("pending stages = %#v, %v", stages, err)
	}
}
