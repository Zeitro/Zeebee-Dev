package logic

import (
	"context"
	"fmt"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/repo"
)

func TestValidateProjectManagementDependency(t *testing.T) {
	valid := dto.ProjectManagementTaskDependencyDto{TaskId: 2, DependencyTaskId: 1, DependencyType: constants.ProjectManagementDependencyTypePreTask, DependencyRule: constants.ProjectManagementDependencyRuleFS}
	if err := validateProjectManagementDependency(valid); err != nil {
		t.Fatalf("valid dependency rejected: %v", err)
	}
	for _, invalid := range []dto.ProjectManagementTaskDependencyDto{
		{},
		{TaskId: 1, DependencyTaskId: 1, DependencyType: constants.ProjectManagementDependencyTypePreTask, DependencyRule: constants.ProjectManagementDependencyRuleFS},
		{TaskId: 2, DependencyTaskId: 1, DependencyType: "invalid", DependencyRule: constants.ProjectManagementDependencyRuleFS},
	} {
		if err := validateProjectManagementDependency(invalid); err == nil {
			t.Fatalf("invalid dependency accepted: %#v", invalid)
		}
	}
}

func TestCreateTaskDependencyRules(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	for _, statement := range []string{
		"CREATE TABLE project_management_tasks (id integer primary key, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, project_id integer not null, organization_id integer not null, parent_task_id integer not null default 0, stage_id integer not null, creator_id integer not null default 0, creator_name text not null default '', title text not null default '', description text not null default '', priority text not null default '', task_type json not null default '{}')",
		"CREATE TABLE project_management_project_members (id integer primary key, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, project_id integer not null, organization_user_id integer not null, role text not null default 'TaskExecutor')",
		"CREATE TABLE project_management_workflow_stages (id integer primary key, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, workflow_id integer not null default 0, name text not null default '', phase text not null, color text not null default '', sort_order integer not null default 0)",
		"CREATE TABLE project_management_task_dependencies (id integer primary key autoincrement, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, task_id integer not null, dependency_task_id integer not null, dependency_type text not null, dependency_rule text not null)",
		"CREATE TABLE operations (id integer primary key autoincrement, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, entity_id integer not null, entity_type text not null, operator_id integer not null, operator_name text not null default '', operator_type text not null, operation text not null, note text not null default '', metadata json not null default '{}', organization_id integer not null)",
	} {
		if err := db.Exec(statement).Error; err != nil {
			t.Fatal(err)
		}
	}
	for _, statement := range []string{
		"INSERT INTO project_management_workflow_stages (id, phase) VALUES (100, 'Pending'), (101, 'InProgress')",
		"INSERT INTO project_management_tasks (id, project_id, organization_id, stage_id, title) VALUES (1, 10, 1, 100, 'Predecessor'), (2, 10, 1, 100, 'Successor'), (3, 11, 1, 100, 'Other project'), (4, 10, 1, 100, 'Post task')",
		"INSERT INTO project_management_project_members (id, project_id, organization_user_id) VALUES (1, 10, 7)",
	} {
		if err := db.Exec(statement).Error; err != nil {
			t.Fatal(err)
		}
	}
	logic := NewProjectManagementDependencyLogic(NewBaseLogic(repo.NewAllRepos(db), nil, nil, dto.JwtConfig{}))
	actor := utils.JwtUser{OrganizationId: 1, OrganizationUserId: 7, Username: "Actor"}
	input := dto.ProjectManagementTaskDependencyDto{TaskId: 2, DependencyTaskId: 1, DependencyType: constants.ProjectManagementDependencyTypePreTask, DependencyRule: constants.ProjectManagementDependencyRuleFS}

	dependency, err := logic.Create(context.Background(), actor, input)
	if err != nil || dependency == nil || dependency.Id == 0 {
		t.Fatalf("valid dependency rejected: %#v, %v", dependency, err)
	}
	postDependency, err := logic.Create(context.Background(), actor, dto.ProjectManagementTaskDependencyDto{TaskId: 2, DependencyTaskId: 4, DependencyType: constants.ProjectManagementDependencyTypePostTask, DependencyRule: constants.ProjectManagementDependencyRuleFS})
	if err != nil {
		t.Fatalf("mixed post-task dependency rejected: %v", err)
	}
	if _, err = logic.Update(context.Background(), actor, dto.ProjectManagementTaskDependencyDto{Id: dependency.Id, TaskId: 2, DependencyTaskId: 3, DependencyType: constants.ProjectManagementDependencyTypePreTask, DependencyRule: constants.ProjectManagementDependencyRuleFS}); err == nil {
		t.Fatal("dependency updated to a task from another project")
	}
	if _, err = logic.Update(context.Background(), actor, dto.ProjectManagementTaskDependencyDto{Id: postDependency.Id, TaskId: 2, DependencyTaskId: 1, DependencyType: constants.ProjectManagementDependencyTypePostTask, DependencyRule: constants.ProjectManagementDependencyRuleFS}); err == nil {
		t.Fatal("dependency update created a cycle")
	}
	if _, err = logic.Create(context.Background(), actor, dto.ProjectManagementTaskDependencyDto{TaskId: 1, DependencyTaskId: 4, DependencyType: constants.ProjectManagementDependencyTypePreTask, DependencyRule: constants.ProjectManagementDependencyRuleFS}); err == nil {
		t.Fatal("mixed pre/post dependency cycle accepted")
	}
	if _, err = logic.Create(context.Background(), actor, input); err == nil {
		t.Fatal("duplicate dependency accepted")
	}
	if _, err = logic.Create(context.Background(), actor, dto.ProjectManagementTaskDependencyDto{TaskId: 1, DependencyTaskId: 2, DependencyType: constants.ProjectManagementDependencyTypePreTask, DependencyRule: constants.ProjectManagementDependencyRuleFS}); err == nil {
		t.Fatal("dependency cycle accepted")
	}
	if _, err = logic.Create(context.Background(), actor, dto.ProjectManagementTaskDependencyDto{TaskId: 2, DependencyTaskId: 1, DependencyType: constants.ProjectManagementDependencyTypePostTask, DependencyRule: constants.ProjectManagementDependencyRuleFS}); err == nil {
		t.Fatal("post-task dependency cycle accepted")
	}
	if _, err = logic.Create(context.Background(), actor, dto.ProjectManagementTaskDependencyDto{TaskId: 2, DependencyTaskId: 3, DependencyType: constants.ProjectManagementDependencyTypePreTask, DependencyRule: constants.ProjectManagementDependencyRuleFS}); err == nil {
		t.Fatal("cross-project dependency accepted")
	}
	if err = db.Exec("UPDATE project_management_tasks SET stage_id = 101 WHERE id = 2").Error; err != nil {
		t.Fatal(err)
	}
	if _, err = logic.Create(context.Background(), actor, dto.ProjectManagementTaskDependencyDto{TaskId: 2, DependencyTaskId: 1, DependencyType: constants.ProjectManagementDependencyTypePreTask, DependencyRule: constants.ProjectManagementDependencyRuleFF}); err == nil {
		t.Fatal("dependency added after successor left pending")
	}
	if err = db.Exec("UPDATE project_management_tasks SET stage_id = 100 WHERE id = 2").Error; err != nil {
		t.Fatal(err)
	}
	if err = logic.repos.ProjectManagementDependencyRepo.Transaction(context.Background(), func(tx *gorm.DB) error {
		return logic.SyncWithTx(context.Background(), tx, actor, 2, []*dto.ProjectManagementTaskDependencyDto{{
			Id: dependency.Id, TaskId: 2, DependencyTaskId: 1,
			DependencyType: constants.ProjectManagementDependencyTypePreTask,
			DependencyRule: constants.ProjectManagementDependencyRuleFF,
		}})
	}); err != nil {
		t.Fatalf("dependency sync failed: %v", err)
	}
	var count int64
	if err = db.Model(&struct{ Id int64 }{}).Table("project_management_task_dependencies").Where("task_id = ? AND is_del = 0", 2).Count(&count).Error; err != nil || count != 1 {
		t.Fatalf("dependency sync did not replace set: count=%d err=%v", count, err)
	}
	if err = db.Exec("INSERT INTO project_management_task_dependencies (task_id, dependency_task_id, dependency_type, dependency_rule) VALUES (4, 2, 'PreTask', 'FS')").Error; err != nil {
		t.Fatal(err)
	}
	if err = logic.repos.ProjectManagementDependencyRepo.Transaction(context.Background(), func(tx *gorm.DB) error {
		return logic.SyncWithTx(context.Background(), tx, actor, 2, nil)
	}); err != nil {
		t.Fatalf("empty dependency sync failed: %v", err)
	}
	if err = db.Model(&struct{ Id int64 }{}).Table("project_management_task_dependencies").Where("task_id = ? AND is_del = 0", 2).Count(&count).Error; err != nil || count != 0 {
		t.Fatalf("empty dependency sync did not clear task dependencies: count=%d err=%v", count, err)
	}
	if err = db.Model(&struct{ Id int64 }{}).Table("project_management_task_dependencies").Where("task_id = ? AND dependency_task_id = ? AND is_del = 0", 4, 2).Count(&count).Error; err != nil || count != 1 {
		t.Fatalf("empty dependency sync removed a reverse dependency: count=%d err=%v", count, err)
	}
}
