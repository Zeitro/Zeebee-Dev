package logic

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/repo"
)

func TestProjectManagementNotFound(t *testing.T) {
	if got := projectManagementNotFound(gorm.ErrRecordNotFound, "missing"); got.Error() != "missing" {
		t.Fatalf("record-not-found mapping = %v", got)
	}

	sqlErr := errors.New("database connection lost")
	if got := projectManagementNotFound(sqlErr, "missing"); got != sqlErr {
		t.Fatalf("query failure was hidden: %v", got)
	}
}

func TestProjectManagementAccessScopeAndRemovedMember(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	for _, statement := range []string{
		"CREATE TABLE organization_users (id integer primary key, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, organization_id integer not null, user_id integer not null, full_name text not null default '', status text not null default '')",
		"CREATE TABLE users (id integer primary key, is_del integer not null default 0)",
		"CREATE TABLE organization_user_and_roles (organization_user_id integer not null, role_id integer not null)",
		"CREATE TABLE project_management_projects (id integer primary key, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, organization_id integer not null, property_id integer not null default 0, name text not null default '', description text not null default '', cover_url text not null default '', owner_id integer not null default 0)",
		"CREATE TABLE project_management_project_members (id integer primary key, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, project_id integer not null, organization_user_id integer not null, role text not null default 'executor')",
	} {
		if err := db.Exec(statement).Error; err != nil {
			t.Fatal(err)
		}
	}
	for _, statement := range []string{
		"INSERT INTO organization_users (id, organization_id, user_id, full_name) VALUES (20, 1, 99, 'Member')",
		"INSERT INTO project_management_projects (id, organization_id, name) VALUES (10, 1, 'Project')",
		"INSERT INTO project_management_project_members (id, project_id, organization_user_id) VALUES (30, 10, 20)",
	} {
		if err := db.Exec(statement).Error; err != nil {
			t.Fatal(err)
		}
	}
	logic := NewProjectManagementLogic(NewBaseLogic(repo.NewAllRepos(db), nil, nil, dto.JwtConfig{}))
	actor, err := logic.RequireActor(context.Background(), 1, 99)
	if err != nil || actor.OrganizationUserId != 20 || actor.Name != "Member" {
		t.Fatalf("organization actor lookup failed: %#v, %v", actor, err)
	}
	if _, err = logic.RequireProject(context.Background(), 10, 2); err == nil {
		t.Fatal("cross-organization project accepted")
	}
	if _, err = logic.RequireMember(context.Background(), 10, 20); err != nil {
		t.Fatalf("project member rejected: %v", err)
	}
	if err = db.Exec("UPDATE project_management_project_members SET is_del = 1 WHERE id = 30").Error; err != nil {
		t.Fatal(err)
	}
	if _, err = logic.RequireMember(context.Background(), 10, 20); err == nil {
		t.Fatal("removed project member retained access")
	}
}
