package repo

import (
	"context"

	"gorm.io/gorm"

	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/model"
)

type OrganizationSettingRepo struct {
	*gormx.BaseRepo[*model.OrganizationSetting]
	m *model.OrganizationSetting
}

func NewOrganizationSettingRepo(db *gorm.DB) *OrganizationSettingRepo {
	return &OrganizationSettingRepo{
		BaseRepo: gormx.NewBaseRepo[*model.OrganizationSetting](db),
		m:        new(model.OrganizationSetting),
	}
}

// Detail retrieves an organization template by its ID.
// It optionally preloads the associated BaseSetting if withPreload is true.
func (repo *OrganizationSettingRepo) Detail(ctx context.Context, id int64, withPreload bool) (result *model.OrganizationSetting, err error) {
	query := repo.Client(ctx).Where("id = ?", id)
	if withPreload {
		query = query.Preload("Organization")
	}
	err = query.First(&result).Error
	return result, errorx.WithStack(err)
}

// DetailByPurpose retrieves an organization template by its purpose and organization ID.
func (repo *OrganizationSettingRepo) DetailByPurpose(ctx context.Context, organizationId int64, purpose string, withPreload bool) (result *model.OrganizationSetting, err error) {
	query := repo.Client(ctx).Where("organization_id = ?", organizationId).Where("purpose = ?", purpose)
	if withPreload {
		query = query.Preload("Organization")
	}
	err = query.First(&result).Error
	return result, errorx.WithStack(err)
}

// DeleteByOrganizationIdWithTx
func (repo *OrganizationSettingRepo) DeleteByOrganizationIdWithTx(ctx context.Context, tx *gorm.DB, organizationId int64) (err error) {
	err = tx.WithContext(ctx).
		Where("organization_id = ?", organizationId).
		Delete(&model.OrganizationSetting{}).Error
	return errorx.WithStack(err)
}
