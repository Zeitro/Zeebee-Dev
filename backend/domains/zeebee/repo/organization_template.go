package repo

import (
	"context"

	"gorm.io/gorm"

	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
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

// Detail retrieves an organization template by its ID.
// It optionally preloads the associated BaseTemplate if withPreload is true.
func (repo *OrganizationTemplateRepo) Detail(ctx context.Context, id int64, withPreload bool) (result *model.OrganizationTemplate, err error) {
	query := repo.Client(ctx).Where("id = ?", id)
	if withPreload {
		query = query.Preload("BaseTemplate")
	}
	err = query.First(&result).Error
	return result, errorx.WithStack(err)
}

func (repo *OrganizationTemplateRepo) Page(ctx context.Context, pageReq dto.OrganizationTemplatePageRequestDto) (page gormx.PageResponse[*model.OrganizationTemplate], err error) {
	filter := pageReq.Filter
	sorter := pageReq.Sorter
	baseTemplateTable := model.OrganizationTemplate{}.TableName()

	query := repo.Client(ctx).Model(&model.OrganizationTemplate{}).Where("organization_id = ?", filter.OrganizationId)

	// Apply filters
	if len(filter.Keyword) > 0 {
		keyword := "%" + filter.Keyword + "%"
		query = query.Where("name LIKE ? OR purpose LIKE ?", keyword, keyword)
	}
	if len(filter.Type) > 0 {
		query = query.Where("type = ?", filter.Type)
	}
	if len(filter.Category) > 0 {
		query = query.Where("category = ?", filter.Category)
	}
	if len(filter.Status) > 0 {
		query = query.Where("status = ?", filter.Status)
	}

	// Apply sorting
	if len(sorter.Column) > 0 {
		query = query.Order(sorter.Column + " " + sorter.Sort)
	}
	query = query.Order(baseTemplateTable + ".created_at DESC")

	return repo.BaseRepo.Page(query, gormx.PageRequest{
		Page:  pageReq.Page,
		Limit: pageReq.Limit,
	})
}
