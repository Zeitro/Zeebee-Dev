package logic

import (
	"context"
	"fmt"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/repo"
)

func TestValidateProjectManagementProjectWorkflow(t *testing.T) {
	workflow := dto.ProjectManagementWorkflowDto{
		Name: "Project workflow",
		Stages: []dto.ProjectManagementWorkflowStageDto{
			{Id: 12, Name: "Pending", Phase: constants.ProjectManagementStagePhasePending, SortOrder: 10},
			{Id: -1, Name: "In progress", Phase: constants.ProjectManagementStagePhaseInProgress, SortOrder: 20},
			{Id: -2, Name: "Finished", Phase: constants.ProjectManagementStagePhaseFinished, SortOrder: 30},
		},
		Transitions: []dto.ProjectManagementWorkflowTransitionDto{{Name: "Start", FromStageId: 12, ToStageId: -1}},
	}
	if err := validateProjectManagementProjectWorkflow(workflow); err != nil {
		t.Fatalf("valid workflow replacement rejected: %v", err)
	}
	workflow.Name = ""
	if err := validateProjectManagementProjectWorkflow(workflow); err != nil {
		t.Fatalf("project workflow definition should not require a name: %v", err)
	}
	workflow.Transitions = append(workflow.Transitions, dto.ProjectManagementWorkflowTransitionDto{Name: "Restart", FromStageId: -1, ToStageId: 12})
	if err := validateProjectManagementProjectWorkflow(workflow); err != nil {
		t.Fatalf("workflow replacement rejected a transition loop: %v", err)
	}
	workflow.Transitions = workflow.Transitions[:1]
	workflow.Stages[1].Id = 0
	if err := validateProjectManagementProjectWorkflow(workflow); err == nil {
		t.Fatal("workflow replacement accepted a stage without a stable or temporary id")
	}
}

func TestUpdateProjectWorkflowRequiresProjectManagerAndPreservesMetadata(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())), &gorm.Config{})
	require.NoError(t, err)
	for _, statement := range []string{
		"CREATE TABLE project_management_projects (id INTEGER PRIMARY KEY, created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, is_del INTEGER NOT NULL DEFAULT 0, organization_id INTEGER NOT NULL, entity_id INTEGER NOT NULL DEFAULT 0, entity_type TEXT NOT NULL DEFAULT '', name TEXT NOT NULL DEFAULT '', description TEXT NOT NULL DEFAULT '', cover_url TEXT NOT NULL DEFAULT '', owner_id INTEGER NOT NULL DEFAULT 0, owner_name TEXT NOT NULL DEFAULT '')",
		"CREATE TABLE project_management_project_members (id INTEGER PRIMARY KEY, created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, is_del INTEGER NOT NULL DEFAULT 0, project_id INTEGER NOT NULL, organization_user_id INTEGER NOT NULL, role TEXT NOT NULL)",
		"CREATE TABLE project_management_workflows (id INTEGER PRIMARY KEY, created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, is_del INTEGER NOT NULL DEFAULT 0, kind TEXT NOT NULL, organization_id INTEGER NOT NULL, project_id INTEGER NOT NULL DEFAULT 0, name TEXT NOT NULL DEFAULT '', description TEXT NOT NULL DEFAULT '', is_default BOOLEAN NOT NULL DEFAULT FALSE, source_template_id INTEGER NOT NULL DEFAULT 0, source_template_name TEXT NOT NULL DEFAULT '')",
		"CREATE TABLE project_management_workflow_stages (id INTEGER PRIMARY KEY, created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, is_del INTEGER NOT NULL DEFAULT 0, workflow_id INTEGER NOT NULL, name TEXT NOT NULL DEFAULT '', phase TEXT NOT NULL, color TEXT NOT NULL DEFAULT '', sort_order INTEGER NOT NULL)",
		"CREATE TABLE project_management_workflow_transitions (id INTEGER PRIMARY KEY, created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, is_del INTEGER NOT NULL DEFAULT 0, workflow_id INTEGER NOT NULL, name TEXT NOT NULL DEFAULT '', from_stage_id INTEGER NOT NULL, to_stage_id INTEGER NOT NULL, allowed_actor_types JSON NOT NULL DEFAULT '[]')",
		"CREATE TABLE operations (id INTEGER PRIMARY KEY AUTOINCREMENT, created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, is_del INTEGER NOT NULL DEFAULT 0, entity_id INTEGER NOT NULL, entity_type TEXT NOT NULL, operator_id INTEGER NOT NULL, operator_name TEXT NOT NULL DEFAULT '', operator_type TEXT NOT NULL, operation TEXT NOT NULL, note TEXT NOT NULL DEFAULT '', metadata JSON NOT NULL DEFAULT '{}', organization_id INTEGER NOT NULL)",
	} {
		require.NoError(t, db.Exec(statement).Error)
	}
	require.NoError(t, db.Exec("CREATE TABLE project_management_tasks (id INTEGER PRIMARY KEY, updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, is_del INTEGER NOT NULL DEFAULT 0, stage_id INTEGER NOT NULL, phase TEXT NOT NULL DEFAULT 'Pending')").Error)
	require.NoError(t, db.Exec("INSERT INTO project_management_projects (id, organization_id, name, owner_id) VALUES (1, 1, 'Project', 10)").Error)
	require.NoError(t, db.Exec("INSERT INTO project_management_project_members (id, project_id, organization_user_id, role) VALUES (1, 1, 20, ?), (2, 1, 10, ?)", constants.ProjectManagementMemberRoleTaskExecutor, constants.ProjectManagementMemberRoleProjectOwner).Error)
	require.NoError(t, db.Exec("INSERT INTO project_management_workflows (id, kind, organization_id, project_id, name, description) VALUES (100, ?, 1, 1, 'Workflow', 'Original description')", constants.ProjectManagementWorkflowKindProject).Error)
	require.NoError(t, db.Exec("INSERT INTO project_management_workflow_stages (id, workflow_id, name, phase, sort_order) VALUES (101, 100, 'Pending', ?, 1), (102, 100, 'Finished', ?, 2)", constants.ProjectManagementStagePhasePending, constants.ProjectManagementStagePhaseFinished).Error)

	logic := NewProjectManagementWorkflowLogic(NewBaseLogic(repo.NewAllRepos(db), nil, nil, dto.JwtConfig{}))
	input := dto.ProjectManagementWorkflowUpdateDto{
		Stages: []dto.ProjectManagementWorkflowUpdateStageDto{
			{Id: 101, Name: "Not started", Phase: constants.ProjectManagementStagePhasePending, SortOrder: 1},
			{Id: 102, Name: "Finished", Phase: constants.ProjectManagementStagePhaseFinished, SortOrder: 2},
		},
	}
	_, err = logic.ProjectUpdate(context.Background(), utils.JwtUser{OrganizationId: 1, OrganizationUserId: 20}, 1, input)
	require.Error(t, err)
	require.Contains(t, err.Error(), "project member role is not allowed")

	var workflow model.ProjectManagementWorkflow
	require.NoError(t, db.First(&workflow, 100).Error)
	require.Equal(t, "Workflow", workflow.Name)

	updated, err := logic.ProjectUpdate(context.Background(), utils.JwtUser{OrganizationId: 1, OrganizationUserId: 10}, 1, input)
	require.NoError(t, err)
	require.Equal(t, "Workflow", updated.Name)
	require.Equal(t, "Original description", updated.Description)
	require.NoError(t, db.First(&workflow, 100).Error)
	require.Equal(t, "Workflow", workflow.Name)
	require.Equal(t, "Original description", workflow.Description)
	var stage model.ProjectManagementWorkflowStage
	require.NoError(t, db.First(&stage, 101).Error)
	require.Equal(t, "Not started", stage.Name)
}
