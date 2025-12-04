package repo

import (
	"context"

	"gorm.io/gorm"

	"zeebee-backend/common/constants"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
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
func (repo *OrganizationRepo) GetOrganization(ctx context.Context, id int64, organizationType constants.UserType) (result *model.Organization, err error) {
	err = repo.Client(ctx).
		Where("id = ? and type = ?", id, organizationType).
		First(&result).Error
	return result, errorx.WithStack(err)
}
