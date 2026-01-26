package repo

import (
	"context"

	"gorm.io/gorm"

	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"
)

type OrganizationTemplateRepo struct {
	*gormx.BaseRepo[*model.OrganizationTemplate]
	m *model.OrganizationTemplate
}

func NewOrganizationTemplateRepo(db *gorm.DB) *OrganizationTemplateRepo {
	return &OrganizationTemplateRepo{
		BaseRepo: gormx.NewBaseRepo[*model.OrganizationTemplate](db),
		m:        new(model.OrganizationTemplate),
	}
}

// CascadeStatus updates the status of all organization templates derived from the given base template
func (repo *OrganizationTemplateRepo) CascadeStatus(ctx context.Context, tx *gorm.DB, baseTemplateID int64, status constants.CommonStatus) (err error) {
	err = tx.WithContext(ctx).
		Model(&model.OrganizationTemplate{}).
		Where("base_template_id = ?", baseTemplateID).
		Update("status", status).Error
	return errorx.WithStack(err)
}

// CascadeDelete performs a hard delete on all organization templates derived from the given base template
func (repo *OrganizationTemplateRepo) CascadeDelete(ctx context.Context, tx *gorm.DB, baseTemplateID int64) (err error) {
	err = tx.WithContext(ctx).
		Where("base_template_id = ?", baseTemplateID).
		Delete(&model.OrganizationTemplate{}).Error
	return errorx.WithStack(err)
}
