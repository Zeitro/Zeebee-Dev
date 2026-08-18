package logic

import (
	"context"
	"fmt"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/repo"
)

func TestValidateProjectManagementMemberRole(t *testing.T) {
	require.NoError(t, validateProjectManagementMemberRole(constants.ProjectManagementMemberRoleProjectAdmin, false))
	require.NoError(t, validateProjectManagementMemberRole(constants.ProjectManagementMemberRoleTaskExecutor, false))
	require.Error(t, validateProjectManagementMemberRole(constants.ProjectManagementMemberRoleProjectOwner, false))
	require.Error(t, validateProjectManagementMemberRole(constants.ProjectManagementMemberRoleTaskCreator, false))
	require.NoError(t, validateProjectManagementMemberRole(constants.ProjectManagementMemberRoleProjectOwner, true))
}

func TestProjectMemberTransferRemoveAndHistoricalSnapshots(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())), &gorm.Config{})
	require.NoError(t, err)
	for _, statement := range []string{
		"CREATE TABLE users (id INTEGER PRIMARY KEY, is_del INTEGER NOT NULL DEFAULT 0, email TEXT NOT NULL DEFAULT '')",
		"CREATE TABLE organization_users (id INTEGER PRIMARY KEY, is_del INTEGER NOT NULL DEFAULT 0, organization_id INTEGER NOT NULL DEFAULT 0, user_id INTEGER NOT NULL DEFAULT 0, full_name TEXT NOT NULL DEFAULT '', status TEXT NOT NULL DEFAULT '')",
		"CREATE TABLE project_management_projects (id INTEGER PRIMARY KEY, created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, is_del INTEGER NOT NULL DEFAULT 0, organization_id INTEGER NOT NULL DEFAULT 0, property_id INTEGER NOT NULL DEFAULT 0, name TEXT NOT NULL DEFAULT '', description TEXT NOT NULL DEFAULT '', cover_url TEXT NOT NULL DEFAULT '', owner_id INTEGER NOT NULL DEFAULT 0)",
		"CREATE TABLE project_management_project_members (id INTEGER PRIMARY KEY, created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, is_del INTEGER NOT NULL DEFAULT 0, project_id INTEGER NOT NULL DEFAULT 0, organization_user_id INTEGER NOT NULL DEFAULT 0, role TEXT NOT NULL DEFAULT '')",
		"CREATE TABLE project_management_tasks (id INTEGER PRIMARY KEY, is_del INTEGER NOT NULL DEFAULT 0, project_id INTEGER NOT NULL DEFAULT 0, organization_id INTEGER NOT NULL DEFAULT 0, stage_id INTEGER NOT NULL DEFAULT 0, creator_id INTEGER NOT NULL DEFAULT 0, creator_name TEXT NOT NULL DEFAULT '', title TEXT NOT NULL DEFAULT '')",
		"CREATE TABLE project_management_task_assignees (id INTEGER PRIMARY KEY, is_del INTEGER NOT NULL DEFAULT 0, task_id INTEGER NOT NULL DEFAULT 0, organization_user_id INTEGER NOT NULL DEFAULT 0, organization_user_name TEXT NOT NULL DEFAULT '')",
		"CREATE TABLE operations (id INTEGER PRIMARY KEY AUTOINCREMENT, created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, is_del INTEGER NOT NULL DEFAULT 0, entity_id INTEGER NOT NULL, entity_type TEXT NOT NULL, operator_id INTEGER NOT NULL, operator_name TEXT NOT NULL DEFAULT '', operator_type TEXT NOT NULL, operation TEXT NOT NULL, note TEXT NOT NULL DEFAULT '', metadata JSON NOT NULL DEFAULT '{}', organization_id INTEGER NOT NULL)",
	} {
		require.NoError(t, db.Exec(statement).Error)
	}
	require.NoError(t, db.Exec("INSERT INTO users (id, email) VALUES (100, 'owner@example.test'), (101, 'admin@example.test'), (102, 'candidate@example.test'), (103, 'other@example.test')").Error)
	require.NoError(t, db.Exec("INSERT INTO organization_users (id, organization_id, user_id, full_name, status) VALUES (10, 1, 100, 'Owner', ?), (11, 1, 101, 'Admin', ?), (12, 1, 102, 'Candidate', ?), (13, 2, 103, 'Other organization', ?)", constants.OrganizationUserStatusActive, constants.OrganizationUserStatusActive, constants.OrganizationUserStatusActive, constants.OrganizationUserStatusActive).Error)
	require.NoError(t, db.Exec("INSERT INTO project_management_projects (id, organization_id, name, owner_id) VALUES (1, 1, 'Project', 10)").Error)
	require.NoError(t, db.Exec("INSERT INTO project_management_project_members (id, project_id, organization_user_id, role) VALUES (1, 1, 10, ?), (2, 1, 11, ?)", constants.ProjectManagementMemberRoleProjectOwner, constants.ProjectManagementMemberRoleProjectAdmin).Error)
	require.NoError(t, db.Exec("INSERT INTO project_management_tasks (id, project_id, organization_id, stage_id, creator_id, creator_name, title) VALUES (1, 1, 1, 1, 10, 'Former owner', 'Historical task')").Error)
	require.NoError(t, db.Exec("INSERT INTO project_management_task_assignees (id, task_id, organization_user_id, organization_user_name) VALUES (1, 1, 10, 'Former owner')").Error)

	logic := NewProjectManagementLogic(NewBaseLogic(repo.NewAllRepos(db), nil, nil, dto.JwtConfig{}))
	actor := &dto.ProjectManagementActor{OrganizationId: 1, OrganizationUserId: 10, Name: "Owner"}
	candidates, err := logic.repos.ProjectManagementProjectRepo.CandidateMembers(context.Background(), 1, 1, "")
	require.NoError(t, err)
	require.Len(t, candidates, 1)
	require.Equal(t, int64(12), candidates[0].OrganizationUserId)

	require.NoError(t, logic.TransferProjectOwner(context.Background(), actor, 1, 11))
	require.NoError(t, logic.RemoveProjectMember(context.Background(), &dto.ProjectManagementActor{OrganizationId: 1, OrganizationUserId: 11, Name: "Admin"}, 1, 10))

	var project model.ProjectManagementProject
	require.NoError(t, db.First(&project, 1).Error)
	require.Equal(t, int64(11), project.OwnerId)
	var formerOwner model.ProjectManagementProjectMember
	require.ErrorIs(t, db.Where("project_id = ? AND organization_user_id = ?", 1, 10).First(&formerOwner).Error, gorm.ErrRecordNotFound)
	var task model.ProjectManagementTask
	require.NoError(t, db.First(&task, 1).Error)
	require.Equal(t, int64(10), task.CreatorId)
	require.Equal(t, "Former owner", task.CreatorName)
	var assignee model.ProjectManagementTaskAssignee
	require.NoError(t, db.First(&assignee, 1).Error)
	require.Equal(t, int64(10), assignee.OrganizationUserId)
	require.Equal(t, "Former owner", assignee.OrganizationUserName)
}
