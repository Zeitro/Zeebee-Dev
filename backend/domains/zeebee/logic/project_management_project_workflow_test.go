package logic

import (
	"context"
	"fmt"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

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
		},
		Transitions: []dto.ProjectManagementWorkflowTransitionDto{{FromStageId: 12, ToStageId: -1}},
	}
	if err := validateProjectManagementProjectWorkflow(workflow); err != nil {
		t.Fatalf("valid workflow replacement rejected: %v", err)
	}
	workflow.Stages[1].Id = 0
	if err := validateProjectManagementProjectWorkflow(workflow); err == nil {
		t.Fatal("workflow replacement accepted a stage without a stable or temporary id")
	}
}

func TestUpdateProjectWorkflowRequiresProjectManager(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())), &gorm.Config{})
	require.NoError(t, err)
	for _, statement := range []string{
		"CREATE TABLE project_management_projects (id INTEGER PRIMARY KEY, created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, is_del INTEGER NOT NULL DEFAULT 0, organization_id INTEGER NOT NULL, property_id INTEGER NOT NULL DEFAULT 0, name TEXT NOT NULL DEFAULT '', description TEXT NOT NULL DEFAULT '', cover_url TEXT NOT NULL DEFAULT '', owner_id INTEGER NOT NULL DEFAULT 0)",
		"CREATE TABLE project_management_project_members (id INTEGER PRIMARY KEY, created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, is_del INTEGER NOT NULL DEFAULT 0, project_id INTEGER NOT NULL, organization_user_id INTEGER NOT NULL, role TEXT NOT NULL)",
		"CREATE TABLE project_management_workflows (id INTEGER PRIMARY KEY, created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, is_del INTEGER NOT NULL DEFAULT 0, kind TEXT NOT NULL, organization_id INTEGER NOT NULL, project_id INTEGER NOT NULL DEFAULT 0, name TEXT NOT NULL DEFAULT '', description TEXT NOT NULL DEFAULT '', enabled BOOLEAN NOT NULL DEFAULT TRUE, is_default BOOLEAN NOT NULL DEFAULT FALSE, source_template_id INTEGER NOT NULL DEFAULT 0, source_template_name TEXT NOT NULL DEFAULT '')",
		"CREATE TABLE project_management_workflow_stages (id INTEGER PRIMARY KEY, created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, is_del INTEGER NOT NULL DEFAULT 0, workflow_id INTEGER NOT NULL, name TEXT NOT NULL DEFAULT '', phase TEXT NOT NULL, sort_order INTEGER NOT NULL)",
		"CREATE TABLE project_management_workflow_transitions (id INTEGER PRIMARY KEY, created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, is_del INTEGER NOT NULL DEFAULT 0, workflow_id INTEGER NOT NULL, from_stage_id INTEGER NOT NULL, to_stage_id INTEGER NOT NULL, allowed_actor_types TEXT NOT NULL DEFAULT '{}')",
	} {
		require.NoError(t, db.Exec(statement).Error)
	}
	require.NoError(t, db.Exec("INSERT INTO project_management_projects (id, organization_id, name, owner_id) VALUES (1, 1, 'Project', 10)").Error)
	require.NoError(t, db.Exec("INSERT INTO project_management_project_members (id, project_id, organization_user_id, role) VALUES (1, 1, 20, ?)", constants.ProjectManagementMemberRoleTaskExecutor).Error)
	require.NoError(t, db.Exec("INSERT INTO project_management_workflows (id, kind, organization_id, project_id, name) VALUES (100, ?, 1, 1, 'Workflow')", constants.ProjectManagementWorkflowKindProject).Error)
	require.NoError(t, db.Exec("INSERT INTO project_management_workflow_stages (id, workflow_id, name, phase, sort_order) VALUES (101, 100, 'Pending', ?, 1)", constants.ProjectManagementStagePhasePending).Error)

	logic := NewProjectManagementLogic(NewBaseLogic(repo.NewAllRepos(db), nil, nil, dto.JwtConfig{}))
	input := dto.ProjectManagementWorkflowDto{
		Name:   "Updated workflow",
		Stages: []dto.ProjectManagementWorkflowStageDto{{Id: 101, Name: "Pending", Phase: constants.ProjectManagementStagePhasePending, SortOrder: 1}},
	}
	_, err = logic.UpdateProjectWorkflow(context.Background(), &dto.ProjectManagementActor{OrganizationId: 1, OrganizationUserId: 20}, 1, input)
	require.Error(t, err)
	require.Contains(t, err.Error(), "project member role is not allowed")

	var workflow model.ProjectManagementWorkflow
	require.NoError(t, db.First(&workflow, 100).Error)
	require.Equal(t, "Workflow", workflow.Name)
}
