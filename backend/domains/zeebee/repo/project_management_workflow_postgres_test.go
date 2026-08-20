package repo_test

import (
	"context"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	domainlogic "zeebee-backend/domains/zeebee/logic"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/repo"
)

func projectManagementPostgresDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := os.Getenv("ZEEBEE_TEST_POSTGRES_DSN")
	if dsn == "" {
		t.Skip("set ZEEBEE_TEST_POSTGRES_DSN to a migrated zeebee_pm_test database")
	}
	if !strings.Contains(dsn, "zeebee_pm_test") {
		t.Fatal("ZEEBEE_TEST_POSTGRES_DSN must target a zeebee_pm_test database")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.Exec(`CREATE TABLE IF NOT EXISTS operations (
		id bigserial PRIMARY KEY, created_at timestamptz NOT NULL DEFAULT now(),
		updated_at timestamptz NOT NULL DEFAULT now(), is_del int2 NOT NULL DEFAULT 0,
		entity_id int4 NOT NULL DEFAULT 0, entity_type varchar(100) NOT NULL DEFAULT '',
		operator_id int4 NOT NULL DEFAULT 0, operator_name varchar(50) NOT NULL DEFAULT '',
		operator_type varchar(50) NOT NULL DEFAULT '', operation varchar(50) NOT NULL DEFAULT '',
		note varchar(255) NOT NULL DEFAULT '', metadata jsonb NOT NULL DEFAULT '{}',
		organization_id int4 NOT NULL DEFAULT 0)`).Error)
	require.NoError(t, db.Exec(`TRUNCATE operations, project_management_task_dependencies,
		project_management_task_assignees, project_management_tasks, project_management_project_members,
		project_management_projects, project_management_workflow_transitions,
		project_management_workflow_stages, project_management_workflows RESTART IDENTITY`).Error)
	return db
}

func TestProjectManagementWorkflowPostgresConcurrencyAndConstraints(t *testing.T) {
	db := projectManagementPostgresDB(t)
	projectLogic := domainlogic.NewProjectManagementWorkflowLogic(
		domainlogic.NewBaseLogic(repo.NewAllRepos(db), nil, nil, dto.JwtConfig{}),
	)
	var wg sync.WaitGroup
	errs := make(chan error, 2)
	for range 2 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := projectLogic.ListEnabled(context.Background(), 1)
			errs <- err
		}()
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		require.NoError(t, err)
	}

	var count int64
	require.NoError(t, db.Model(&model.ProjectManagementWorkflow{}).
		Where("organization_id = ? AND kind = ? AND is_default = true", 1, constants.ProjectManagementWorkflowKindOrganizationTemplate).
		Count(&count).Error)
	require.EqualValues(t, 1, count)

	workflow := &model.ProjectManagementWorkflow{Kind: constants.ProjectManagementWorkflowKindOrganizationTemplate, OrganizationId: 2, Name: "Constraint workflow", Enabled: true}
	require.NoError(t, db.Create(workflow).Error)
	pending := &model.ProjectManagementWorkflowStage{WorkflowId: workflow.Id, Name: "Pending", Phase: constants.ProjectManagementStagePhasePending, SortOrder: 10}
	inProgress := &model.ProjectManagementWorkflowStage{WorkflowId: workflow.Id, Name: "In progress", Phase: constants.ProjectManagementStagePhaseInProgress, SortOrder: 20}
	require.NoError(t, db.Create(pending).Error)
	require.NoError(t, db.Create(inProgress).Error)
	require.Error(t, db.Create(&model.ProjectManagementWorkflowStage{WorkflowId: workflow.Id, Name: "Other pending", Phase: constants.ProjectManagementStagePhasePending, SortOrder: 30}).Error)
	require.Error(t, db.Create(&model.ProjectManagementWorkflowStage{WorkflowId: workflow.Id, Name: "Same order", Phase: constants.ProjectManagementStagePhaseFinished, SortOrder: 20}).Error)

	transition := &model.ProjectManagementWorkflowTransition{WorkflowId: workflow.Id, FromStageId: pending.Id, ToStageId: inProgress.Id}
	require.NoError(t, db.Create(transition).Error)
	require.Error(t, db.Create(&model.ProjectManagementWorkflowTransition{WorkflowId: workflow.Id, FromStageId: pending.Id, ToStageId: inProgress.Id}).Error)
}
