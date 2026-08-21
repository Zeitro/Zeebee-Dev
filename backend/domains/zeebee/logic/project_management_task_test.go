package logic

import (
	"context"
	"fmt"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"zeebee-backend/common/utils"
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

func TestValidateProjectManagementTaskDelete(t *testing.T) {
	if err := validateProjectManagementTaskDelete(false); err == nil {
		t.Fatal("unconfirmed task deletion accepted")
	}
	if err := validateProjectManagementTaskDelete(true); err != nil {
		t.Fatalf("confirmed task deletion rejected: %v", err)
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
		"CREATE TABLE project_management_workflows (id integer primary key, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, kind text not null, organization_id integer not null, project_id integer not null, name text not null default '', description text not null default '', enabled boolean not null default true, is_default boolean not null default false, source_template_id integer not null default 0, source_template_name text not null default '')",
		"CREATE TABLE project_management_workflow_stages (id integer primary key, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, workflow_id integer not null, name text not null default '', phase text not null, sort_order integer not null default 0)",
		"CREATE TABLE project_management_workflow_transitions (id integer primary key, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, workflow_id integer not null, from_stage_id integer not null, to_stage_id integer not null, allowed_actor_types json not null default '[]')",
		"CREATE TABLE project_management_task_types (id integer primary key, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, organization_id integer not null, name text not null default '', description text not null default '', color text not null default '', enabled boolean not null default true)",
		"CREATE TABLE project_management_tasks (id integer primary key autoincrement, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, project_id integer not null, organization_id integer not null, parent_task_id integer not null default 0, stage_id integer not null, creator_id integer not null, creator_name text not null default '', title text not null default '', description text not null default '', priority text not null default '', start_at datetime, end_at datetime, task_type json not null default '{}')",
		"CREATE TABLE project_management_task_assignees (id integer primary key autoincrement, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, task_id integer not null, organization_user_id integer not null, organization_user_name text not null default '')",
		"CREATE TABLE operations (id integer primary key autoincrement, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, entity_id integer not null, entity_type text not null, operator_id integer not null, operator_name text not null default '', operator_type text not null, operation text not null, note text not null default '', metadata json not null default '{}', organization_id integer not null)",
	} {
		if err := db.Exec(statement).Error; err != nil {
			t.Fatal(err)
		}
	}
	for _, statement := range []string{
		"INSERT INTO organization_users (id, organization_id, full_name, status) VALUES (20, 1, 'Creator', 'Active'), (21, 1, 'Executor', 'Active'), (22, 1, 'Not member', 'Active')",
		"INSERT INTO project_management_projects (id, organization_id, name, owner_id) VALUES (10, 1, 'Project', 20)",
		"INSERT INTO project_management_project_members (id, project_id, organization_user_id, role) VALUES (1, 10, 20, 'ProjectOwner'), (2, 10, 21, 'TaskExecutor')",
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
		ProjectManagementTaskDto: dto.ProjectManagementTaskDto{ProjectId: 10, Title: "Inspect roof"},
		TaskTypeId:               500,
		AssigneeIds:              []int64{21},
	})
	if err != nil || task.StageId != 200 || len(task.Assignees) != 1 || task.Assignees[0].OrganizationUserName != "Executor" || task.TaskType == nil || task.TaskType.Id != 500 || task.TaskType.Name != "Inspection" || task.TaskType.Description != "Initial visit" || task.TaskType.Color != "#123456" {
		t.Fatalf("task snapshots were not created: %#v, %v", task, err)
	}
	if _, err = logic.Create(context.Background(), utils.JwtUser{OrganizationId: 1, OrganizationUserId: 22, Username: "Not member"}, dto.ProjectManagementTaskCreateDto{ProjectManagementTaskDto: dto.ProjectManagementTaskDto{ProjectId: 10, Title: "Denied"}}); err == nil {
		t.Fatal("non-member created a task")
	}
	if err = db.Exec("UPDATE project_management_tasks SET stage_id = 201 WHERE id = ?", task.Id).Error; err != nil {
		t.Fatal(err)
	}
	if _, err = logic.Create(context.Background(), actor, dto.ProjectManagementTaskCreateDto{ProjectManagementTaskDto: dto.ProjectManagementTaskDto{ProjectId: 10, ParentTaskId: task.Id, Title: "Child"}}); err == nil {
		t.Fatal("child task was added to a finished parent")
	}
}
