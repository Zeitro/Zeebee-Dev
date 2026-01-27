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

type BaseTemplateRepo struct {
	*gormx.BaseRepo[*model.BaseTemplate]
	m *model.BaseTemplate
}

func NewBaseTemplateRepo(db *gorm.DB) *BaseTemplateRepo {
	return &BaseTemplateRepo{
		BaseRepo: gormx.NewBaseRepo[*model.BaseTemplate](db),
		m:        new(model.BaseTemplate),
	}
}

// Create creates a new base template
func (repo *BaseTemplateRepo) Create(ctx context.Context, data *model.BaseTemplate) (result *model.BaseTemplate, err error) {
	// Validate variable group schema before creating
	if err = data.BeforeValidate(); err != nil {
		return nil, errorx.WithStack(err)
	}
	created, err := repo.BaseRepo.Create(ctx, data)
	return created, errorx.WithStack(err)
}

// GetByID retrieves a base template by ID
func (repo *BaseTemplateRepo) GetByID(ctx context.Context, id int64) (result *model.BaseTemplate, err error) {
	return repo.BaseRepo.FindOneById(ctx, id)
}

// Update updates an existing base template
func (repo *BaseTemplateRepo) Update(ctx context.Context, data *model.BaseTemplate) (result *model.BaseTemplate, err error) {
	// Validate variable group schema before updating
	// if err = data.BeforeValidate(); err != nil {
	// 	return nil, errorx.WithStack(err)
	// }
	updated, err := repo.BaseRepo.UpdateOneById(ctx, data.Id, data)
	return updated, errorx.WithStack(err)
}

// UpdateStatus updates the status of a base template
func (repo *BaseTemplateRepo) UpdateStatus(ctx context.Context, id int64, status constants.CommonStatus) (err error) {
	return repo.BaseRepo.UpdateColumnsById(ctx, id, map[string]interface{}{
		"status": status,
	})
}

// UpdateStatusWithTx updates the status of a base template within a transaction
func (repo *BaseTemplateRepo) UpdateStatusWithTx(ctx context.Context, tx *gorm.DB, id int64, status constants.CommonStatus) (err error) {
	return repo.BaseRepo.UpdateColumnsByIdWithTx(ctx, tx, id, map[string]interface{}{
		"status": status,
	})
}

// Delete performs a hard delete on a base template
func (repo *BaseTemplateRepo) Delete(ctx context.Context, id int64) (err error) {
	err = repo.Client(ctx).Delete(&model.BaseTemplate{}, id).Error
	return errorx.WithStack(err)
}

// DeleteWithTx performs a hard delete on a base template within a transaction
func (repo *BaseTemplateRepo) DeleteWithTx(ctx context.Context, tx *gorm.DB, id int64) (err error) {
	err = tx.WithContext(ctx).Delete(&model.BaseTemplate{}, id).Error
	return errorx.WithStack(err)
}

// Page retrieves a paginated list of base templates
func (repo *BaseTemplateRepo) Page(ctx context.Context, pageReq dto.BaseTemplatePageRequestDto) (page gormx.PageResponse[*model.BaseTemplate], err error) {
	filter := pageReq.Filter
	sorter := pageReq.Sorter
	baseTemplateTable := model.BaseTemplate{}.TableName()

	query := repo.Client(ctx).Model(&model.BaseTemplate{})

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

// ExistsByPurpose checks if a template with the given tag exists within the same type and category
func (repo *BaseTemplateRepo) ExistsByPurpose(ctx context.Context, purpose string, templateType constants.UserType) (exists bool, err error) {
	return repo.BaseRepo.IsExistedByColumns(ctx, map[string]interface{}{
		"purpose": purpose,
		"type":    templateType,
	})
}
