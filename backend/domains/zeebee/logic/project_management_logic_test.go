package logic

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
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

func TestValidateProjectWorkflowInvariant(t *testing.T) {
	projectID := int64(10)
	workflows := []*model.ProjectManagementWorkflow{{Base: model.Base{Id: 1}, EntityType: constants.ProjectManagementOwnerTypeProject, EntityId: projectID, IsDefault: true}}
	stages := []*model.ProjectManagementStage{
		{Base: model.Base{Id: 11}, EntityType: constants.ProjectManagementOwnerTypeProject, EntityId: projectID, Phase: constants.ProjectManagementStagePhasePending},
		{Base: model.Base{Id: 12}, EntityType: constants.ProjectManagementOwnerTypeProject, EntityId: projectID, Phase: constants.ProjectManagementStagePhaseFinished},
	}
	links := []*model.ProjectManagementWorkflowStage{{WorkflowId: 1, StageId: 11}, {WorkflowId: 1, StageId: 12}}
	if err := validateProjectWorkflowInvariant(projectID, workflows, stages, links); err != nil {
		t.Fatalf("valid project workflow rejected: %v", err)
	}
	workflows = append(workflows, &model.ProjectManagementWorkflow{Base: model.Base{Id: 2}, EntityType: constants.ProjectManagementOwnerTypeProject, EntityId: projectID, IsDefault: true})
	if err := validateProjectWorkflowInvariant(projectID, workflows, stages, links); err == nil {
		t.Fatal("duplicate default workflow accepted")
	}
	if err := validateProjectWorkflowInvariant(projectID, workflows[:1], stages, append(links, &model.ProjectManagementWorkflowStage{WorkflowId: 1, StageId: 99})); err == nil {
		t.Fatal("foreign workflow-stage link accepted")
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
		"CREATE TABLE project_management_projects (id integer primary key, created_at datetime not null default CURRENT_TIMESTAMP, updated_at datetime not null default CURRENT_TIMESTAMP, is_del integer not null default 0, organization_id integer not null, entity_id integer not null default 0, entity_type text not null default '', name text not null default '', description text not null default '', cover_url text not null default '', owner_id integer not null default 0, owner_name text not null default '')",
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

func TestProjectManagementTaskOverdue(t *testing.T) {
	now := time.Date(2026, 8, 18, 12, 0, 0, 0, time.UTC)
	cases := []struct {
		name    string
		endAt   *sql.NullTime
		phase   constants.ProjectManagementStagePhase
		overdue bool
		seconds int64
	}{
		{"no deadline", nil, constants.ProjectManagementStagePhaseInProgress, false, 0},
		{"not due", &sql.NullTime{Time: now.Add(time.Hour), Valid: true}, constants.ProjectManagementStagePhaseInProgress, false, 0},
		{"finished", &sql.NullTime{Time: now.Add(-time.Hour), Valid: true}, constants.ProjectManagementStagePhaseFinished, false, 0},
		{"reopened after deadline day", &sql.NullTime{Time: now.AddDate(0, 0, -1).Truncate(24 * time.Hour), Valid: true}, constants.ProjectManagementStagePhasePending, true, int64(12 * time.Hour.Seconds())},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			overdue, seconds := projectManagementTaskOverdue(&model.ProjectManagementTask{EndAt: tc.endAt}, tc.phase, now)
			if overdue != tc.overdue || seconds != tc.seconds {
				t.Fatalf("got overdue=%t seconds=%d", overdue, seconds)
			}
		})
	}
	endAt := &sql.NullTime{Time: time.Date(2026, time.May, 3, 0, 0, 0, 0, time.UTC), Valid: true}
	if overdue, _ := projectManagementTaskOverdue(&model.ProjectManagementTask{EndAt: endAt}, constants.ProjectManagementStagePhaseInProgress, time.Date(2026, time.May, 3, 23, 59, 59, 0, time.UTC)); overdue {
		t.Fatal("task became overdue before the day after its deadline")
	}
	if overdue, _ := projectManagementTaskOverdue(&model.ProjectManagementTask{EndAt: endAt}, constants.ProjectManagementStagePhaseInProgress, time.Date(2026, time.May, 4, 0, 0, 0, 0, time.UTC)); !overdue {
		t.Fatal("task was not overdue at the start of the day after its deadline")
	}
}
