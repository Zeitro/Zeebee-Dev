package repo

import (
	"context"

	"gorm.io/gorm"

	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"
)

type OrganizationRepo struct {
	*gormx.BaseRepo[*model.Organization]
	m *model.Organization
}

func NewOrganizationRepo(db *gorm.DB) *OrganizationRepo {
	return &OrganizationRepo{
		BaseRepo: gormx.NewBaseRepo[*model.Organization](db),
		m:        new(model.Organization),
	}
}

// UpdateOwnerIdWithTx update owner id with transaction
func (r *OrganizationRepo) UpdateOwnerIdWithTx(ctx context.Context, tx *gorm.DB, id int64, ownerId int64) (err error) {
	err = tx.WithContext(ctx).
		Model(&model.Organization{}).
		Where("id = ?", id).
		Update("owner_id", ownerId).Error
	return errorx.WithStack(err)
}

// UpdateOwnerOrganizationStatusWithTx update owner's organization status with transaction
func (r *OrganizationRepo) UpdateOwnerOrganizationStatusWithTx(ctx context.Context, tx *gorm.DB, ownerId int64, status constants.OrganizationStatus) (err error) {
	err = tx.WithContext(ctx).
		Model(&model.Organization{}).
		Where("owner_id = ?", ownerId).
		Update("status", status).Error
	return errorx.WithStack(err)
}

// GetOrganization return organization by id and type
func (repo *OrganizationRepo) GetOrganization(ctx context.Context, id int64, organizationType constants.UserType, withPreloads bool) (result *model.Organization, err error) {
	query := repo.Client(ctx).
		Where("id = ? and type = ?", id, organizationType)
	if withPreloads {
		query = query.Preload("OrganizationEmails")
	}
	err = query.First(&result).Error
	return result, errorx.WithStack(err)
}

// GetOrganizationWithTx return organization with transaction
func (repo *OrganizationRepo) GetOrganizationWithTx(ctx context.Context, tx *gorm.DB, id int64, organizationType constants.UserType, withPreloads bool) (result *model.Organization, err error) {
	query := tx.WithContext(ctx).Model(&model.Organization{}).
		Where("id = ? and type = ?", id, organizationType)
	if withPreloads {
		query = query.Preload("OrganizationEmails")
	}
	err = query.First(&result).Error
	return result, errorx.WithStack(err)
}

// All
func (repo *OrganizationRepo) All(ctx context.Context, organizationType constants.UserType) (result []*model.Organization, err error) {
	err = repo.Client(ctx).
		Where("type = ? and status = ?", organizationType, constants.OrganizationStatusActive).
		Order("id asc").
		Find(&result).Error
	return result, errorx.WithStack(err)
}
