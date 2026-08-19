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
		"CREATE TABLE project_management_workflow_stages (id integer primary key autoincrement, created_at datetime, updated_at datetime, is_del integer not null default 0, workflow_id integer not null, name text not null, phase text not null, sort_order integer not null)",
		"CREATE TABLE project_management_workflow_transitions (id integer primary key autoincrement, created_at datetime, updated_at datetime, is_del integer not null default 0, workflow_id integer not null, from_stage_id integer not null, to_stage_id integer not null, allowed_actor_types json not null default '[]')",
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
