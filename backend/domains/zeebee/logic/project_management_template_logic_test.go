package logic

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"zeebee-backend/common/copierx"
	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/repo"
)

func TestProjectManagementTemplateLifecycleAndDiscovery(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.Exec("CREATE TABLE project_management_templates (id integer primary key autoincrement, created_at datetime, updated_at datetime, is_del integer not null default 0, template_type text not null, organization_id integer not null default 0, name text not null, description text not null default '', status text not null default 'Active', is_default boolean not null default false, sort_order integer not null default 0, project_type text, start_at datetime, end_at datetime, period_value integer, period_unit text, source_project_id integer not null default 0)").Error)
	require.NoError(t, db.Exec("CREATE TABLE project_management_workflows (id integer primary key autoincrement, created_at datetime, updated_at datetime, is_del integer not null default 0, entity_type text not null, entity_id integer not null, source_workflow_id integer not null default 0, source_workflow_name text not null default '', name text not null default '', description text not null default '', is_default boolean not null default false)").Error)
	require.NoError(t, db.Exec("CREATE TABLE project_management_stages (id integer primary key autoincrement, created_at datetime, updated_at datetime, is_del integer not null default 0, entity_type text not null, entity_id integer not null, name text not null default '', description text not null default '', phase text not null, color text not null default '', sort_order integer not null)").Error)
	require.NoError(t, db.Exec("CREATE TABLE project_management_workflow_stages (id integer primary key autoincrement, created_at datetime, updated_at datetime, is_del integer not null default 0, workflow_id integer not null, stage_id integer not null)").Error)
	require.NoError(t, db.Exec("CREATE TABLE project_management_workflow_transitions (id integer primary key autoincrement, created_at datetime, updated_at datetime, is_del integer not null default 0, workflow_id integer not null, name text not null default '', from_stage_id integer not null, to_stage_id integer not null, allowed_actor_types json not null default '[]')").Error)
	logic := NewProjectManagementTemplateLogic(NewBaseLogic(repo.NewAllRepos(db), nil, nil, dto.JwtConfig{}))
	adminClaim := utils.JwtUser{UserType: constants.UserTypeAdmin}
	start := &sql.NullTime{Time: time.Date(2026, time.January, 31, 0, 0, 0, 0, time.UTC), Valid: true}
	end := &sql.NullTime{Time: time.Date(2026, time.January, 31, 0, 0, 0, 0, time.UTC), Valid: true}
	period := 1
	created, err := logic.Create(context.Background(), utils.JwtUser{OrganizationId: 9}, dto.ProjectManagementTemplateCreateDto{
		Name: "Public", Status: "Disabled", ProjectType: constants.ProjectManagementProjectTypePlanned, StartAt: start, EndAt: end, PeriodValue: period, PeriodUnit: constants.ProjectManagementPeriodUnitMonth,
	})
	require.NoError(t, err)
	require.Equal(t, constants.CommonStatusActive, created.Status)
	require.Equal(t, "2026-01-31", created.EndAt.Time.Format(time.DateOnly))
	stages, err := logic.Stages(t.Context(), adminClaim, created.Id)
	require.NoError(t, err)
	inputStages, err := copierx.Trans[[]*model.ProjectManagementStage, []dto.ProjectManagementTemplateStageDto](&stages)
	require.NoError(t, err)
	sortOrder := 12
	updated, err := logic.Update(t.Context(), adminClaim, created.Id, dto.ProjectManagementTemplateUpdateDto{Name: "Renamed", Description: "Updated", SortOrder: sortOrder, ProjectType: constants.ProjectManagementProjectTypePlanned, StartAt: start, EndAt: end, PeriodValue: period, PeriodUnit: constants.ProjectManagementPeriodUnitMonth, Stages: *inputStages})
	require.NoError(t, err)
	require.Equal(t, "Renamed", updated.Name)
	require.Equal(t, "Updated", updated.Description)
	require.Equal(t, 12, updated.SortOrder)
	require.Equal(t, (*inputStages)[0].Id, updated.Stages[0].Id)
	_, err = logic.Update(t.Context(), adminClaim, created.Id, dto.ProjectManagementTemplateUpdateDto{ProjectType: constants.ProjectManagementProjectTypeOrdinary, Stages: *inputStages})
	require.ErrorContains(t, err, "project type cannot be changed")
	persisted, err := logic.Detail(t.Context(), adminClaim, created.Id)
	require.NoError(t, err)
	require.Equal(t, constants.ProjectManagementProjectTypePlanned, persisted.ProjectType)
	require.NoError(t, db.Create(&model.ProjectManagementStage{EntityType: constants.ProjectManagementOwnerTypeTemplate, EntityId: created.Id, Name: "Obsolete", Phase: constants.ProjectManagementStagePhaseInProgress, SortOrder: 99}).Error)
	replacementStages := append(*inputStages, dto.ProjectManagementTemplateStageDto{Name: "Obsolete", Phase: constants.ProjectManagementStagePhaseInProgress, SortOrder: 99})
	require.NoError(t, logic.repos.ProjectManagementTemplateRepo.Transaction(t.Context(), func(tx *gorm.DB) error {
		return logic.syncTemplateStagesWithTx(t.Context(), tx, created.Id, replacementStages)
	}))
	stages, err = logic.Stages(t.Context(), adminClaim, created.Id)
	require.NoError(t, err)
	require.Len(t, stages, 4)
	paged, err := logic.Page(t.Context(), dto.ProjectManagementTemplatePageRequestDto{BasePageRequestDto: dto.BasePageRequestDto{Page: 1, Limit: 20}})
	require.NoError(t, err)
	require.Len(t, paged.List, 1)
	require.Len(t, paged.List[0].Stages, 4)
	workflows, err := logic.Workflows(context.Background(), adminClaim, created.Id)
	require.NoError(t, err)
	require.Len(t, workflows, 1)
	require.True(t, workflows[0].IsDefault)
	require.Len(t, workflows[0].Stages, 3)
	require.Len(t, workflows[0].Transitions, 2)
	transitions, err := copierx.Trans[[]*model.ProjectManagementWorkflowStageTransition, []dto.ProjectManagementWorkflowTransitionDto](&workflows[0].Transitions)
	require.NoError(t, err)
	(*transitions)[0].Name = "Renamed transition"
	workflow, err := logic.UpdateWorkflow(t.Context(), adminClaim, created.Id, workflows[0].Id, dto.ProjectManagementWorkflowDto{Name: workflows[0].Name, Description: workflows[0].Description, Transitions: *transitions})
	require.NoError(t, err)
	require.Equal(t, workflows[0].Transitions[0].Id, workflow.Transitions[0].Id)
	require.Equal(t, "Renamed transition", workflow.Transitions[0].Name)
	orgTemplate := &model.ProjectManagementTemplate{TemplateType: constants.ProjectManagementTemplateTypeOrganization, OrganizationId: 9, Name: "Organization", Status: "Active", ProjectType: constants.ProjectManagementProjectTypeOrdinary}
	require.NoError(t, db.Create(orgTemplate).Error)
	foreignTemplate := &model.ProjectManagementTemplate{TemplateType: constants.ProjectManagementTemplateTypeOrganization, OrganizationId: 10, Name: "Foreign", Status: "Active", ProjectType: constants.ProjectManagementProjectTypeOrdinary}
	require.NoError(t, db.Create(foreignTemplate).Error)
	paged, err = logic.Page(t.Context(), dto.ProjectManagementTemplatePageRequestDto{BasePageRequestDto: dto.BasePageRequestDto{Page: 1, Limit: 20}, Filter: dto.ProjectManagementTemplatePageRequestFilter{OrganizationId: 9}})
	require.NoError(t, err)
	require.Equal(t, []int64{orgTemplate.Id}, []int64{paged.List[0].Id})
	discovered, err := logic.Discover(context.Background(), utils.JwtUser{OrganizationId: 9})
	require.NoError(t, err)
	require.Len(t, discovered, 2)
	_, err = logic.SetStatus(context.Background(), adminClaim, created.Id, "Disabled")
	require.NoError(t, err)
	pageFilter := dto.ProjectManagementTemplatePageRequestFilter{OrganizationId: 9, Status: "Active"}
	// ponytail: SQLite does not support ILIKE; PostgreSQL integration tests cover keyword matching.
	if db.Dialector.Name() == "postgres" {
		pageFilter.Keyword = "organization"
	}
	paged, err = logic.Page(t.Context(), dto.ProjectManagementTemplatePageRequestDto{BasePageRequestDto: dto.BasePageRequestDto{Page: 1, Limit: 20}, Filter: pageFilter})
	require.NoError(t, err)
	require.Equal(t, []int64{orgTemplate.Id}, []int64{paged.List[0].Id})
	_, err = logic.Detail(context.Background(), adminClaim, foreignTemplate.Id)
	require.Error(t, err)
	discovered, err = logic.Discover(context.Background(), utils.JwtUser{OrganizationId: 9})
	require.NoError(t, err)
	require.Len(t, discovered, 1)
	require.Equal(t, "Organization", discovered[0].Name)
	other := &model.ProjectManagementTemplate{TemplateType: constants.ProjectManagementTemplateTypePublic, Name: "Other", Status: "Active", ProjectType: constants.ProjectManagementProjectTypeOrdinary}
	require.NoError(t, db.Create(other).Error)
	require.NoError(t, logic.SetDefault(t.Context(), adminClaim, created.Id))
	_, err = logic.SetStatus(t.Context(), adminClaim, created.Id, "Disabled")
	require.ErrorContains(t, err, "default project template cannot be disabled")
	_, err = logic.Update(t.Context(), adminClaim, created.Id, dto.ProjectManagementTemplateUpdateDto{Status: "Disabled", ProjectType: constants.ProjectManagementProjectTypePlanned, StartAt: start, EndAt: end, PeriodValue: period, PeriodUnit: constants.ProjectManagementPeriodUnitMonth, Stages: *inputStages})
	require.ErrorContains(t, err, "default project template cannot be disabled")
	require.NoError(t, logic.SetDefault(t.Context(), adminClaim, other.Id))
	var defaults []model.ProjectManagementTemplate
	require.NoError(t, db.Where("template_type = ? AND organization_id = ? AND is_default = ?", constants.ProjectManagementTemplateTypePublic, 0, true).Find(&defaults).Error)
	require.Len(t, defaults, 1)
	require.Equal(t, other.Id, defaults[0].Id)
	lenderClaim := utils.JwtUser{OrganizationId: 9, UserType: constants.UserTypeLender}
	_, err = logic.SetStatus(t.Context(), lenderClaim, orgTemplate.Id, constants.CommonStatusDisabled)
	require.NoError(t, err)
	_, err = logic.SetStatus(t.Context(), lenderClaim, foreignTemplate.Id, constants.CommonStatusDisabled)
	require.ErrorContains(t, err, "project template not found")
	require.Error(t, logic.Delete(t.Context(), lenderClaim, created.Id))
	require.NoError(t, logic.Delete(t.Context(), lenderClaim, orgTemplate.Id))
}

func TestValidateProjectManagementTemplateStages(t *testing.T) {
	stages := []dto.ProjectManagementTemplateStageDto{
		{Id: 1, Name: "Pending", Phase: constants.ProjectManagementStagePhasePending},
		{Id: 2, Name: "Finished", Phase: constants.ProjectManagementStagePhaseFinished},
	}
	require.NoError(t, validateProjectManagementTemplateStages(stages))
	stages[1].Name = "Pending"
	require.Error(t, validateProjectManagementTemplateStages(stages))
	stages[1].Name = "Finished"
	stages[1].Phase = constants.ProjectManagementStagePhaseInProgress
	require.Error(t, validateProjectManagementTemplateStages(stages))
}

func TestTemplateWorkflowStageIDs(t *testing.T) {
	ids := templateWorkflowStageIDs([]dto.ProjectManagementWorkflowTransitionDto{
		{FromStageId: 2, ToStageId: 3},
		{FromStageId: 3, ToStageId: 5},
	})
	if len(ids) != 3 || ids[0] != 2 || ids[1] != 3 || ids[2] != 5 {
		t.Fatalf("stage ids = %v", ids)
	}
}

func TestProjectManagementTemplateCreateFromProjectSnapshotsWorkflows(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())), &gorm.Config{})
	require.NoError(t, err)
	for _, statement := range []string{
		"CREATE TABLE project_management_templates (id integer primary key autoincrement, created_at datetime, updated_at datetime, is_del integer not null default 0, template_type text not null, organization_id integer not null default 0, name text not null, description text not null default '', status text not null default 'Active', is_default boolean not null default false, sort_order integer not null default 0, project_type text, start_at datetime, end_at datetime, period_value integer, period_unit text, source_project_id integer not null default 0)",
		"CREATE TABLE project_management_projects (id integer primary key, created_at datetime, updated_at datetime, is_del integer not null default 0, organization_id integer not null, entity_id integer not null default 0, entity_type text not null default '', name text not null default '', description text not null default '', priority text not null default '', cover_url text not null default '', owner_id integer not null default 0, owner_name text not null default '', project_type text not null default 'Ordinary', start_at datetime, end_at datetime, period_value integer not null default 0, period_unit text not null default '')",
		"CREATE TABLE project_management_project_members (id integer primary key, created_at datetime, updated_at datetime, is_del integer not null default 0, project_id integer not null, organization_user_id integer not null, role text not null)",
		"CREATE TABLE project_management_workflows (id integer primary key autoincrement, created_at datetime, updated_at datetime, is_del integer not null default 0, entity_type text not null, entity_id integer not null, source_workflow_id integer not null default 0, source_workflow_name text not null default '', name text not null default '', description text not null default '', is_default boolean not null default false)",
		"CREATE TABLE project_management_stages (id integer primary key autoincrement, created_at datetime, updated_at datetime, is_del integer not null default 0, entity_type text not null, entity_id integer not null, name text not null default '', description text not null default '', phase text not null, color text not null default '', sort_order integer not null)",
		"CREATE TABLE project_management_workflow_stages (id integer primary key autoincrement, created_at datetime, updated_at datetime, is_del integer not null default 0, workflow_id integer not null, stage_id integer not null)",
		"CREATE TABLE project_management_workflow_transitions (id integer primary key autoincrement, created_at datetime, updated_at datetime, is_del integer not null default 0, workflow_id integer not null, name text not null default '', from_stage_id integer not null, to_stage_id integer not null, allowed_actor_types json not null default '[]')",
		"CREATE TABLE project_management_tasks (id integer primary key autoincrement, created_at datetime, updated_at datetime, is_del integer not null default 0, entity_type text not null, entity_id integer not null, workflow_id integer not null default 0, organization_id integer not null, parent_task_id integer not null default 0, stage_id integer, phase text, creator_id integer not null default 0, creator_name text not null default '', title text not null default '', description text not null default '', priority text not null default '', work_hours numeric not null default 0, start_at datetime, end_at datetime, task_type json not null default '{}')",
		"CREATE TABLE project_management_task_assignees (id integer primary key autoincrement, created_at datetime, updated_at datetime, is_del integer not null default 0, task_id integer not null, stage_id integer, organization_user_id integer not null, organization_user_full_name text not null default '')",
		"CREATE TABLE project_management_task_dependencies (id integer primary key autoincrement, created_at datetime, updated_at datetime, is_del integer not null default 0, task_id integer not null, dependency_task_id integer not null, dependency_type text not null, dependency_rule text not null)",
	} {
		require.NoError(t, db.Exec(statement).Error)
	}
	require.NoError(t, db.Exec("INSERT INTO project_management_projects (id, organization_id, name, priority, project_type) VALUES (1, 9, 'Source', 'High', 'Ordinary')").Error)
	require.NoError(t, db.Exec("INSERT INTO project_management_project_members (id, project_id, organization_user_id, role) VALUES (1, 1, 7, 'ProjectOwner')").Error)
	require.NoError(t, db.Exec("INSERT INTO project_management_workflows (id, entity_type, entity_id, name, is_default) VALUES (10, 'Project', 1, 'Source flow', true)").Error)
	require.NoError(t, db.Exec("INSERT INTO project_management_stages (id, entity_type, entity_id, name, phase, sort_order) VALUES (20, 'Project', 1, 'Pending', 'Pending', 1), (21, 'Project', 1, 'Finished', 'Finished', 2)").Error)
	require.NoError(t, db.Exec("INSERT INTO project_management_workflow_stages (id, workflow_id, stage_id) VALUES (30, 10, 20), (31, 10, 21)").Error)
	require.NoError(t, db.Exec("INSERT INTO project_management_workflow_transitions (workflow_id, name, from_stage_id, to_stage_id) VALUES (10, 'Finish', 20, 21)").Error)
	require.NoError(t, db.Exec("INSERT INTO project_management_tasks (id, entity_type, entity_id, workflow_id, organization_id, stage_id, phase, title, priority) VALUES (40, 'Project', 1, 10, 9, 20, 'Pending', 'Parent', 'High'), (41, 'Project', 1, 10, 9, 20, 'Pending', 'Child', 'Low')").Error)
	require.NoError(t, db.Exec("UPDATE project_management_tasks SET parent_task_id = 40 WHERE id = 41").Error)
	require.NoError(t, db.Exec("INSERT INTO project_management_task_dependencies (task_id, dependency_task_id, dependency_type, dependency_rule) VALUES (41, 40, 'PreTask', 'FS')").Error)
	require.NoError(t, db.Exec("INSERT INTO project_management_task_assignees (task_id, stage_id, organization_user_id, organization_user_full_name) VALUES (40, 20, 8, 'Executor')").Error)

	logic := NewProjectManagementTemplateLogic(NewBaseLogic(repo.NewAllRepos(db), nil, nil, dto.JwtConfig{}))
	template, err := logic.Create(t.Context(), utils.JwtUser{OrganizationId: 9, OrganizationUserId: 7}, dto.ProjectManagementTemplateCreateDto{Name: "Snapshot", SourceProjectId: 1, CopyTasks: true, CopyAssignees: true})
	require.NoError(t, err)
	require.Equal(t, int64(9), template.OrganizationId)
	workflows, err := logic.Workflows(t.Context(), utils.JwtUser{OrganizationId: 9}, template.Id)
	require.NoError(t, err)
	require.Len(t, workflows, 1)
	require.Equal(t, "Source flow", workflows[0].Name)
	var stageID sql.NullInt64
	var phase sql.NullString
	require.NoError(t, db.Raw("SELECT stage_id, phase FROM project_management_tasks WHERE entity_type = ? AND entity_id = ?", constants.ProjectManagementOwnerTypeTemplate, template.Id).Row().Scan(&stageID, &phase))
	require.False(t, stageID.Valid)
	require.False(t, phase.Valid)
	var assigneeStageID sql.NullInt64
	require.NoError(t, db.Raw("SELECT a.stage_id FROM project_management_task_assignees a JOIN project_management_tasks t ON t.id = a.task_id WHERE t.entity_type = ? AND t.entity_id = ?", constants.ProjectManagementOwnerTypeTemplate, template.Id).Row().Scan(&assigneeStageID))
	require.False(t, assigneeStageID.Valid)
	require.Len(t, workflows[0].Stages, 2)
	require.Len(t, workflows[0].Transitions, 1)
	var tasks []*model.ProjectManagementTask
	require.NoError(t, db.Where("entity_type = ? AND entity_id = ?", constants.ProjectManagementOwnerTypeTemplate, template.Id).Order("id ASC").Find(&tasks).Error)
	require.Len(t, tasks, 2)
	require.Equal(t, tasks[0].Id, tasks[1].ParentTaskId)
	var dependencies []*model.ProjectManagementTaskDependency
	require.NoError(t, db.Where("task_id = ?", tasks[1].Id).Find(&dependencies).Error)
	require.Len(t, dependencies, 1)
	require.Equal(t, tasks[1].Id, dependencies[0].TaskId)
	require.Equal(t, tasks[0].Id, dependencies[0].DependencyTaskId)
	var assignees []*model.ProjectManagementTaskAssignee
	require.NoError(t, db.Where("task_id = ? AND stage_id IS NULL", tasks[0].Id).Find(&assignees).Error)
	require.Len(t, assignees, 1)
	require.Equal(t, int64(8), assignees[0].OrganizationUserId)
	require.NoError(t, db.Model(&model.ProjectManagementWorkflow{}).Where("id = ?", 10).Update("name", "Changed source").Error)
	workflows, err = logic.Workflows(t.Context(), utils.JwtUser{OrganizationId: 9}, template.Id)
	require.NoError(t, err)
	require.Equal(t, "Source flow", workflows[0].Name)
	require.NoError(t, db.Exec("INSERT INTO project_management_projects (id, organization_id, name, priority, project_type) VALUES (2, 10, 'Foreign source', 'High', 'Ordinary')").Error)
	_, err = logic.Create(t.Context(), utils.JwtUser{OrganizationId: 9, OrganizationUserId: 7}, dto.ProjectManagementTemplateCreateDto{Name: "Foreign snapshot", SourceProjectId: 2})
	require.Error(t, err)
	withoutTasks, err := logic.Create(t.Context(), utils.JwtUser{OrganizationId: 9, OrganizationUserId: 7}, dto.ProjectManagementTemplateCreateDto{Name: "Structure only", SourceProjectId: 1})
	require.NoError(t, err)
	var withoutTaskCount int64
	require.NoError(t, db.Model(&model.ProjectManagementTask{}).Where("entity_type = ? AND entity_id = ?", constants.ProjectManagementOwnerTypeTemplate, withoutTasks.Id).Count(&withoutTaskCount).Error)
	require.Zero(t, withoutTaskCount)
	_, err = logic.Create(t.Context(), utils.JwtUser{OrganizationId: 9, OrganizationUserId: 7}, dto.ProjectManagementTemplateCreateDto{Name: "Invalid", SourceProjectId: 1, CopyAssignees: true})
	require.Error(t, err)
}
