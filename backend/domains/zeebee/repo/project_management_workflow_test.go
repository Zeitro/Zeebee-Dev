package repo

import (
	"context"
	"fmt"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/proto"
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

func TestReplaceProjectDefinitionUpdatesTransitions(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	for _, statement := range []string{
		"CREATE TABLE project_management_workflow_stages (id integer primary key autoincrement, created_at datetime, updated_at datetime, is_del integer not null default 0, workflow_id integer not null, name text not null, phase text not null, color text not null default '', sort_order integer not null)",
		"CREATE TABLE project_management_workflow_transitions (id integer primary key autoincrement, created_at datetime, updated_at datetime, is_del integer not null default 0, workflow_id integer not null, name text not null default '', from_stage_id integer not null, to_stage_id integer not null, allowed_actor_types json not null default '[]')",
		"CREATE TABLE project_management_tasks (id integer primary key autoincrement, created_at datetime, updated_at datetime, is_del integer not null default 0, stage_id integer not null)",
		"INSERT INTO project_management_workflow_stages (id, workflow_id, name, phase, sort_order) VALUES (1, 1, 'Pending', 'Pending', 1), (2, 1, 'Finished', 'Finished', 2), (3, 1, 'Removed', 'InProgress', 3)",
		"INSERT INTO project_management_workflow_transitions (id, workflow_id, name, from_stage_id, to_stage_id) VALUES (10, 1, 'Old', 1, 2), (11, 1, 'Remove', 2, 3)",
	} {
		if err := db.Exec(statement).Error; err != nil {
			t.Fatal(err)
		}
	}

	repo := NewProjectManagementWorkflowRepo(db)
	if err := db.Transaction(func(tx *gorm.DB) error {
		return repo.ReplaceProjectDefinitionWithTx(context.Background(), tx, 1,
			[]*model.ProjectManagementWorkflowStage{
				{Base: model.Base{Id: 1}, Name: "Not started", Phase: constants.ProjectManagementStagePhasePending, SortOrder: 1},
				{Base: model.Base{Id: 2}, Name: "Finished", Phase: constants.ProjectManagementStagePhaseFinished, SortOrder: 2},
			},
			[]*model.ProjectManagementWorkflowTransition{
				{Base: model.Base{Id: 10}, Name: "Complete", FromStageId: 1, ToStageId: 2},
				{Base: model.Base{Id: -1}, Name: "Reopen", FromStageId: 2, ToStageId: 1},
			})
	}); err != nil {
		t.Fatalf("replace project definition: %v", err)
	}
	var updated model.ProjectManagementWorkflowTransition
	if err := db.First(&updated, 10).Error; err != nil || updated.Name != "Complete" {
		t.Fatalf("existing transition was not updated: %#v, %v", updated, err)
	}
	if err := db.First(&model.ProjectManagementWorkflowTransition{}, 11).Error; err == nil {
		t.Fatal("removed transition still exists")
	}
	var transitions []model.ProjectManagementWorkflowTransition
	if err := db.Where("workflow_id = ?", 1).Order("id ASC").Find(&transitions).Error; err != nil || len(transitions) != 2 || transitions[1].Name != "Reopen" {
		t.Fatalf("new transition was not created: %#v, %v", transitions, err)
	}
	if err := db.First(&model.ProjectManagementWorkflowStage{}, 3).Error; err == nil {
		t.Fatal("removed stage still exists")
	}
	if err := db.Transaction(func(tx *gorm.DB) error {
		return repo.ReplaceProjectDefinitionWithTx(context.Background(), tx, 1,
			[]*model.ProjectManagementWorkflowStage{
				{Base: model.Base{Id: 1}, Name: "Not started", Phase: constants.ProjectManagementStagePhasePending, SortOrder: 1},
				{Base: model.Base{Id: 2}, Name: "Finished", Phase: constants.ProjectManagementStagePhaseFinished, SortOrder: 2},
			},
			[]*model.ProjectManagementWorkflowTransition{{Base: model.Base{Id: 999}, FromStageId: 1, ToStageId: 2}})
	}); err == nil {
		t.Fatal("foreign workflow transition accepted")
	}
}

func TestWorkflowDefinitionsPreloadTransitionStages(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	for _, statement := range []string{
		"CREATE TABLE project_management_workflows (id integer primary key, is_del integer not null default 0, organization_id integer not null, project_id integer not null default 0, kind text not null)",
		"CREATE TABLE project_management_workflow_stages (id integer primary key, created_at datetime, updated_at datetime, is_del integer not null default 0, workflow_id integer not null, name text not null, phase text not null, color text not null default '', sort_order integer not null)",
		"CREATE TABLE project_management_workflow_transitions (id integer primary key, is_del integer not null default 0, workflow_id integer not null, name text not null default '', from_stage_id integer not null, to_stage_id integer not null, allowed_actor_types json not null default '[]')",
		"INSERT INTO project_management_workflows (id, organization_id, project_id, kind) VALUES (1, 4, 0, 'OrganizationTemplate'), (2, 4, 5, 'Project')",
		"INSERT INTO project_management_workflow_stages (id, workflow_id, name, phase, sort_order) VALUES (10, 1, 'Pending', 'Pending', 1), (11, 1, 'In progress', 'InProgress', 2), (12, 1, 'Finished', 'Finished', 3)",
		"INSERT INTO project_management_workflow_stages (id, workflow_id, name, phase, sort_order) VALUES (20, 2, 'Project pending', 'Pending', 1), (21, 2, 'Project done', 'Finished', 2)",
		"INSERT INTO project_management_workflow_transitions (id, workflow_id, name, from_stage_id, to_stage_id) VALUES (2, 1, 'Pause', 10, 11), (1, 1, 'Finish', 10, 12), (3, 1, 'Finish work', 11, 12), (4, 2, 'Complete', 20, 21)",
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
	definition, err := NewProjectManagementWorkflowRepo(db).FindDefinitionByOrganization(context.Background(), 1, 4)
	if err != nil || len(definition.Stages) != 3 || len(definition.Stages[0].Transitions) != 2 || definition.Stages[0].Transitions[0].FromStage == nil || definition.Stages[0].Transitions[0].ToStage == nil {
		t.Fatalf("unexpected template workflow stages: %#v, %v", definition, err)
	}
	workflow, err := NewProjectManagementWorkflowRepo(db).FindProjectDefinition(context.Background(), 5, 4)
	if err != nil {
		t.Fatalf("find project workflow: %v", err)
	}
	if len(workflow.Stages) != 2 || len(workflow.Stages[0].Transitions) != 1 || workflow.Stages[0].Transitions[0].FromStage == nil || workflow.Stages[0].Transitions[0].ToStage == nil || workflow.Stages[0].Transitions[0].FromStage.Name != "Project pending" || workflow.Stages[0].Transitions[0].ToStage.Name != "Project done" {
		t.Fatalf("unexpected project workflow stages: %#v, %v", workflow, err)
	}
	entity, err := copierx.Trans[model.ProjectManagementWorkflow, proto.ProjectManagementWorkflowEntity](workflow)
	if err != nil || len(entity.Stages) != 2 || len(entity.Stages[0].Transitions) != 1 || entity.Stages[0].Transitions[0].ToStage == nil {
		t.Fatalf("unexpected workflow entity stages: %#v, %v", entity, err)
	}
}

func TestWorkflowLocksRejectMismatchedKinds(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Exec("CREATE TABLE project_management_workflows (id integer primary key, is_del integer not null default 0, organization_id integer not null, project_id integer not null default 0, kind text not null)").Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Exec("INSERT INTO project_management_workflows (id, organization_id, project_id, kind) VALUES (1, 4, 0, 'OrganizationTemplate'), (2, 4, 5, 'Project')").Error; err != nil {
		t.Fatal(err)
	}

	repo := NewProjectManagementWorkflowRepo(db)
	err = db.Transaction(func(tx *gorm.DB) error {
		if _, err := repo.LockTemplateWithTx(context.Background(), tx, 2, 4); err == nil {
			t.Fatal("project workflow accepted as organization template")
		}
		if _, err := repo.LockProjectWithTx(context.Background(), tx, 1, 4); err == nil {
			t.Fatal("organization template accepted as project workflow")
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestStageSortOrderUsesWorkflowScope(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	for _, statement := range []string{
		"CREATE TABLE project_management_workflows (id integer primary key, is_del integer not null default 0, organization_id integer not null, kind text not null)",
		"CREATE TABLE project_management_workflow_stages (id integer primary key, created_at datetime, updated_at datetime, is_del integer not null default 0, workflow_id integer not null, name text not null, phase text not null, color text not null default '', sort_order integer not null)",
		"INSERT INTO project_management_workflows (id, organization_id, kind) VALUES (1, 4, 'OrganizationTemplate'), (2, 4, 'OrganizationTemplate')",
		"INSERT INTO project_management_workflow_stages (id, workflow_id, name, phase, sort_order) VALUES (10, 1, 'Pending', 'Pending', 10), (11, 1, 'Finished', 'Finished', 20), (12, 2, 'Other', 'Pending', 30)",
	} {
		if err := db.Exec(statement).Error; err != nil {
			t.Fatal(err)
		}
	}
	repo := NewProjectManagementWorkflowRepo(db)
	if err := db.Transaction(func(tx *gorm.DB) error {
		maxSortOrder, err := repo.MaxStageSortOrderWithTx(context.Background(), tx, 1)
		if err != nil || maxSortOrder != 20 {
			t.Fatalf("unexpected max sort order: %d, %v", maxSortOrder, err)
		}
		return repo.SwapTemplateStageSortOrdersWithTx(context.Background(), tx, 1, 4, 10, 11)
	}); err != nil {
		t.Fatal(err)
	}
	var stages []model.ProjectManagementWorkflowStage
	if err := db.Where("workflow_id = ?", 1).Order("id ASC").Find(&stages).Error; err != nil || stages[0].SortOrder != 20 || stages[1].SortOrder != 10 {
		t.Fatalf("unexpected swapped stages: %#v, %v", stages, err)
	}
	if err := db.Transaction(func(tx *gorm.DB) error {
		return repo.SwapTemplateStageSortOrdersWithTx(context.Background(), tx, 1, 4, 10, 12)
	}); err == nil {
		t.Fatal("cross-workflow stage swap accepted")
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
