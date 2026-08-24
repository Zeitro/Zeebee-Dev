package repo

import (
	"context"
	"fmt"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
)

func TestProjectManagementQueryFilters(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	for _, statement := range []string{
		"CREATE TABLE project_management_workflow_stages (id integer primary key, is_del integer not null default 0, sort_order integer not null, phase text not null, color text not null default '')",
		"CREATE TABLE project_management_tasks (id integer primary key, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, project_id integer not null, organization_id integer not null, parent_task_id integer not null default 0, stage_id integer not null, creator_id integer not null, creator_name text not null default '', title text not null, description text not null default '', priority text not null default '', start_at datetime, end_at datetime, task_type json not null default '{}')",
		"CREATE TABLE project_management_task_assignees (id integer primary key, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, task_id integer not null, stage_id integer not null default 0, organization_user_id integer not null, organization_user_name text not null default '')",
		"CREATE TABLE operations (id integer primary key, is_del integer not null default 0, created_at datetime not null, updated_at datetime not null, entity_id integer not null, entity_type text not null, operator_id integer not null, operator_name text not null default '', operator_type text not null default '', operation text not null, note text not null default '', metadata json not null default '{}', organization_id integer not null)",
	} {
		if err := db.Exec(statement).Error; err != nil {
			t.Fatal(err)
		}
	}
	for _, statement := range []string{
		"INSERT INTO project_management_workflow_stages (id, sort_order, phase) VALUES (100, 10, 'in_progress'), (101, 20, 'finished')",
		"INSERT INTO project_management_tasks (id, project_id, organization_id, stage_id, creator_id, creator_name, title, end_at, task_type) VALUES (1, 10, 1, 100, 7, 'Creator', 'Roof inspection', datetime('now', '-1 hour'), '{\"id\":5}')",
		"INSERT INTO project_management_tasks (id, project_id, organization_id, stage_id, creator_id, creator_name, title, task_type) VALUES (2, 10, 1, 100, 8, 'Other', 'Assigned work', '{\"id\":6}')",
		"INSERT INTO project_management_tasks (id, project_id, organization_id, stage_id, creator_id, creator_name, title, task_type) VALUES (3, 10, 1, 100, 8, 'Other', 'Not mine', '{\"id\":5}')",
		"INSERT INTO project_management_tasks (id, project_id, organization_id, stage_id, creator_id, creator_name, title, end_at, task_type) VALUES (4, 11, 1, 101, 7, 'Creator', 'Finished task', datetime('now', '-1 hour'), '{\"id\":5}')",
		"INSERT INTO project_management_tasks (id, project_id, organization_id, stage_id, creator_id, creator_name, title, end_at, task_type) VALUES (5, 10, 1, 100, 7, 'Creator', 'Future roof work', datetime('now', '+1 hour'), '{\"id\":5}')",
		"INSERT INTO project_management_task_assignees (id, task_id, stage_id, organization_user_id, organization_user_name) VALUES (1, 2, 100, 7, 'Assignee')",
		"INSERT INTO operations (id, created_at, updated_at, entity_id, entity_type, operator_id, operation, metadata, organization_id) VALUES (1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 10, 'ProjectManagementProject', 7, 'Project task updated', '{\"task_id\":1}', 1), (2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 10, 'ProjectManagementProject', 7, 'Project task updated', '{\"task_id\":2}', 1), (3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 11, 'ProjectManagementProject', 7, 'Project task updated', '{\"task_id\":1}', 1)",
	} {
		if err := db.Exec(statement).Error; err != nil {
			t.Fatal(err)
		}
	}
	repo := NewProjectManagementTaskRepo(db)
	overdue := true
	page, err := repo.PageMine(context.Background(), 1, 7, dto.ProjectManagementMyTaskPageRequestDto{
		Page: 1, Limit: 20,
		Filter: dto.ProjectManagementTaskFilterDto{ProjectId: 10, StageId: 100, TaskTypeId: 5, Keyword: "roof", Overdue: &overdue},
	})
	if err != nil || page.Total != 1 || len(page.List) != 1 || page.List[0].Id != 1 {
		t.Fatalf("mine filters failed: %#v, %v", page, err)
	}
	page, err = repo.PageMine(context.Background(), 1, 7, dto.ProjectManagementMyTaskPageRequestDto{Page: 1, Limit: 20, Filter: dto.ProjectManagementTaskFilterDto{AssigneeId: 7}})
	if err != nil || page.Total != 1 || len(page.List) != 1 || page.List[0].Id != 2 {
		t.Fatalf("assignee filter failed: %#v, %v", page, err)
	}
	operations, err := NewOperationRepo(db).PageProjectManagementOperations(context.Background(), 1, dto.ProjectManagementOperationPageRequestDto{
		Page: 1, Limit: 20,
		Filter: dto.ProjectManagementOperationFilterDto{ProjectId: 10, TaskId: 1, OperationType: string(constants.OperationActionProjectManagementTaskUpdate)},
	})
	if err != nil || operations.Total != 1 || len(operations.List) != 1 || operations.List[0].Id != 1 {
		t.Fatalf("operation filters failed: %#v, %v", operations, err)
	}
}
