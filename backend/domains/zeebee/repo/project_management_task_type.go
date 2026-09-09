package repo

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/model"
)

type ProjectManagementTaskTypeRepo struct {
	*gormx.BaseRepo[*model.ProjectManagementTaskType]
}

func NewProjectManagementTaskTypeRepo(db *gorm.DB) *ProjectManagementTaskTypeRepo {
	return &ProjectManagementTaskTypeRepo{BaseRepo: gormx.NewBaseRepo[*model.ProjectManagementTaskType](db)}
}

func (r *ProjectManagementTaskTypeRepo) FindByOrganizationWithTx(ctx context.Context, tx *gorm.DB, id, organizationId int64) (*model.ProjectManagementTaskType, error) {
	return r.FindOneByColumnsWithTx(ctx, tx, map[string]any{"id": id, "organization_id": organizationId})
}

func (r *ProjectManagementTaskTypeRepo) FindByOrganization(ctx context.Context, id, organizationId int64) (*model.ProjectManagementTaskType, error) {
	return r.FindOneByColumns(ctx, map[string]any{"id": id, "organization_id": organizationId})
}

func (r *ProjectManagementTaskTypeRepo) ListByOrganization(ctx context.Context, organizationId int64) ([]*model.ProjectManagementTaskType, error) {
	taskTypes := make([]*model.ProjectManagementTaskType, 0)
	err := r.Client(ctx).Where("organization_id = ?", organizationId).Order("sort_order ASC, id ASC").Find(&taskTypes).Error
	return taskTypes, errorx.WithStack(err)
}

func (r *ProjectManagementTaskTypeRepo) ListByOrganizationWithTx(ctx context.Context, tx *gorm.DB, organizationId int64) ([]*model.ProjectManagementTaskType, error) {
	taskTypes := make([]*model.ProjectManagementTaskType, 0)
	err := tx.WithContext(ctx).Where("organization_id = ?", organizationId).Order("sort_order ASC, id ASC").Find(&taskTypes).Error
	return taskTypes, errorx.WithStack(err)
}

func (r *ProjectManagementTaskTypeRepo) LockByOrganizationWithTx(ctx context.Context, tx *gorm.DB, organizationId int64) ([]*model.ProjectManagementTaskType, error) {
	taskTypes := make([]*model.ProjectManagementTaskType, 0)
	err := tx.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).Where("organization_id = ?", organizationId).Order("id ASC").Find(&taskTypes).Error
	return taskTypes, errorx.WithStack(err)
}

func (r *ProjectManagementTaskTypeRepo) UpdateSortOrderWithTx(ctx context.Context, tx *gorm.DB, id, organizationId int64, sortOrder int) error {
	return r.UpdateColumnsWithTx(ctx, tx, map[string]any{"sort_order": sortOrder}, map[string]any{"id": id, "organization_id": organizationId}, nil)
}

func (r *ProjectManagementTaskTypeRepo) DeleteByOrganizationWithTx(ctx context.Context, tx *gorm.DB, id, organizationId int64) error {
	_, err := r.DeleteWithTx(ctx, tx, map[string]any{"id": id, "organization_id": organizationId})
	return err
}
