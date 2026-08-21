package logic

import (
	"context"
	"fmt"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/repo"
)

func TestProjectManagementTaskTypeSnapshot(t *testing.T) {
	taskType := &model.ProjectManagementTaskType{Base: model.Base{Id: 9}, Name: "Inspection", Description: "Initial visit", Color: "#123456"}
	snapshot := ProjectManagementTaskTypeSnapshot(taskType)
	taskType.Name = "Changed"

	if snapshot.Id != 9 || snapshot.Name != "Inspection" || snapshot.Description != "Initial visit" || snapshot.Color != "#123456" {
		t.Fatalf("snapshot changed: %#v", snapshot)
	}
}

func TestValidateProjectManagementTaskType(t *testing.T) {
	if err := validateProjectManagementTaskType(dto.ProjectManagementTaskTypeDto{Name: "Inspection"}); err != nil {
		t.Fatalf("valid task type rejected: %v", err)
	}
	if err := validateProjectManagementTaskType(dto.ProjectManagementTaskTypeDto{Name: " \t"}); err == nil {
		t.Fatal("blank task type accepted")
	}
}

func TestProjectManagementTaskTypeAllCopiesGlobalTemplatesAndKeepsOne(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	for _, statement := range []string{
		"CREATE TABLE project_management_task_types (id integer primary key, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, organization_id integer not null, name text not null default '', description text not null default '', color text not null default '', sort_order integer not null default 0)",
		"CREATE TABLE operations (id integer primary key autoincrement, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, entity_id integer not null, entity_type text not null, operator_id integer not null, operator_name text not null default '', operator_type text not null, operation text not null, note text not null default '', metadata json not null default '{}', organization_id integer not null)",
		"INSERT INTO project_management_task_types (id, organization_id, name, description, color, sort_order) VALUES (1, 0, 'Inspection', 'Initial visit', '#111111', 20), (2, 0, 'Documentation', 'Gather documents', '#222222', 10)",
	} {
		if err := db.Exec(statement).Error; err != nil {
			t.Fatal(err)
		}
	}
	logic := NewProjectManagementTaskTypeLogic(NewBaseLogic(repo.NewAllRepos(db), nil, nil, dto.JwtConfig{}))
	actor := utils.JwtUser{OrganizationId: 1, OrganizationUserId: 7, Username: "Admin"}

	taskTypes, err := logic.All(context.Background(), actor.OrganizationId)
	if err != nil {
		t.Fatalf("all task types: %v", err)
	}
	if len(taskTypes) != 2 || taskTypes[0].Name != "Documentation" || taskTypes[0].SortOrder != 10 || taskTypes[1].OrganizationId != actor.OrganizationId {
		t.Fatalf("global templates were not copied and sorted: %#v", taskTypes)
	}
	detail, err := logic.Detail(context.Background(), actor.OrganizationId, taskTypes[0].Id)
	if err != nil || detail.Id != taskTypes[0].Id {
		t.Fatalf("detail task type: %#v, %v", detail, err)
	}
	if _, err = logic.Detail(context.Background(), 2, taskTypes[0].Id); err == nil {
		t.Fatal("task type was visible to another organization")
	}
	updated, err := logic.Update(context.Background(), actor, taskTypes[1].Id, dto.ProjectManagementTaskTypeDto{Name: "Updated inspection", Description: "Initial visit", Color: "#111111", SortOrder: 7})
	if err != nil {
		t.Fatalf("update task type: %v", err)
	}
	if updated.SortOrder != 7 {
		t.Fatalf("task type update did not change sort order: %#v", updated)
	}
	if err = logic.Order(context.Background(), actor, updated.Id, 5); err != nil {
		t.Fatalf("order task type: %v", err)
	}
	taskTypes, err = logic.All(context.Background(), actor.OrganizationId)
	if err != nil {
		t.Fatalf("all copied task types: %v", err)
	}
	if taskTypes[0].Id != updated.Id || taskTypes[0].SortOrder != 5 {
		t.Fatalf("task types were not reordered: %#v", taskTypes)
	}
	if err = logic.Delete(context.Background(), actor, taskTypes[0].Id); err != nil {
		t.Fatalf("delete task type: %v", err)
	}
	if err = logic.Delete(context.Background(), actor, taskTypes[1].Id); err == nil {
		t.Fatal("last task type was deleted")
	}
}
