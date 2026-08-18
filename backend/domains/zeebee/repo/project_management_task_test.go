package repo

import (
	"context"
	"fmt"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"zeebee-backend/domains/zeebee/constants"
)

func TestTaskCascadeDeleteAndHistoricalSnapshots(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	for _, statement := range []string{
		"CREATE TABLE project_management_tasks (id integer primary key, is_del integer not null default 0, project_id integer not null, organization_id integer not null, parent_task_id integer not null default 0, stage_id integer not null default 0, creator_id integer not null default 0, creator_name text not null default '', title text not null default '')",
		"CREATE TABLE project_management_task_assignees (id integer primary key, is_del integer not null default 0, task_id integer not null, organization_user_id integer not null, organization_user_name text not null default '')",
		"CREATE TABLE project_management_task_dependencies (id integer primary key, is_del integer not null default 0, successor_task_id integer not null, predecessor_task_id integer not null)",
		"CREATE TABLE files (id integer primary key, is_del integer not null default 0, entity_id integer not null, entity_type text not null)",
		"CREATE TABLE project_management_project_members (id integer primary key, is_del integer not null default 0, project_id integer not null, organization_user_id integer not null, role text not null default '')",
		"CREATE TABLE operations (id integer primary key, is_del integer not null default 0, entity_id integer not null, entity_type text not null, operation text not null default '')",
	} {
		if err := db.Exec(statement).Error; err != nil {
			t.Fatal(err)
		}
	}
	for _, statement := range []string{
		"INSERT INTO project_management_tasks (id, project_id, organization_id, stage_id, creator_id, creator_name, title) VALUES (1, 10, 1, 100, 8, 'Former creator', 'Historical task')",
		"INSERT INTO project_management_task_assignees (id, task_id, organization_user_id, organization_user_name) VALUES (2, 1, 9, 'Former executor')",
		"INSERT INTO project_management_project_members (id, project_id, organization_user_id, role) VALUES (3, 10, 9, 'executor')",
	} {
		if err := db.Exec(statement).Error; err != nil {
			t.Fatal(err)
		}
	}
	taskRepo, fileRepo := NewProjectManagementTaskRepo(db), NewFileRepo(db)
	if err := db.Exec("UPDATE project_management_project_members SET is_del = 1 WHERE organization_user_id = 9").Error; err != nil {
		t.Fatal(err)
	}
	historical, err := taskRepo.FindDetailByOrganization(context.Background(), 1, 1)
	if err != nil || historical.CreatorId != 8 || historical.CreatorName != "Former creator" || len(historical.Assignees) != 1 || historical.Assignees[0].OrganizationUserId != 9 || historical.Assignees[0].OrganizationUserName != "Former executor" {
		t.Fatalf("historical task identity was lost: %#v, %v", historical, err)
	}
	for _, statement := range []string{
		"INSERT INTO project_management_tasks (id, project_id, organization_id, parent_task_id, stage_id, creator_id, creator_name, title) VALUES (11, 10, 1, 0, 100, 8, 'Creator', 'Parent')",
		"INSERT INTO project_management_tasks (id, project_id, organization_id, parent_task_id, stage_id, creator_id, creator_name, title) VALUES (12, 10, 1, 11, 100, 8, 'Creator', 'Child')",
		"INSERT INTO project_management_task_assignees (id, task_id, organization_user_id, organization_user_name) VALUES (13, 12, 9, 'Executor')",
		"INSERT INTO project_management_task_dependencies (id, successor_task_id, predecessor_task_id) VALUES (14, 12, 11)",
		"INSERT INTO files (id, entity_id, entity_type) VALUES (15, 12, 'ProjectManagementTask')",
		"INSERT INTO operations (id, entity_id, entity_type, operation) VALUES (16, 10, 'ProjectManagementProject', 'Project task deleted')",
	} {
		if err := db.Exec(statement).Error; err != nil {
			t.Fatal(err)
		}
	}
	if err := db.Transaction(func(tx *gorm.DB) error {
		ids, err := taskRepo.DescendantIdsWithTx(context.Background(), tx, 11)
		if err != nil {
			return err
		}
		if err := fileRepo.DeleteByEntityIdWithTx(context.Background(), tx, ids, constants.FileEntityTypeProjectManagementTask); err != nil {
			return err
		}
		return taskRepo.DeleteByIdsWithTx(context.Background(), tx, ids)
	}); err != nil {
		t.Fatal(err)
	}
	for _, target := range []struct {
		table string
		id    int64
	}{
		{"project_management_tasks", 11}, {"project_management_tasks", 12}, {"project_management_task_assignees", 13}, {"project_management_task_dependencies", 14}, {"files", 15},
	} {
		var isDel int
		if err := db.Raw(fmt.Sprintf("SELECT is_del FROM %s WHERE id = ?", target.table), target.id).Scan(&isDel).Error; err != nil || isDel != 1 {
			t.Fatalf("cascade target not soft deleted: %#v, is_del=%d, err=%v", target, isDel, err)
		}
	}
	var operationIsDel int
	if err := db.Raw("SELECT is_del FROM operations WHERE id = ?", 16).Scan(&operationIsDel).Error; err != nil || operationIsDel != 0 {
		t.Fatalf("task deletion changed audit history: is_del=%d, err=%v", operationIsDel, err)
	}
}
