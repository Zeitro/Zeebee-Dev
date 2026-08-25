package logic

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/repo"
)

func TestValidateProjectManagementTaskCreate(t *testing.T) {
	if err := validateProjectManagementTaskCreate(dto.ProjectManagementTaskCreateDto{}); err == nil {
		t.Fatal("invalid task creation accepted")
	}
	if err := validateProjectManagementTaskCreate(dto.ProjectManagementTaskCreateDto{ProjectManagementTaskDto: dto.ProjectManagementTaskDto{ProjectId: 1, Title: "Task"}}); err != nil {
		t.Fatalf("valid task creation rejected: %v", err)
	}
}

func TestValidateProjectManagementTaskTimeRange(t *testing.T) {
	startAt := &sql.NullTime{Time: time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC), Valid: true}
	endAt := &sql.NullTime{Time: time.Date(2026, time.January, 2, 0, 0, 0, 0, time.UTC), Valid: true}
	if err := validateProjectManagementTaskTimeRange(startAt, endAt); err != nil {
		t.Fatalf("valid time range rejected: %v", err)
	}
	if err := validateProjectManagementTaskTimeRange(endAt, startAt); err == nil {
		t.Fatal("reversed time range accepted")
	}
	if err := validateProjectManagementTaskTimeRange(startAt, startAt); err == nil {
		t.Fatal("equal time range accepted")
	}
}

func TestCreateTaskKeepsWorkflowAndTypeSnapshots(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	for _, statement := range []string{
		"CREATE TABLE organization_users (id integer primary key, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, organization_id integer not null, user_id integer not null default 0, full_name text not null default '', status text not null default '')",
		"CREATE TABLE users (id integer primary key, is_del integer not null default 0)",
		"CREATE TABLE organization_user_and_roles (organization_user_id integer not null, role_id integer not null)",
		"CREATE TABLE project_management_projects (id integer primary key, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, organization_id integer not null, entity_id integer not null default 0, entity_type text not null default '', name text not null default '', description text not null default '', cover_url text not null default '', owner_id integer not null default 0, owner_name text not null default '')",
		"CREATE TABLE project_management_project_members (id integer primary key, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, project_id integer not null, organization_user_id integer not null, role text not null default 'TaskExecutor')",
		"CREATE TABLE project_management_workflows (id integer primary key, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, kind text not null, organization_id integer not null, project_id integer not null, name text not null default '', description text not null default '', is_default boolean not null default false, source_template_id integer not null default 0, source_template_name text not null default '')",
		"CREATE TABLE project_management_workflow_stages (id integer primary key, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, workflow_id integer not null, name text not null default '', phase text not null, color text not null default '', sort_order integer not null default 0)",
		"CREATE TABLE project_management_workflow_transitions (id integer primary key, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, workflow_id integer not null, name text not null default '', from_stage_id integer not null, to_stage_id integer not null, allowed_actor_types json not null default '[]')",
		"CREATE TABLE project_management_task_types (id integer primary key, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, organization_id integer not null, name text not null default '', description text not null default '', color text not null default '', enabled boolean not null default true)",
		"CREATE TABLE project_management_tasks (id integer primary key autoincrement, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, project_id integer not null, organization_id integer not null, parent_task_id integer not null default 0, stage_id integer not null, creator_id integer not null, creator_name text not null default '', title text not null default '', description text not null default '', priority text not null default '', start_at datetime, end_at datetime, task_type json not null default '{}')",
		"CREATE TABLE project_management_task_assignees (id integer primary key autoincrement, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, task_id integer not null, stage_id integer not null default 0, organization_user_id integer not null, organization_user_name text not null default '')",
		"CREATE TABLE project_management_task_dependencies (id integer primary key autoincrement, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, task_id integer not null, dependency_task_id integer not null, dependency_type text not null default 'PreTask', dependency_rule text not null default 'FS')",
		"CREATE TABLE files (id integer primary key autoincrement, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, entity_type text not null default '', entity_id integer not null default 0, file_url text not null default '', file_name text not null default '', storage_type text not null default '', external_id text not null default '', upload_type text not null default '', drive_url text not null default '', drive_synced_at datetime)",
		"CREATE TABLE operations (id integer primary key autoincrement, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, entity_id integer not null, entity_type text not null, operator_id integer not null, operator_name text not null default '', operator_type text not null, operation text not null, note text not null default '', metadata json not null default '{}', organization_id integer not null)",
	} {
		if err := db.Exec(statement).Error; err != nil {
			t.Fatal(err)
		}
	}
	if err := db.Exec("ALTER TABLE project_management_tasks ADD COLUMN phase text not null default 'Pending'").Error; err != nil {
		t.Fatal(err)
	}
	for _, statement := range []string{
		"INSERT INTO organization_users (id, organization_id, full_name, status) VALUES (20, 1, 'Creator', 'Active'), (21, 1, 'Executor', 'Active'), (22, 1, 'Second executor', 'Active'), (23, 1, 'Not member', 'Active')",
		"INSERT INTO project_management_projects (id, organization_id, name, owner_id) VALUES (10, 1, 'Project', 20)",
		"INSERT INTO project_management_project_members (id, project_id, organization_user_id, role) VALUES (1, 10, 20, 'ProjectOwner'), (2, 10, 21, 'TaskExecutor'), (3, 10, 22, 'TaskExecutor')",
		"INSERT INTO project_management_workflows (id, kind, organization_id, project_id, name) VALUES (100, 'Project', 1, 10, 'Workflow')",
		"INSERT INTO project_management_workflow_stages (id, workflow_id, name, phase, sort_order) VALUES (200, 100, 'Pending', 'Pending', 1), (201, 100, 'Finished', 'Finished', 2)",
		"INSERT INTO project_management_task_types (id, organization_id, name, description, color, enabled) VALUES (500, 1, 'Inspection', 'Initial visit', '#123456', true)",
	} {
		if err := db.Exec(statement).Error; err != nil {
			t.Fatal(err)
		}
	}
	logic := NewProjectManagementTaskLogic(NewBaseLogic(repo.NewAllRepos(db), nil, nil, dto.JwtConfig{}))
	actor := utils.JwtUser{OrganizationId: 1, OrganizationUserId: 20, Username: "Creator"}
	task, err := logic.Create(context.Background(), actor, dto.ProjectManagementTaskCreateDto{
		ProjectManagementTaskDto:     dto.ProjectManagementTaskDto{ProjectId: 10, Title: "Inspect roof"},
		TaskTypeId:                   500,
		AssigneeOrganizationUserIds: []int64{21},
	})
	if err != nil || task.StageId != 200 || task.Phase != constants.ProjectManagementStagePhasePending || len(task.Assignees) != 1 || task.Assignees[0].OrganizationUserName != "Executor" || task.TaskType == nil || task.TaskType.Id != 500 || task.TaskType.Name != "Inspection" || task.TaskType.Description != "Initial visit" || task.TaskType.Color != "#123456" {
		t.Fatalf("task snapshots were not created: %#v, %v", task, err)
	}
	if err = db.Exec("INSERT INTO files (entity_type, entity_id, file_url, file_name, storage_type) VALUES (?, ?, 'https://example.test/task.pdf', 'task.pdf', 'Aws S3')", constants.FileEntityTypeProjectManagementTask, task.Id).Error; err != nil {
		t.Fatal(err)
	}
	task, err = logic.Detail(context.Background(), actor, task.Id)
	if err != nil || len(task.Attachments) != 1 || task.Attachments[0].FileName != "task.pdf" {
		t.Fatalf("task attachments were not returned: %#v, %v", task, err)
	}
	assigneeId := task.Assignees[0].Id
	task, err = logic.Update(context.Background(), actor, task.Id, dto.ProjectManagementTaskUpdateDto{Title: task.Title, TaskTypeId: 500, AssigneeOrganizationUserIds: []int64{21, 22}})
	if err != nil || task.TaskType == nil || task.TaskType.Id != 500 || len(task.Assignees) != 2 || task.Assignees[0].OrganizationUserId != 21 || task.Assignees[0].Id != assigneeId || task.Assignees[1].OrganizationUserId != 22 {
		t.Fatalf("task assignee was not updated: %#v, %v", task, err)
	}
	task, err = logic.Update(context.Background(), actor, task.Id, dto.ProjectManagementTaskUpdateDto{Title: task.Title, AssigneeOrganizationUserIds: []int64{22}})
	if err != nil || len(task.Assignees) != 1 || task.Assignees[0].OrganizationUserId != 22 {
		t.Fatalf("task assignee was not removed: %#v, %v", task, err)
	}
	if _, err = logic.Update(context.Background(), actor, task.Id, dto.ProjectManagementTaskUpdateDto{Title: task.Title, AssigneeOrganizationUserIds: []int64{22, 23}}); err == nil {
		t.Fatal("non-project member was assigned")
	}
	task, err = logic.Update(context.Background(), actor, task.Id, dto.ProjectManagementTaskUpdateDto{Title: task.Title})
	if err != nil || len(task.Assignees) != 0 {
		t.Fatalf("task assignee was not removed: %#v, %v", task, err)
	}
	unassignedTask, err := logic.Create(context.Background(), actor, dto.ProjectManagementTaskCreateDto{
		ProjectManagementTaskDto: dto.ProjectManagementTaskDto{ProjectId: 10, Title: "Unassigned task"},
		TaskTypeId:               500,
	})
	if err != nil || len(unassignedTask.Assignees) != 0 {
		t.Fatalf("unassigned task was not created: %#v, %v", unassignedTask, err)
	}
	if _, err = logic.Create(context.Background(), utils.JwtUser{OrganizationId: 1, OrganizationUserId: 23, Username: "Not member"}, dto.ProjectManagementTaskCreateDto{ProjectManagementTaskDto: dto.ProjectManagementTaskDto{ProjectId: 10, Title: "Denied"}}); err == nil {
		t.Fatal("non-member created a task")
	}
	if err = db.Exec("UPDATE project_management_tasks SET stage_id = 201, phase = 'Finished' WHERE id = ?", task.Id).Error; err != nil {
		t.Fatal(err)
	}
	if _, err = logic.Create(context.Background(), actor, dto.ProjectManagementTaskCreateDto{ProjectManagementTaskDto: dto.ProjectManagementTaskDto{ProjectId: 10, ParentTaskId: task.Id, Title: "Child"}}); err == nil {
		t.Fatal("child task was added to a finished parent")
	}
	if err = logic.Delete(context.Background(), actor, unassignedTask.Id); err != nil {
		t.Fatalf("task deletion without confirmation rejected: %v", err)
	}
}
