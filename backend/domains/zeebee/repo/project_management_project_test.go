package repo

import (
	"context"
	"testing"

	"zeebee-backend/domains/zeebee/constants"
)

func TestProjectManagementProjectRepoCandidates(t *testing.T) {
	db := newRepoTestDB(t)
	execTestSQL(t, db, []string{
		"CREATE TABLE users (id integer primary key, is_del integer not null default 0, email text not null)",
		"CREATE TABLE organization_users (id integer primary key, is_del integer not null default 0, organization_id integer not null, user_id integer not null, status text not null, full_name text not null)",
		"CREATE TABLE project_management_project_members (id integer primary key, is_del integer not null default 0, project_id integer not null, organization_user_id integer not null, role text not null)",
		"INSERT INTO users (id, email) VALUES (1001, 'member@example.com')",
		"INSERT INTO organization_users (id, organization_id, user_id, status, full_name) VALUES (501, 1, 1001, 'Active', 'Member One')",
	})
	members, err := NewProjectManagementProjectRepo(db).CandidateMembers(context.Background(), 1, 10, "501", constants.ProjectManagementMemberCandidatePurposeAddMember)
	if err != nil || len(members) != 1 || members[0].OrganizationUserId != 501 {
		t.Fatalf("member ID search = %#v, %v", members, err)
	}
}
