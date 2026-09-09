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

func TestProjectManagementCandidateSearchesByID(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	for _, statement := range []string{
		"CREATE TABLE project_management_tasks (id integer primary key, is_del integer not null default 0, organization_id integer not null, project_id integer not null, phase text not null, title text not null)",
		"CREATE TABLE users (id integer primary key, is_del integer not null default 0, email text not null)",
		"CREATE TABLE organization_users (id integer primary key, is_del integer not null default 0, organization_id integer not null, user_id integer not null, status text not null, full_name text not null)",
		"CREATE TABLE project_management_project_members (id integer primary key, is_del integer not null default 0, project_id integer not null, organization_user_id integer not null, role text not null)",
		"INSERT INTO project_management_tasks (id, organization_id, project_id, phase, title) VALUES (101, 1, 10, 'Pending', 'Task One'), (102, 1, 10, 'Pending', 'Task Two')",
		"INSERT INTO users (id, email) VALUES (1001, 'member@example.com')",
		"INSERT INTO organization_users (id, organization_id, user_id, status, full_name) VALUES (501, 1, 1001, 'Active', 'Member One')",
	} {
		if err := db.Exec(statement).Error; err != nil {
			t.Fatal(err)
		}
	}

	taskRepo := NewProjectManagementTaskRepo(db)
	tasks, err := taskRepo.CandidateByProject(context.Background(), 1, dto.ProjectManagementTaskCandidateDto{ProjectId: 10, Keyword: "101"})
	if err != nil || len(tasks) != 1 || tasks[0].Id != 101 {
		t.Fatalf("task ID search = %#v, %v", tasks, err)
	}

	projectRepo := NewProjectManagementProjectRepo(db)
	members, err := projectRepo.CandidateMembers(context.Background(), 1, 10, "501", constants.ProjectManagementMemberCandidatePurposeAddMember)
	if err != nil || len(members) != 1 || members[0].OrganizationUserId != 501 {
		t.Fatalf("member ID search = %#v, %v", members, err)
	}
}
