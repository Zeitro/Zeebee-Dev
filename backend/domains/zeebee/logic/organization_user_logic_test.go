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

func setupOrganizationUserDeleteLogic(t *testing.T) (*OrganizationUserLogic, *gorm.DB) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())), &gorm.Config{})
	require.NoError(t, err)
	for _, statement := range []string{
		`CREATE TABLE organizations (id INTEGER PRIMARY KEY, is_del INTEGER NOT NULL DEFAULT 0, owner_id INTEGER NOT NULL DEFAULT 0, type TEXT NOT NULL DEFAULT '')`,
		`CREATE TABLE users (id INTEGER PRIMARY KEY, is_del INTEGER NOT NULL DEFAULT 0)`,
		`CREATE TABLE roles (id INTEGER PRIMARY KEY, is_del INTEGER NOT NULL DEFAULT 0)`,
		`CREATE TABLE organization_users (id INTEGER PRIMARY KEY, is_del INTEGER NOT NULL DEFAULT 0, organization_id INTEGER NOT NULL DEFAULT 0, user_id INTEGER NOT NULL DEFAULT 0, full_name TEXT NOT NULL DEFAULT '', status TEXT NOT NULL DEFAULT '')`,
		`CREATE TABLE organization_user_and_roles (organization_user_id INTEGER NOT NULL, role_id INTEGER NOT NULL)`,
		`CREATE TABLE project_management_project_members (id INTEGER PRIMARY KEY, is_del INTEGER NOT NULL DEFAULT 0, project_id INTEGER NOT NULL DEFAULT 0, organization_user_id INTEGER NOT NULL DEFAULT 0, role TEXT NOT NULL DEFAULT '')`,
		`CREATE TABLE project_management_tasks (id INTEGER PRIMARY KEY, is_del INTEGER NOT NULL DEFAULT 0, entity_type TEXT NOT NULL DEFAULT 'Project', entity_id INTEGER NOT NULL DEFAULT 0, organization_id INTEGER NOT NULL DEFAULT 0, stage_id INTEGER NOT NULL DEFAULT 0, creator_id INTEGER NOT NULL DEFAULT 0, creator_name TEXT NOT NULL DEFAULT '', title TEXT NOT NULL DEFAULT '')`,
		`CREATE TABLE project_management_task_assignees (id INTEGER PRIMARY KEY, is_del INTEGER NOT NULL DEFAULT 0, task_id INTEGER NOT NULL DEFAULT 0, stage_id INTEGER NOT NULL DEFAULT 0, organization_user_id INTEGER NOT NULL DEFAULT 0, organization_user_full_name TEXT NOT NULL DEFAULT '')`,
	} {
		require.NoError(t, db.Exec(statement).Error)
	}
	base := NewBaseLogic(repo.NewAllRepos(db), nil, nil, dto.JwtConfig{})
	return NewOrganizationUserLogic(base, nil, nil, map[constants.UserType]string{}), db
}

func TestOrganizationUserDelete_RemovesProjectMembershipOnly(t *testing.T) {
	logic, db := setupOrganizationUserDeleteLogic(t)
	require.NoError(t, db.Exec(`INSERT INTO organizations (id, owner_id, type) VALUES (1, 100, ?)`, constants.UserTypeBorrower).Error)
	require.NoError(t, db.Exec(`INSERT INTO users (id) VALUES (200)`).Error)
	require.NoError(t, db.Exec(`INSERT INTO organization_users (id, organization_id, user_id, full_name, status) VALUES (2, 1, 200, 'Former member', ?)`, constants.OrganizationUserStatusActive).Error)
	require.NoError(t, db.Exec(`INSERT INTO project_management_project_members (id, project_id, organization_user_id, role) VALUES (3, 10, 2, ?)`, constants.ProjectManagementMemberRoleTaskExecutor).Error)
	require.NoError(t, db.Exec(`INSERT INTO project_management_tasks (id, entity_type, entity_id, organization_id, stage_id, creator_id, creator_name, title) VALUES (4, 'Project', 10, 1, 20, 2, 'Former member', 'Historical task')`).Error)
	require.NoError(t, db.Exec(`INSERT INTO project_management_task_assignees (id, task_id, organization_user_id, organization_user_full_name) VALUES (5, 4, 2, 'Former member')`).Error)

	require.NoError(t, logic.Delete(context.Background(), 2, 1, constants.UserTypeBorrower))

	var removedMember model.ProjectManagementProjectMember
	require.NoError(t, db.Where("id = ?", 3).Find(&removedMember).Error)
	require.Zero(t, removedMember.Id)
	require.NoError(t, db.Unscoped().Where("id = ?", 3).First(&removedMember).Error)
	require.NotZero(t, removedMember.IsDel)
	var removedOrganizationUser model.OrganizationUser
	require.NoError(t, db.Where("id = ?", 2).Find(&removedOrganizationUser).Error)
	require.Zero(t, removedOrganizationUser.Id)

	var historicalTask model.ProjectManagementTask
	require.NoError(t, db.Where("id = ?", 4).First(&historicalTask).Error)
	require.Equal(t, int64(2), historicalTask.CreatorId)
	require.Equal(t, "Former member", historicalTask.CreatorName)
	var historicalAssignee model.ProjectManagementTaskAssignee
	require.NoError(t, db.Where("id = ?", 5).First(&historicalAssignee).Error)
	require.Equal(t, int64(2), historicalAssignee.OrganizationUserId)
	require.Equal(t, "Former member", historicalAssignee.OrganizationUserFullName)
}
