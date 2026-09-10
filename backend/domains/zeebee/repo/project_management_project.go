package repo

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
)

type ProjectManagementProjectRepo struct {
	*gormx.BaseRepo[*model.ProjectManagementProject]
	memberRepo *gormx.BaseRepo[*model.ProjectManagementProjectMember]
}

type ProjectManagementProjectMemberEntity struct {
	Id                 int64                                 `gorm:"column:id" json:"id"`
	ProjectId          int64                                 `gorm:"column:project_id" json:"project_id"`
	OrganizationUserId int64                                 `gorm:"column:organization_user_id" json:"organization_user_id"`
	Role               constants.ProjectManagementMemberRole `gorm:"column:role" json:"role"`
	FullName           string                                `gorm:"column:full_name" json:"full_name"`
	Email              string                                `gorm:"column:email" json:"email"`
	CreatedAt          time.Time                             `gorm:"column:created_at" json:"created_at"`
	UpdatedAt          time.Time                             `gorm:"column:updated_at" json:"updated_at"`
}

func NewProjectManagementProjectRepo(db *gorm.DB) *ProjectManagementProjectRepo {
	return &ProjectManagementProjectRepo{
		BaseRepo:   gormx.NewBaseRepo[*model.ProjectManagementProject](db),
		memberRepo: gormx.NewBaseRepo[*model.ProjectManagementProjectMember](db),
	}
}

func (r *ProjectManagementProjectRepo) Page(ctx context.Context, organizationId int64, input dto.ProjectManagementProjectPageRequestDto) (gormx.PageResponse[*model.ProjectManagementProject], error) {
	query := r.Client(ctx).Where("organization_id = ?", organizationId).Order("id DESC")
	if input.Filter.Keyword != "" {
		query = query.Where("name ILIKE ?", "%"+input.Filter.Keyword+"%")
	}
	if input.Filter.EntityId != 0 {
		query = query.Where("entity_id = ?", input.Filter.EntityId)
	}
	if input.Filter.EntityType != "" {
		query = query.Where("entity_type = ?", input.Filter.EntityType)
	}
	if input.Filter.OrganizationUserId > 0 {
		query = query.Where("id IN (?)", r.memberRepo.Client(ctx).
			Select("project_id").
			Where("organization_user_id = ?", input.Filter.OrganizationUserId))
	}
	return r.BaseRepo.Page(query, gormx.PageRequest{Page: input.Page, Limit: input.Limit})
}

func (r *ProjectManagementProjectRepo) PageMembers(ctx context.Context, projectId int64, input dto.ProjectManagementMemberPageRequestDto) (gormx.PageResponse[*ProjectManagementProjectMemberEntity], error) {
	membersTable := model.ProjectManagementProjectMember{}.TableName()
	organizationUsersTable := model.OrganizationUser{}.TableName()
	usersTable := model.User{}.TableName()
	query := r.Client(ctx).Table(membersTable).
		Select(fmt.Sprintf("%s.id, %s.project_id, %s.organization_user_id, %s.role, %s.full_name, %s.email, %s.created_at, %s.updated_at", membersTable, membersTable, membersTable, membersTable, organizationUsersTable, usersTable, membersTable, membersTable)).
		Joins(fmt.Sprintf("JOIN %s ON %s.id = %s.organization_user_id", organizationUsersTable, organizationUsersTable, membersTable)).
		Joins(fmt.Sprintf("JOIN %s ON %s.id = %s.user_id", usersTable, usersTable, organizationUsersTable)).
		Where(fmt.Sprintf("%s.project_id = ?", membersTable), projectId).
		Order(fmt.Sprintf("CASE %s.role WHEN '%s' THEN 1 WHEN '%s' THEN 2 WHEN '%s' THEN 3 WHEN '%s' THEN 4 ELSE 5 END, %s.id DESC", membersTable, constants.ProjectManagementMemberRoleProjectOwner, constants.ProjectManagementMemberRoleProjectAdmin, constants.ProjectManagementMemberRoleTaskExecutor, constants.ProjectManagementMemberRoleTaskCreator, membersTable))
	if input.Filter.Keyword != "" {
		query = query.Where(fmt.Sprintf("%s.full_name ILIKE ? OR %s.email ILIKE ?", organizationUsersTable, usersTable), "%"+input.Filter.Keyword+"%", "%"+input.Filter.Keyword+"%")
	}
	if input.Filter.Role != "" {
		query = query.Where(fmt.Sprintf("%s.role = ?", membersTable), input.Filter.Role)
	}
	return gormx.Page[*ProjectManagementProjectMemberEntity](query, gormx.PageRequest{Page: input.Page, Limit: input.Limit, CountColumn: fmt.Sprintf("%s.id", membersTable)})
}

func (r *ProjectManagementProjectRepo) DetailMember(ctx context.Context, projectId, memberId int64) (*ProjectManagementProjectMemberEntity, error) {
	membersTable := model.ProjectManagementProjectMember{}.TableName()
	organizationUsersTable := model.OrganizationUser{}.TableName()
	usersTable := model.User{}.TableName()
	member := new(ProjectManagementProjectMemberEntity)
	err := r.Client(ctx).Table(membersTable).
		Select(fmt.Sprintf("%s.id, %s.project_id, %s.organization_user_id, %s.role, %s.full_name, %s.email, %s.created_at, %s.updated_at", membersTable, membersTable, membersTable, membersTable, organizationUsersTable, usersTable, membersTable, membersTable)).
		Joins(fmt.Sprintf("JOIN %s ON %s.id = %s.organization_user_id", organizationUsersTable, organizationUsersTable, membersTable)).
		Joins(fmt.Sprintf("JOIN %s ON %s.id = %s.user_id", usersTable, usersTable, organizationUsersTable)).
		Where(fmt.Sprintf("%s.id = ? AND %s.project_id = ?", membersTable, membersTable), memberId, projectId).
		First(member).Error
	return member, errorx.WithStack(err)
}

func (r *ProjectManagementProjectRepo) CandidateMembers(ctx context.Context, organizationId, projectId int64, keyword string, purpose constants.ProjectManagementMemberCandidatePurpose) ([]*dto.ProjectManagementMemberCandidateEntityDto, error) {
	membersTable := model.ProjectManagementProjectMember{}.TableName()
	organizationUsersTable := model.OrganizationUser{}.TableName()
	usersTable := model.User{}.TableName()
	var query *gorm.DB
	switch purpose {
	case constants.ProjectManagementMemberCandidatePurposeAddMember:
		query = r.Client(ctx).Model(&model.OrganizationUser{}).Table(organizationUsersTable).
			Select(fmt.Sprintf("%s.id AS organization_user_id, %s.full_name, %s.email", organizationUsersTable, organizationUsersTable, usersTable)).
			Joins(fmt.Sprintf("JOIN %s ON %s.id = %s.user_id", usersTable, usersTable, organizationUsersTable)).
			Where(fmt.Sprintf("%s.organization_id = ? AND %s.status = ? AND %s.is_del = 0", organizationUsersTable, organizationUsersTable, usersTable), organizationId, constants.OrganizationUserStatusActive).
			Order(fmt.Sprintf("%s.full_name ASC, %s.id ASC", organizationUsersTable, organizationUsersTable))
		query = query.Where(fmt.Sprintf("NOT EXISTS (SELECT 1 FROM %s WHERE %s.project_id = ? AND %s.organization_user_id = %s.id AND %s.is_del = 0)", membersTable, membersTable, membersTable, organizationUsersTable, membersTable), projectId)
	case constants.ProjectManagementMemberCandidatePurposeTransferOwner, constants.ProjectManagementMemberCandidatePurposeProjectMember:
		query = r.Client(ctx).Table(membersTable).
			Select(fmt.Sprintf("%s.organization_user_id, %s.full_name, %s.email", membersTable, organizationUsersTable, usersTable)).
			Joins(fmt.Sprintf("JOIN %s ON %s.id = %s.organization_user_id", organizationUsersTable, organizationUsersTable, membersTable)).
			Joins(fmt.Sprintf("JOIN %s ON %s.id = %s.user_id", usersTable, usersTable, organizationUsersTable)).
			Where(fmt.Sprintf("%s.project_id = ? AND %s.is_del = 0", membersTable, membersTable), projectId).
			Order(fmt.Sprintf("%s.full_name ASC, %s.id ASC", organizationUsersTable, membersTable))
		if purpose == constants.ProjectManagementMemberCandidatePurposeTransferOwner {
			query = query.Where(fmt.Sprintf("%s.role <> ?", membersTable), constants.ProjectManagementMemberRoleProjectOwner)
		}
	default:
		return nil, errorx.New("invalid project member candidate purpose")
	}
	keyword = strings.TrimSpace(keyword)
	if keyword != "" {
		pattern := "%" + keyword + "%"
		if id, err := strconv.ParseInt(keyword, 10, 64); err == nil && id > 0 {
			query = query.Where(fmt.Sprintf("%s.id = ? OR LOWER(%s.full_name) LIKE LOWER(?) OR LOWER(%s.email) LIKE LOWER(?)", organizationUsersTable, organizationUsersTable, usersTable), id, pattern, pattern)
		} else {
			query = query.Where(fmt.Sprintf("LOWER(%s.full_name) LIKE LOWER(?) OR LOWER(%s.email) LIKE LOWER(?)", organizationUsersTable, usersTable), pattern, pattern)
		}
	}
	candidates := make([]*dto.ProjectManagementMemberCandidateEntityDto, 0)
	err := query.Find(&candidates).Error
	return candidates, errorx.WithStack(err)
}

func (r *ProjectManagementProjectRepo) HasUnfinishedTasksForEntity(ctx context.Context, entityId int64, entityType constants.ProjectManagementProjectEntityType) (bool, error) {
	projectsTable := model.ProjectManagementProject{}.TableName()
	tasksTable := model.ProjectManagementTask{}.TableName()
	stagesTable := model.ProjectManagementWorkflowStage{}.TableName()
	var count int64
	err := r.Client(ctx).Table(projectsTable).
		Joins(fmt.Sprintf("JOIN %s ON %s.project_id = %s.id AND %s.is_del = 0", tasksTable, tasksTable, projectsTable, tasksTable)).
		Joins(fmt.Sprintf("JOIN %s ON %s.id = %s.stage_id AND %s.is_del = 0", stagesTable, stagesTable, tasksTable, stagesTable)).
		Where(fmt.Sprintf("%s.entity_id = ? AND %s.entity_type = ? AND %s.is_del = 0 AND %s.phase <> ?", projectsTable, projectsTable, projectsTable, stagesTable), entityId, entityType, constants.ProjectManagementStagePhaseFinished).
		Count(&count).Error
	return count > 0, errorx.WithStack(err)
}

func (r *ProjectManagementProjectRepo) FindByOrganization(ctx context.Context, id, organizationId int64) (*model.ProjectManagementProject, error) {
	return r.FindOneByColumns(ctx, map[string]any{"id": id, "organization_id": organizationId})
}

func (r *ProjectManagementProjectRepo) FindMember(ctx context.Context, projectId, organizationUserId int64) (*model.ProjectManagementProjectMember, error) {
	return r.memberRepo.FindOneByColumns(ctx, map[string]any{"project_id": projectId, "organization_user_id": organizationUserId})
}

func (r *ProjectManagementProjectRepo) FindMemberById(ctx context.Context, projectId, memberId int64) (*model.ProjectManagementProjectMember, error) {
	return r.memberRepo.FindOneByColumns(ctx, map[string]any{"id": memberId, "project_id": projectId})
}

func (r *ProjectManagementProjectRepo) LockWithTx(ctx context.Context, tx *gorm.DB, id int64) (*model.ProjectManagementProject, error) {
	project := new(model.ProjectManagementProject)
	err := tx.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).First(project, id).Error
	return project, errorx.WithStack(err)
}

func (r *ProjectManagementProjectRepo) LockByOrganizationWithTx(ctx context.Context, tx *gorm.DB, id, organizationId int64) (*model.ProjectManagementProject, error) {
	project := new(model.ProjectManagementProject)
	err := tx.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ? AND organization_id = ?", id, organizationId).First(project).Error
	return project, errorx.WithStack(err)
}

func (r *ProjectManagementProjectRepo) LockMemberWithTx(ctx context.Context, tx *gorm.DB, projectId, organizationUserId int64) (*model.ProjectManagementProjectMember, error) {
	member := new(model.ProjectManagementProjectMember)
	err := tx.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("project_id = ? AND organization_user_id = ?", projectId, organizationUserId).First(member).Error
	return member, errorx.WithStack(err)
}

func (r *ProjectManagementProjectRepo) LockMemberByIdWithTx(ctx context.Context, tx *gorm.DB, projectId, memberId int64) (*model.ProjectManagementProjectMember, error) {
	member := new(model.ProjectManagementProjectMember)
	err := tx.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ? AND project_id = ?", memberId, projectId).First(member).Error
	return member, errorx.WithStack(err)
}

func (r *ProjectManagementProjectRepo) DeleteByOrganizationUserWithTx(ctx context.Context, tx *gorm.DB, organizationUserId int64) error {
	_, err := r.memberRepo.DeleteWithTx(ctx, tx, map[string]any{"organization_user_id": organizationUserId})
	return err
}

func (r *ProjectManagementProjectRepo) CreateMemberWithTx(ctx context.Context, tx *gorm.DB, member *model.ProjectManagementProjectMember) error {
	_, err := r.memberRepo.CreateWithTx(ctx, tx, member, false)
	return err
}

func (r *ProjectManagementProjectRepo) UpdateMemberRoleWithTx(ctx context.Context, tx *gorm.DB, member *model.ProjectManagementProjectMember, role constants.ProjectManagementMemberRole) error {
	member.Role = role
	return r.memberRepo.UpdateColumnsWithTx(ctx, tx, map[string]any{"role": role}, map[string]any{"id": member.Id, "project_id": member.ProjectId}, nil)
}

func (r *ProjectManagementProjectRepo) UpdateOwnerWithTx(ctx context.Context, tx *gorm.DB, projectId, ownerId int64, ownerName string) error {
	return r.UpdateColumnsWithTx(ctx, tx, map[string]any{"owner_id": ownerId, "owner_name": ownerName}, map[string]any{"id": projectId}, nil)
}

func (r *ProjectManagementProjectRepo) DeleteMemberWithTx(ctx context.Context, tx *gorm.DB, member *model.ProjectManagementProjectMember) error {
	_, err := r.memberRepo.DeleteWithTx(ctx, tx, map[string]any{"id": member.Id, "project_id": member.ProjectId})
	return err
}

func (r *ProjectManagementProjectRepo) DeleteProjectWithTx(ctx context.Context, tx *gorm.DB, projectId int64) error {
	if _, err := r.memberRepo.DeleteWithTx(ctx, tx, map[string]any{"project_id": projectId}); err != nil {
		return err
	}
	_, err := r.DeleteByIDWithTx(ctx, tx, projectId)
	return err
}
