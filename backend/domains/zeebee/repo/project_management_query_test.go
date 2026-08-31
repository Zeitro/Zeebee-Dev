package repo

import (
	"context"
	"fmt"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
)

func TestProjectManagementQueryFilters(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	for _, statement := range []string{
		"CREATE TABLE project_management_workflow_stages (id integer primary key, is_del integer not null default 0, sort_order integer not null, name text not null default '', phase text not null, color text not null default '')",
		"CREATE TABLE project_management_projects (id integer primary key, is_del integer not null default 0, name text not null default '')",
		"CREATE TABLE project_management_tasks (id integer primary key, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, project_id integer not null, organization_id integer not null, parent_task_id integer not null default 0, stage_id integer not null, creator_id integer not null, creator_name text not null default '', title text not null, description text not null default '', priority text not null default '', start_at datetime, end_at datetime, task_type json not null default '{}')",
		"CREATE TABLE project_management_task_assignees (id integer primary key, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, task_id integer not null, stage_id integer not null default 0, organization_user_id integer not null, organization_user_full_name text not null default '')",
		"CREATE TABLE operations (id integer primary key, is_del integer not null default 0, created_at datetime not null, updated_at datetime not null, entity_id integer not null, entity_type text not null, operator_id integer not null, operator_name text not null default '', operator_type text not null default '', operation text not null, note text not null default '', metadata json not null default '{}', organization_id integer not null)",
	} {
		if err := db.Exec(statement).Error; err != nil {
			t.Fatal(err)
		}
	}
	for _, statement := range []string{
		"INSERT INTO project_management_workflow_stages (id, sort_order, name, phase) VALUES (100, 10, 'In progress', 'in_progress'), (101, 20, 'Finished', 'finished')",
		"INSERT INTO project_management_projects (id, name) VALUES (10, 'Project')",
		"INSERT INTO project_management_tasks (id, project_id, organization_id, stage_id, creator_id, creator_name, title, priority, end_at, task_type) VALUES (1, 10, 1, 100, 7, 'Creator', 'Roof inspection', 'High', datetime('now', '-1 hour'), '{\"id\":5}')",
		"INSERT INTO project_management_tasks (id, project_id, organization_id, stage_id, creator_id, creator_name, title, priority, task_type) VALUES (2, 10, 1, 100, 8, 'Other', 'Assigned work', 'Low', '{\"id\":6}')",
		"INSERT INTO project_management_tasks (id, project_id, organization_id, stage_id, creator_id, creator_name, title, task_type) VALUES (3, 10, 1, 100, 8, 'Other', 'Not mine', '{\"id\":5}')",
		"INSERT INTO project_management_tasks (id, project_id, organization_id, stage_id, creator_id, creator_name, title, end_at, task_type) VALUES (4, 10, 1, 101, 7, 'Creator', 'Finished task', datetime('now', '-1 hour'), '{\"id\":5}')",
		"INSERT INTO project_management_tasks (id, project_id, organization_id, stage_id, creator_id, creator_name, title, end_at, task_type) VALUES (5, 10, 1, 100, 7, 'Creator', 'Future roof work', datetime('now', '+1 hour'), '{\"id\":5}')",
		"INSERT INTO project_management_task_assignees (id, task_id, stage_id, organization_user_id, organization_user_full_name) VALUES (1, 1, 100, 7, 'Assignee'), (2, 2, 100, 7, 'Assignee')",
		fmt.Sprintf("INSERT INTO operations (id, created_at, updated_at, entity_id, entity_type, operator_id, operation, metadata, organization_id) VALUES (1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 10, 'ProjectManagementProject', 7, '%s', '{\"task_id\":1}', 1), (2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 10, 'ProjectManagementProject', 7, '%s', '{\"task_id\":2}', 1), (3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 11, 'ProjectManagementProject', 7, '%s', '{\"task_id\":1}', 1)", constants.OperationActionProjectManagementTaskUpdate, constants.OperationActionProjectManagementTaskUpdate, constants.OperationActionProjectManagementTaskUpdate),
	} {
		if err := db.Exec(statement).Error; err != nil {
			t.Fatal(err)
		}
	}
	if err := db.Exec("ALTER TABLE project_management_tasks ADD COLUMN phase text not null default 'Pending'").Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Exec("UPDATE project_management_tasks SET phase = 'InProgress' WHERE id = 1").Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Exec("UPDATE project_management_tasks SET phase = 'Finished' WHERE id = 4").Error; err != nil {
		t.Fatal(err)
	}
	repo := NewProjectManagementTaskRepo(db)
	overdue := true
	page, err := repo.PageMine(context.Background(), 1, 7, dto.ProjectManagementMyTaskPageRequestDto{
		Page: 1, Limit: 20,
		Filter: dto.ProjectManagementTaskFilterDto{ProjectId: 10, StageId: 100, Phase: []constants.ProjectManagementStagePhase{constants.ProjectManagementStagePhaseInProgress}, TaskTypeId: 5, Keyword: "roof", Overdue: &overdue},
	})
	if err != nil || page.Total != 1 || len(page.List) != 1 || page.List[0].Id != 1 {
		t.Fatalf("mine filters failed: %#v, %v", page, err)
	}
	page, err = repo.PageMine(context.Background(), 1, 7, dto.ProjectManagementMyTaskPageRequestDto{Page: 1, Limit: 20, Filter: dto.ProjectManagementTaskFilterDto{AssigneeId: 7, TaskTypeId: 6}})
	if err != nil || page.Total != 1 || len(page.List) != 1 || page.List[0].Id != 2 {
		t.Fatalf("assignee filter failed: %#v, %v", page, err)
	}
	page, err = repo.PageMine(context.Background(), 1, 7, dto.ProjectManagementMyTaskPageRequestDto{Page: 1, Limit: 20, Filter: dto.ProjectManagementTaskFilterDto{ProjectId: 10, Phase: []constants.ProjectManagementStagePhase{constants.ProjectManagementStagePhasePending, constants.ProjectManagementStagePhaseInProgress}}})
	if err != nil || page.Total != 2 {
		t.Fatalf("mine should only contain current-stage assignments: %#v, %v", page, err)
	}
	page, err = repo.PageMine(context.Background(), 1, 7, dto.ProjectManagementMyTaskPageRequestDto{Page: 1, Limit: 20, Filter: dto.ProjectManagementTaskFilterDto{Priority: constants.ProjectManagementPriorityHigh}})
	if err != nil || page.Total != 1 || len(page.List) != 1 || page.List[0].Id != 1 {
		t.Fatalf("mine priority filter failed: %#v, %v", page, err)
	}
	page, err = repo.PageByProject(context.Background(), 1, dto.ProjectManagementTaskPageRequestDto{Page: 1, Limit: 20, Filter: dto.ProjectManagementTaskFilterDto{ProjectId: 10, Phase: []constants.ProjectManagementStagePhase{constants.ProjectManagementStagePhaseInProgress}}})
	if err != nil || page.Total != 1 || len(page.List) != 1 || page.List[0].Id != 1 {
		t.Fatalf("multiple phase filter failed: %#v, %v", page, err)
	}
	page, err = repo.PageByProject(context.Background(), 1, dto.ProjectManagementTaskPageRequestDto{Page: 1, Limit: 20, Filter: dto.ProjectManagementTaskFilterDto{ProjectId: 10, Priority: constants.ProjectManagementPriorityLow}})
	if err != nil || page.Total != 1 || len(page.List) != 1 || page.List[0].Id != 2 {
		t.Fatalf("project priority filter failed: %#v, %v", page, err)
	}
	candidates, err := repo.CandidateByProject(context.Background(), 1, dto.ProjectManagementTaskCandidateDto{ProjectId: 10})
	if err != nil || len(candidates) != 4 {
		t.Fatalf("task candidates should exclude finished tasks: %#v, %v", candidates, err)
	}
	candidates, err = repo.CandidateByProject(context.Background(), 1, dto.ProjectManagementTaskCandidateDto{ProjectId: 10, TaskId: 1, Keyword: "roof"})
	if err != nil || len(candidates) != 1 || candidates[0].Id != 5 {
		t.Fatalf("task candidates should filter keyword and current task: %#v, %v", candidates, err)
	}
	operations, err := NewOperationRepo(db).PageProjectManagementOperations(context.Background(), 1, dto.ProjectManagementOperationPageRequestDto{
		Page: 1, Limit: 20,
		Filter: dto.ProjectManagementOperationFilterDto{ProjectId: 10, TaskId: 1, OperationType: constants.OperationActionProjectManagementTaskUpdate},
	})
	if err != nil || operations.Total != 1 || len(operations.List) != 1 || operations.List[0].Id != 1 {
		t.Fatalf("operation filters failed: %#v, %v", operations, err)
	}
}

func TestFindDetailByOrganizationLoadsCurrentStageAssignees(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	for _, statement := range []string{
		"CREATE TABLE project_management_tasks (id integer primary key, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, project_id integer not null, organization_id integer not null, parent_task_id integer not null default 0, stage_id integer not null, phase text not null default 'Pending', creator_id integer not null default 0, creator_name text not null default '', title text not null default '', description text not null default '', priority text not null default '', start_at datetime, end_at datetime, task_type json not null default '{}')",
		"CREATE TABLE project_management_task_assignees (id integer primary key, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, task_id integer not null, stage_id integer not null, organization_user_id integer not null, organization_user_full_name text not null default '')",
		"CREATE TABLE project_management_task_dependencies (id integer primary key, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, task_id integer not null, dependency_task_id integer not null, dependency_type text not null default 'PreTask', dependency_rule text not null default 'FS')",
		"CREATE TABLE files (id integer primary key, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, entity_id integer not null, entity_type text not null, file_url text not null default '', file_name text not null default '', storage_type text not null default '')",
		"INSERT INTO project_management_tasks (id, project_id, organization_id, stage_id, title) VALUES (1, 10, 1, 101, 'Task')",
		"INSERT INTO project_management_task_assignees (id, task_id, stage_id, organization_user_id, organization_user_full_name) VALUES (1, 1, 100, 7, 'Previous stage'), (2, 1, 101, 8, 'Current stage')",
	} {
		if err := db.Exec(statement).Error; err != nil {
			t.Fatal(err)
		}
	}
	task, err := NewProjectManagementTaskRepo(db).FindDetailByOrganization(context.Background(), 1, 1)
	if err != nil || len(task.Assignees) != 1 || task.Assignees[0].OrganizationUserId != 8 {
		t.Fatalf("detail assignees should only include the current stage: %#v, %v", task, err)
	}
}

func TestCopyCurrentStageAssigneesIfTargetEmpty(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	for _, statement := range []string{
		"CREATE TABLE project_management_task_assignees (id integer primary key autoincrement, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, task_id integer not null, stage_id integer not null, organization_user_id integer not null, organization_user_full_name text not null default '')",
		"INSERT INTO project_management_task_assignees (task_id, stage_id, organization_user_id, organization_user_full_name) VALUES (1, 100, 7, 'Current')",
	} {
		if err := db.Exec(statement).Error; err != nil {
			t.Fatal(err)
		}
	}
	repo := NewProjectManagementTaskRepo(db)
	task := &model.ProjectManagementTask{Base: model.Base{Id: 1}, StageId: 100}
	if err := db.Transaction(func(tx *gorm.DB) error {
		return repo.CopyCurrentStageAssigneesIfTargetEmptyWithTx(context.Background(), tx, task, 101)
	}); err != nil {
		t.Fatal(err)
	}
	var assignees []*model.ProjectManagementTaskAssignee
	if err := db.Where("task_id = ? AND stage_id = ?", 1, 101).Find(&assignees).Error; err != nil || len(assignees) != 1 || assignees[0].OrganizationUserId != 7 {
		t.Fatalf("current stage assignees were not copied: %#v, %v", assignees, err)
	}
	if err := db.Exec("INSERT INTO project_management_task_assignees (task_id, stage_id, organization_user_id, organization_user_full_name) VALUES (1, 100, 8, 'Another current')").Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Transaction(func(tx *gorm.DB) error {
		return repo.CopyCurrentStageAssigneesIfTargetEmptyWithTx(context.Background(), tx, task, 101)
	}); err != nil {
		t.Fatal(err)
	}
	assignees = nil
	if err := db.Where("task_id = ? AND stage_id = ?", 1, 101).Find(&assignees).Error; err != nil || len(assignees) != 1 || assignees[0].OrganizationUserId != 7 {
		t.Fatalf("existing target assignees should not be replaced: %#v, %v", assignees, err)
	}
}
