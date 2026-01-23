package logic

import (
	"context"

	"gorm.io/gorm"

	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/bo"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/repo"
)

// BaseTemplateLogic provides base template business logic
type BaseTemplateLogic struct {
	repos                      *repo.AllRepos
	baseTemplateRepo           *repo.BaseTemplateRepo
	organizationTemplateRepo   *repo.OrganizationTemplateRepo
}

// NewBaseTemplateLogic creates a new base template logic instance
func NewBaseTemplateLogic(repos *repo.AllRepos) *BaseTemplateLogic {
	return &BaseTemplateLogic{
		repos:                    repos,
		baseTemplateRepo:         repos.BaseTemplateRepo,
		organizationTemplateRepo: repos.OrganizationTemplateRepo,
	}
}

// Create creates a new base template
func (l *BaseTemplateLogic) Create(ctx context.Context, req dto.BaseTemplateCreateDto) (result *bo.BaseTemplateBo, err error) {
	// Convert DTO to model
	data, err := copierx.Trans[dto.BaseTemplateCreateDto, model.BaseTemplate](&req)
	if err != nil {
		return nil, err
	}

	created, err := l.baseTemplateRepo.Create(ctx, data)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	result, err = copierx.Trans[model.BaseTemplate, bo.BaseTemplateBo](created)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Detail retrieves base template details by ID
func (l *BaseTemplateLogic) Detail(ctx context.Context, id int64) (result *bo.BaseTemplateBo, err error) {
	template, err := l.baseTemplateRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	result, err = copierx.Trans[model.BaseTemplate, bo.BaseTemplateBo](template)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Update updates an existing base template
func (l *BaseTemplateLogic) Update(ctx context.Context, req dto.BaseTemplateUpdateDto) (result *bo.BaseTemplateBo, err error) {
	// Get existing template
	existing, err := l.baseTemplateRepo.GetByID(ctx, req.Id)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	// Convert DTO to model for update
	updateData, err := copierx.Trans[dto.BaseTemplateUpdateDto, model.BaseTemplate](&req)
	if err != nil {
		return nil, err
	}

	// Update fields
	existing.Name = updateData.Name
	existing.Content = updateData.Content
	existing.Description = updateData.Description
	existing.VariableGroupSchema = updateData.VariableGroupSchema

	updated, err := l.baseTemplateRepo.Update(ctx, existing)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	result, err = copierx.Trans[model.BaseTemplate, bo.BaseTemplateBo](updated)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Page retrieves paginated base template list
func (l *BaseTemplateLogic) Page(ctx context.Context, pageReq dto.BaseTemplatePageRequestDto) (page bo.PageResultBo[bo.BaseTemplatePageBo], err error) {
	results, err := l.baseTemplateRepo.Page(ctx, pageReq)
	if err != nil {
		return page, err
	}

	// Convert model page results to BO page results
	page.Total = results.Total
	page.Page = results.Page
	page.Limit = results.Limit

	// Convert each item
	for _, item := range results.List {
		itemBo, err := copierx.Trans[model.BaseTemplate, bo.BaseTemplatePageBo](item)
		if err != nil {
			return page, err
		}
		page.List = append(page.List, *itemBo)
	}

	return page, nil
}

// SwitchStatus switches the status of a base template
func (l *BaseTemplateLogic) SwitchStatus(ctx context.Context, id int64, status constants.CommonStatus) (err error) {
	// Use transaction to cascade status changes
	err = l.baseTemplateRepo.Transaction(ctx, func(tx *gorm.DB) error {
		// Update base template status
		if err = l.baseTemplateRepo.UpdateStatusWithTx(ctx, tx, id, status); err != nil {
			return err
		}
		// Cascade to organization templates
		return l.organizationTemplateRepo.CascadeStatus(ctx, tx, id, status)
	})
	if err != nil {
		return errorx.WithStack(err)
	}
	return nil
}

// Delete deletes a base template by ID (hard delete with cascade)
func (l *BaseTemplateLogic) Delete(ctx context.Context, id int64) (err error) {
	// Use transaction to cascade delete
	err = l.baseTemplateRepo.Transaction(ctx, func(tx *gorm.DB) error {
		// Cascade delete to organization templates first
		if err = l.organizationTemplateRepo.CascadeDelete(ctx, tx, id); err != nil {
			return err
		}
		// Hard delete base template
		return l.baseTemplateRepo.DeleteWithTx(ctx, tx, id)
	})
	if err != nil {
		return errorx.WithStack(err)
	}
	return nil
}

// ExistsByTag checks if a template tag exists within the same type and category
func (l *BaseTemplateLogic) ExistsByTag(ctx context.Context, tag string, templateType constants.UserType, category constants.TemplateCategory, excludeID int64) (exists bool, err error) {
	return l.baseTemplateRepo.ExistsByTag(ctx, tag, templateType, category, excludeID)
}

// GetByIDForUpdate retrieves a base template model by ID for internal update operations
// This returns model for use in RPC layer where model modifications are needed
func (l *BaseTemplateLogic) GetByIDForUpdate(ctx context.Context, id int64) (*model.BaseTemplate, error) {
	return l.baseTemplateRepo.GetByID(ctx, id)
}
