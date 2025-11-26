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

// GetOneById get one by id
func (repo *OrganizationRepo) GetOneById(ctx context.Context, id int64, typ constants.UserType) (result *model.Organization, err error) {
	err = repo.Client(ctx).
		Where("id = ? and type = ?", id, typ).
		First(&result).Error
	return result, errorx.WithStack(err)
}
