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
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/repo"
)

func TestCreateProjectClonesTemplateAndRollsBack(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	for _, statement := range []string{
		"CREATE TABLE borrower_properties (id integer primary key, is_del integer not null default 0, organization_id integer not null)",
		"CREATE TABLE project_management_projects (id integer primary key autoincrement, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, organization_id integer not null, entity_id integer not null default 0, entity_type text not null default '', name text not null default '', description text not null default '', cover_url text not null default '', owner_id integer not null default 0, owner_name text not null default '')",
		"CREATE TABLE project_management_project_members (id integer primary key autoincrement, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, project_id integer not null, organization_user_id integer not null, role text not null)",
		"CREATE TABLE project_management_workflows (id integer primary key autoincrement, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, kind text not null, organization_id integer not null, project_id integer not null default 0, name text not null default '', description text not null default '', is_default boolean not null default false, source_template_id integer not null default 0, source_template_name text not null default '')",
		"CREATE TABLE project_management_workflow_stages (id integer primary key autoincrement, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, workflow_id integer not null, name text not null default '', phase text not null, color text not null default '', sort_order integer not null)",
		"CREATE TABLE project_management_workflow_transitions (id integer primary key autoincrement, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, workflow_id integer not null, name text not null default '', from_stage_id integer not null, to_stage_id integer not null, allowed_actor_types json not null default '[]')",
		"CREATE TABLE operations (id integer primary key autoincrement, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, entity_id integer not null, entity_type text not null, operator_id integer not null, operator_name text not null default '', operator_type text not null, operation text not null, note text not null default '', metadata json not null default '{}', organization_id integer not null)",
	} {
		if err := db.Exec(statement).Error; err != nil {
			t.Fatal(err)
		}
	}
	for _, statement := range []string{
		"INSERT INTO borrower_properties (id, organization_id) VALUES (50, 1), (51, 2)",
		"INSERT INTO project_management_workflows (id, kind, organization_id, name) VALUES (100, 'OrganizationTemplate', 1, 'Inspection')",
		"INSERT INTO project_management_workflow_stages (workflow_id, name, phase, sort_order) VALUES (100, 'Pending', 'Pending', 1)",
	} {
		if err := db.Exec(statement).Error; err != nil {
			t.Fatal(err)
		}
	}
	logic := NewProjectManagementProjectLogic(NewBaseLogic(repo.NewAllRepos(db), nil, nil, dto.JwtConfig{}))
	actor := utils.JwtUser{OrganizationId: 1, OrganizationUserId: 7, Username: "Owner"}
	project, err := logic.Create(context.Background(), actor, dto.ProjectManagementProjectCreateDto{ProjectManagementProjectDto: dto.ProjectManagementProjectDto{Name: "Renovation", Description: "Initial description", CoverUrl: "https://example.com/cover.png", EntityId: 50, EntityType: constants.ProjectManagementEntityTypeBorrowerProperty}, WorkflowTemplateId: 100})
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
	if err := db.Where("project_id = ? AND kind = ?", project.Id, constants.ProjectManagementWorkflowKindProject).First(&workflow).Error; err != nil {
		t.Fatal(err)
	}
	if workflow.SourceTemplateId != 100 || workflow.SourceTemplateName != "Inspection" {
		t.Fatalf("unexpected workflow snapshot: %#v", workflow)
	}
	transactionProject, err := logic.Create(context.Background(), actor, dto.ProjectManagementProjectCreateDto{ProjectManagementProjectDto: dto.ProjectManagementProjectDto{Name: "Transaction project", EntityId: 51, EntityType: "Transaction"}, WorkflowTemplateId: 100})
	if err != nil || transactionProject.EntityId != 51 || transactionProject.EntityType != "Transaction" {
		t.Fatalf("transaction binding was rejected: %#v, %v", transactionProject, err)
	}
	if _, err := logic.Create(context.Background(), actor, dto.ProjectManagementProjectCreateDto{ProjectManagementProjectDto: dto.ProjectManagementProjectDto{Name: "Rejected"}, WorkflowTemplateId: 999}); err == nil {
		t.Fatal("project with missing template was created")
	}
	if _, err := logic.Create(context.Background(), actor, dto.ProjectManagementProjectCreateDto{ProjectManagementProjectDto: dto.ProjectManagementProjectDto{Name: "Foreign property", EntityId: 51, EntityType: constants.ProjectManagementEntityTypeBorrowerProperty}, WorkflowTemplateId: 100}); err == nil {
		t.Fatal("project with a property from another organization was created")
	}
	updated, err := logic.Update(context.Background(), actor, project.Id, dto.ProjectManagementProjectUpdateDto{Name: "Renovation updated"})
	if err != nil {
		t.Fatalf("project update rejected: %v", err)
	}
	if updated.Name != "Renovation updated" || updated.Description != "" || updated.CoverUrl != "" {
		t.Fatalf("unexpected updated project: %#v", updated)
	}
	var persisted model.ProjectManagementProject
	if err := db.First(&persisted, project.Id).Error; err != nil {
		t.Fatal(err)
	}
	if persisted.Name != "Renovation updated" || persisted.Description != "" || persisted.CoverUrl != "" {
		t.Fatalf("zero-value fields were not persisted: %#v", persisted)
	}
	var count int64
	if err := db.Model(&model.ProjectManagementProject{}).Count(&count).Error; err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Fatalf("failed project creation left %d projects", count)
	}
}
