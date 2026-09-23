package logic

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/repo"
)

func TestCalculateProjectEndAt(t *testing.T) {
	tests := []struct {
		name    string
		start   string
		value   int
		unit    constants.ProjectManagementPeriodUnit
		want    string
		wantErr bool
	}{
		{name: "day", start: "2026-09-01", value: 3, unit: constants.ProjectManagementPeriodUnitDay, want: "2026-09-03"},
		{name: "week", start: "2026-09-01", value: 5, unit: constants.ProjectManagementPeriodUnitWeek, want: "2026-10-05"},
		{name: "month end", start: "2026-01-31", value: 2, unit: constants.ProjectManagementPeriodUnitMonth, want: "2026-02-28"},
		{name: "leap year", start: "2024-01-31", value: 2, unit: constants.ProjectManagementPeriodUnitMonth, want: "2024-02-29"},
		{name: "invalid value", start: "2026-09-01", value: 0, unit: constants.ProjectManagementPeriodUnitDay, wantErr: true},
		{name: "invalid unit", start: "2026-09-01", value: 1, unit: constants.ProjectManagementPeriodUnit("Year"), wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start, err := time.Parse(time.DateOnly, tt.start)
			if err != nil {
				t.Fatal(err)
			}
			got, err := CalculateProjectEndAt(start, tt.value, tt.unit)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if got.Format(time.DateOnly) != tt.want {
				t.Fatalf("end date = %s, want %s", got.Format(time.DateOnly), tt.want)
			}
		})
	}
}

func TestPrepareProjectSchedule(t *testing.T) {
	start := sql.NullTime{Time: time.Date(2026, time.January, 31, 0, 0, 0, 0, time.UTC), Valid: true}
	planned := &model.ProjectManagementProject{ProjectType: constants.ProjectManagementProjectTypePlanned, StartAt: &start, PeriodValue: 2, PeriodUnit: constants.ProjectManagementPeriodUnitMonth}
	if err := prepareProjectSchedule(planned); err != nil {
		t.Fatal(err)
	}
	if planned.EndAt == nil || planned.EndAt.Time.Format(time.DateOnly) != "2026-02-28" {
		t.Fatalf("unexpected end date: %#v", planned.EndAt)
	}
	if err := prepareProjectSchedule(&model.ProjectManagementProject{ProjectType: constants.ProjectManagementProjectTypeOrdinary, PeriodValue: 1}); err == nil {
		t.Fatal("ordinary project schedule was accepted")
	}
}

func TestShiftTemplateTaskTimes(t *testing.T) {
	templateStart := sql.NullTime{Time: time.Date(2026, time.January, 10, 0, 0, 0, 0, time.UTC), Valid: true}
	projectStart := sql.NullTime{Time: time.Date(2026, time.February, 1, 0, 0, 0, 0, time.UTC), Valid: true}
	taskStart := sql.NullTime{Time: time.Date(2026, time.January, 12, 9, 0, 0, 0, time.UTC), Valid: true}
	start, _ := shiftTemplateTaskTimes(&model.ProjectManagementTemplate{StartAt: &templateStart}, &model.ProjectManagementProject{StartAt: &projectStart}, &taskStart, nil)
	if start.Time.Format(time.RFC3339) != "2026-02-03T09:00:00Z" {
		t.Fatalf("shifted start = %s", start.Time.Format(time.RFC3339))
	}
	location, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatal(err)
	}
	templateStart.Time = time.Date(2026, time.March, 8, 0, 0, 0, 0, location)
	projectStart.Time = time.Date(2026, time.March, 9, 0, 0, 0, 0, location)
	taskStart.Time = time.Date(2026, time.March, 8, 9, 0, 0, 0, location)
	start, _ = shiftTemplateTaskTimes(&model.ProjectManagementTemplate{StartAt: &templateStart}, &model.ProjectManagementProject{StartAt: &projectStart}, &taskStart, nil)
	if start.Time.Format(time.RFC3339) != "2026-03-09T09:00:00-04:00" {
		t.Fatalf("DST-shifted start = %s", start.Time.Format(time.RFC3339))
	}
}

func TestCreateProjectClonesTemplateAndRollsBack(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	for _, statement := range []string{
		"CREATE TABLE borrower_properties (id integer primary key, is_del integer not null default 0, organization_id integer not null)",
		"CREATE TABLE project_management_projects (id integer primary key autoincrement, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, organization_id integer not null, entity_id integer not null default 0, entity_type text not null default '', name text not null default '', description text not null default '', priority text not null default '', cover_url text not null default '', owner_id integer not null default 0, owner_name text not null default '', project_type text not null default 'Ordinary', start_at date, end_at date, period_value integer not null default 0, period_unit text not null default '')",
		"CREATE TABLE project_management_project_members (id integer primary key autoincrement, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, project_id integer not null, organization_user_id integer not null, role text not null)",
		"CREATE TABLE project_management_templates (id integer primary key, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, template_type text not null, organization_id integer not null, name text not null default '', description text not null default '', status text not null default 'Active', project_type text)",
		"CREATE TABLE project_management_workflows (id integer primary key autoincrement, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, entity_type text not null default 'Project', entity_id integer not null default 0, source_workflow_id integer not null default 0, source_workflow_name text not null default '', name text not null default '', description text not null default '', is_default boolean not null default false)",
		"CREATE TABLE project_management_stages (id integer primary key autoincrement, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, entity_type text not null, entity_id integer not null, name text not null default '', description text not null default '', phase text not null, color text not null default '', sort_order integer not null)",
		"CREATE TABLE project_management_workflow_stages (id integer primary key autoincrement, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, workflow_id integer not null, stage_id integer not null)",
		"CREATE TABLE project_management_workflow_transitions (id integer primary key autoincrement, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, workflow_id integer not null, name text not null default '', from_stage_id integer not null, to_stage_id integer not null, allowed_actor_types json not null default '[]')",
		"CREATE TABLE project_management_tasks (id integer primary key autoincrement, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, entity_type text not null, entity_id integer not null, workflow_id integer not null default 0, organization_id integer not null default 0, parent_task_id integer not null default 0, stage_id integer not null default 0, phase text not null default 'Pending', creator_id integer not null default 0, creator_name text not null default '', title text not null default '', description text not null default '', priority text not null default '', work_hours numeric not null default 0, start_at datetime, end_at datetime, task_type json not null default '{}')",
		"CREATE TABLE project_management_task_dependencies (id integer primary key autoincrement, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, task_id integer not null, dependency_task_id integer not null, dependency_type text not null, dependency_rule text not null)",
		"CREATE TABLE project_management_task_assignees (id integer primary key autoincrement, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, task_id integer not null, stage_id integer not null default 0, organization_user_id integer not null, organization_user_full_name text not null default '')",
		"CREATE TABLE organization_users (id integer primary key, created_at datetime, updated_at datetime, is_del integer not null default 0, organization_id integer not null, user_id integer not null default 0, status text not null, first_name text not null default '', mid_name text not null default '', last_name text not null default '', full_name text not null default '', phone text not null default '', remark text not null default '')",
		"CREATE TABLE operations (id integer primary key autoincrement, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, entity_id integer not null, entity_type text not null, operator_id integer not null, operator_name text not null default '', operator_type text not null, operation text not null, note text not null default '', metadata json not null default '{}', organization_id integer not null)",
	} {
		if err := db.Exec(statement).Error; err != nil {
			t.Fatal(err)
		}
	}
	for _, statement := range []string{
		"INSERT INTO borrower_properties (id, organization_id) VALUES (50, 1), (51, 2)",
		"INSERT INTO project_management_templates (id, template_type, organization_id, name, status, project_type) VALUES (90, 'Organization', 1, 'Inspection template', 'Active', 'Ordinary')",
		"INSERT INTO project_management_workflows (id, entity_type, entity_id, name, is_default) VALUES (100, 'Template', 90, 'Inspection', true)",
		"INSERT INTO project_management_stages (id, entity_type, entity_id, name, phase, sort_order) VALUES (200, 'Template', 90, 'Pending', 'Pending', 1), (201, 'Template', 90, 'Finished', 'Finished', 2)",
		"INSERT INTO project_management_workflow_stages (id, workflow_id, stage_id) VALUES (300, 100, 200), (301, 100, 201)",
		"INSERT INTO project_management_workflow_transitions (workflow_id, name, from_stage_id, to_stage_id) VALUES (100, 'Finish', 200, 201)",
		"INSERT INTO project_management_tasks (id, entity_type, entity_id, workflow_id, organization_id, title, priority) VALUES (400, 'Template', 90, 100, 1, 'Template task', 'High')",
		"INSERT INTO project_management_task_assignees (task_id, organization_user_id, organization_user_full_name) VALUES (400, 8, 'Executor'), (400, 9, 'Disabled')",
		"INSERT INTO organization_users (id, organization_id, status, full_name) VALUES (8, 1, 'Active', 'Executor'), (9, 1, 'Disabled', 'Disabled')",
	} {
		if err := db.Exec(statement).Error; err != nil {
			t.Fatal(err)
		}
	}
	logic := NewProjectManagementProjectLogic(NewBaseLogic(repo.NewAllRepos(db), nil, nil, dto.JwtConfig{}))
	actor := utils.JwtUser{OrganizationId: 1, OrganizationUserId: 7, Username: "Owner"}
	project, assigneeResult, err := logic.CreateWithAssigneeResult(context.Background(), actor, dto.ProjectManagementProjectCreateDto{Name: "Renovation", Description: "Initial description", Priority: constants.ProjectManagementPriorityHigh, CoverUrl: "https://example.com/cover.png", EntityId: 50, EntityType: constants.ProjectManagementEntityTypeBorrowerProperty, TemplateId: 90, CopyAssignees: true})
	if err != nil {
		t.Fatalf("project creation rejected: %v", err)
	}
	if project.EntityId != 50 || project.EntityType != constants.ProjectManagementEntityTypeBorrowerProperty {
		t.Fatalf("unexpected entity binding: %#v", project)
	}
	var member model.ProjectManagementProjectMember
	if err := db.Where("project_id = ? AND organization_user_id = ?", project.Id, actor.OrganizationUserId).First(&member).Error; err != nil {
		t.Fatal(err)
	}
	if member.Role != constants.ProjectManagementMemberRoleProjectOwner {
		t.Fatalf("expected project owner role, got %q", member.Role)
	}
	var workflow model.ProjectManagementWorkflow
	if err := db.Where("entity_type = ? AND entity_id = ?", constants.ProjectManagementOwnerTypeProject, project.Id).First(&workflow).Error; err != nil {
		t.Fatal(err)
	}
	if workflow.SourceWorkflowId != 100 || workflow.SourceWorkflowName != "Inspection" || !workflow.IsDefault {
		t.Fatalf("unexpected workflow snapshot: %#v", workflow)
	}
	var copiedTask model.ProjectManagementTask
	if err := db.Where("entity_type = ? AND entity_id = ? AND title = ?", constants.ProjectManagementOwnerTypeProject, project.Id, "Template task").First(&copiedTask).Error; err != nil {
		t.Fatal(err)
	}
	if copiedTask.WorkflowId != workflow.Id || copiedTask.Phase == nil || *copiedTask.Phase != constants.ProjectManagementStagePhasePending {
		t.Fatalf("template task was not instantiated as pending: %#v", copiedTask)
	}
	require.Equal(t, []int64{8}, assigneeResult.CopiedOrganizationUserIds)
	require.Len(t, assigneeResult.SkippedAssignees, 1)
	require.Equal(t, int64(9), assigneeResult.SkippedAssignees[0].OrganizationUserId)
	require.Equal(t, "organization_user_not_active", assigneeResult.SkippedAssignees[0].Reason)
	var copiedAssignees []*model.ProjectManagementTaskAssignee
	require.NoError(t, db.Where("task_id = ? AND stage_id = ?", copiedTask.Id, copiedTask.StageId).Find(&copiedAssignees).Error)
	require.Len(t, copiedAssignees, 1)
	require.Equal(t, int64(8), copiedAssignees[0].OrganizationUserId)
	var executor model.ProjectManagementProjectMember
	require.NoError(t, db.Where("project_id = ? AND organization_user_id = ?", project.Id, 8).First(&executor).Error)
	require.Equal(t, constants.ProjectManagementMemberRoleTaskExecutor, executor.Role)
	withoutAssignees, withoutAssigneeResult, err := logic.CreateWithAssigneeResult(context.Background(), actor, dto.ProjectManagementProjectCreateDto{Name: "No assignee clone", Priority: constants.ProjectManagementPriorityHigh, TemplateId: 90})
	require.NoError(t, err)
	require.Empty(t, withoutAssigneeResult.CopiedOrganizationUserIds)
	require.Empty(t, withoutAssigneeResult.SkippedAssignees)
	var withoutAssigneeMembers int64
	require.NoError(t, db.Model(&model.ProjectManagementProjectMember{}).Where("project_id = ? AND organization_user_id IN ?", withoutAssignees.Id, []int64{8, 9}).Count(&withoutAssigneeMembers).Error)
	require.Zero(t, withoutAssigneeMembers)
	var stageCount int64
	if err := db.Model(&model.ProjectManagementStage{}).Where("entity_type = ? AND entity_id = ?", constants.ProjectManagementOwnerTypeProject, project.Id).Count(&stageCount).Error; err != nil || stageCount != 2 {
		t.Fatalf("workflow stages were not cloned: %d, %v", stageCount, err)
	}
	var operation model.Operation
	if err := db.Where("entity_id = ? AND operation = ?", project.Id, string(constants.OperationActionProjectManagementProjectCreate)).Order("id DESC").First(&operation).Error; err != nil {
		t.Fatal(err)
	}
	var metadata map[string]any
	if err := json.Unmarshal([]byte(operation.Metadata), &metadata); err != nil {
		t.Fatal(err)
	}
	if metadata["template_id"] != float64(90) || metadata["template_type"] != "Organization" {
		t.Fatalf("template provenance missing from operation metadata: %#v", metadata)
	}
	transactionProject, err := logic.Create(context.Background(), actor, dto.ProjectManagementProjectCreateDto{Name: "Transaction project", EntityId: 51, EntityType: "Transaction", Priority: constants.ProjectManagementPriorityLow, TemplateId: 90})
	if err != nil || transactionProject.EntityId != 51 || transactionProject.EntityType != "Transaction" {
		t.Fatalf("transaction binding was rejected: %#v, %v", transactionProject, err)
	}
	if _, err := logic.Create(context.Background(), actor, dto.ProjectManagementProjectCreateDto{Name: "Missing priority", TemplateId: 90}); err == nil {
		t.Fatal("project without priority was created")
	}
	if _, err := logic.Create(context.Background(), actor, dto.ProjectManagementProjectCreateDto{Name: "Rejected", Priority: constants.ProjectManagementPriorityLow, TemplateId: 999}); err == nil {
		t.Fatal("project with missing template was created")
	}
	if _, err := logic.Create(context.Background(), actor, dto.ProjectManagementProjectCreateDto{Name: "Foreign property", EntityId: 51, EntityType: constants.ProjectManagementEntityTypeBorrowerProperty, Priority: constants.ProjectManagementPriorityLow, TemplateId: 90}); err == nil {
		t.Fatal("project with a property from another organization was created")
	}
	updated, err := logic.Update(context.Background(), actor, project.Id, dto.ProjectManagementProjectUpdateDto{Name: "Renovation updated", Priority: constants.ProjectManagementPriorityCritical})
	if err != nil {
		t.Fatalf("project update rejected: %v", err)
	}
	if updated.Name != "Renovation updated" || updated.Description != "" || updated.Priority != constants.ProjectManagementPriorityCritical || updated.CoverUrl != "" {
		t.Fatalf("unexpected updated project: %#v", updated)
	}
	var persisted model.ProjectManagementProject
	if err := db.First(&persisted, project.Id).Error; err != nil {
		t.Fatal(err)
	}
	if persisted.Name != "Renovation updated" || persisted.Description != "" || persisted.Priority != constants.ProjectManagementPriorityCritical || persisted.CoverUrl != "" {
		t.Fatalf("zero-value fields were not persisted: %#v", persisted)
	}
	var count int64
	if err := db.Model(&model.ProjectManagementProject{}).Count(&count).Error; err != nil {
		t.Fatal(err)
	}
	if count != 3 {
		t.Fatalf("failed project creation left %d projects", count)
	}
}
