package repo

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/proto"
)

func TestProjectManagementWorkflowRepoListOmitsTransitions(t *testing.T) {
	db := newRepoTestDB(t)
	execTestSQL(t, db, []string{
		"CREATE TABLE project_management_templates (id integer primary key, is_del integer not null default 0, template_type text not null, organization_id integer not null default 0, status text not null default 'Active')",
		"CREATE TABLE project_management_projects (id integer primary key, is_del integer not null default 0, organization_id integer not null default 0)",
		"CREATE TABLE project_management_workflows (id integer primary key autoincrement, created_at datetime, updated_at datetime, is_del integer not null default 0, name text not null default '', description text not null default '', is_default boolean not null default false, entity_type text not null default 'Project', entity_id integer not null default 0, source_workflow_id integer not null default 0, source_workflow_name text not null default '')",
		"CREATE TABLE project_management_stages (id integer primary key, is_del integer not null default 0, sort_order integer not null default 0)",
		"CREATE TABLE project_management_workflow_stages (id integer primary key, is_del integer not null default 0, workflow_id integer not null, stage_id integer not null)",
		"CREATE TABLE project_management_workflow_transitions (id integer primary key autoincrement, created_at datetime, updated_at datetime, is_del integer not null default 0, workflow_id integer not null default 0, name text not null default '', from_stage_id integer not null default 0, to_stage_id integer not null default 0, allowed_actor_types json not null default '[]')",
		"INSERT INTO project_management_templates (id, template_type, organization_id, status) VALUES (1, 'Organization', 1, 'Active')",
		"INSERT INTO project_management_projects (id, organization_id) VALUES (1, 1)",
	})
	workflows := []*model.ProjectManagementWorkflow{
		{EntityType: constants.ProjectManagementOwnerTypeTemplate, EntityId: 1, Name: "Source"},
		{EntityType: constants.ProjectManagementOwnerTypeProject, EntityId: 1, Name: "Project"},
	}
	require.NoError(t, db.Create(workflows).Error)
	require.NoError(t, db.Create(&model.ProjectManagementWorkflowTransition{WorkflowId: workflows[0].Id, Name: "Hidden"}).Error)
	require.NoError(t, db.Create(&model.ProjectManagementWorkflowTransition{WorkflowId: workflows[1].Id, Name: "Hidden"}).Error)

	repo := NewProjectManagementWorkflowRepo(db)
	sources, err := repo.ListWorkflowSources(t.Context(), 1)
	require.NoError(t, err)
	require.Len(t, sources, 1)
	require.Nil(t, sources[0].Transitions)
	templates, err := repo.ListTemplateWorkflows(t.Context(), 1)
	require.NoError(t, err)
	require.Len(t, templates, 1)
	require.Nil(t, templates[0].Transitions)
	projects, err := repo.ListProjectWorkflows(t.Context(), 1, 1)
	require.NoError(t, err)
	require.Len(t, projects, 1)
	require.Nil(t, projects[0].Transitions)
}

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
	if err := db.Create(&model.ProjectManagementWorkflowTransition{WorkflowId: workflow.Id, Name: "Finish", FromStageId: stages[1].Id, ToStageId: stages[0].Id}).Error; err != nil {
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
	if got.Transitions[0].FromStage == nil || got.Transitions[0].FromStage.Name != "Pending" || got.Transitions[0].ToStage == nil || got.Transitions[0].ToStage.Name != "Finished" {
		t.Fatalf("transition stages = %#v", got.Transitions[0])
	}
	stage := new(model.ProjectManagementStage)
	if err := db.Preload("Transitions").First(stage, stages[1].Id).Error; err != nil {
		t.Fatal(err)
	}
	if len(stage.Transitions) != 1 || stage.Transitions[0].ToStageId != stages[0].Id {
		t.Fatalf("stage transitions = %#v", stage.Transitions)
	}
	entity, err := copierx.Trans[model.ProjectManagementWorkflow, proto.ProjectManagementWorkflowEntity](got)
	if err != nil || entity.Transitions[0].FromStage.GetName() != "Pending" || entity.Transitions[0].ToStage.GetName() != "Finished" {
		t.Fatalf("transition entity = %#v, %v", entity.Transitions[0], err)
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
