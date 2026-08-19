package repo

import (
	"context"

	"gorm.io/gorm"

	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
)

type ProjectManagementTaskTypeRepo struct {
	*gormx.BaseRepo[*model.ProjectManagementTaskType]
}

func NewProjectManagementTaskTypeRepo(db *gorm.DB) *ProjectManagementTaskTypeRepo {
	return &ProjectManagementTaskTypeRepo{BaseRepo: gormx.NewBaseRepo[*model.ProjectManagementTaskType](db)}
}

func (r *ProjectManagementTaskTypeRepo) FindEnabledByOrganization(ctx context.Context, id, organizationId int64) (*model.ProjectManagementTaskType, error) {
	taskType := new(model.ProjectManagementTaskType)
	err := r.Client(ctx).Where("id = ? AND organization_id = ? AND enabled = true", id, organizationId).First(taskType).Error
	return taskType, errorx.WithStack(err)
}

func (r *ProjectManagementTaskTypeRepo) FindEnabledByOrganizationWithTx(ctx context.Context, tx *gorm.DB, id, organizationId int64) (*model.ProjectManagementTaskType, error) {
	taskType := new(model.ProjectManagementTaskType)
	err := tx.WithContext(ctx).Where("id = ? AND organization_id = ? AND enabled = ?", id, organizationId, true).First(taskType).Error
	return taskType, errorx.WithStack(err)
}

func (r *ProjectManagementTaskTypeRepo) FindByOrganization(ctx context.Context, id, organizationId int64) (*model.ProjectManagementTaskType, error) {
	taskType := new(model.ProjectManagementTaskType)
	err := r.Client(ctx).Where("id = ? AND organization_id = ?", id, organizationId).First(taskType).Error
	return taskType, errorx.WithStack(err)
}

func (r *ProjectManagementTaskTypeRepo) Page(ctx context.Context, organizationId int64, input dto.ProjectManagementTaskTypePageRequestDto) (gormx.PageResponse[*model.ProjectManagementTaskType], error) {
	query := r.Client(ctx).Where("organization_id = ?", organizationId).Order("id DESC")
	if input.Filter.Keyword != "" {
		query = query.Where("name ILIKE ?", "%"+input.Filter.Keyword+"%")
	}
	if input.Filter.Enabled != nil {
		query = query.Where("enabled = ?", *input.Filter.Enabled)
	}
	return r.BaseRepo.Page(query, gormx.PageRequest{Page: input.Page, Limit: input.Limit})
}

func (r *ProjectManagementTaskTypeRepo) UpdateByOrganization(ctx context.Context, id, organizationId int64, taskType *model.ProjectManagementTaskType) error {
	err := r.Client(ctx).Model(taskType).
		Where("id = ? AND organization_id = ?", id, organizationId).
		Select("name", "description", "color", "enabled").Updates(taskType).Error
	return errorx.WithStack(err)
}

func (r *ProjectManagementTaskTypeRepo) DeleteByOrganization(ctx context.Context, id, organizationId int64) error {
	err := r.Client(ctx).Where("id = ? AND organization_id = ?", id, organizationId).Delete(&model.ProjectManagementTaskType{}).Error
	return errorx.WithStack(err)
}
