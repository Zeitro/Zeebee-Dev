package logic

import (
	"bytes"
	"context"
	"encoding/json"
	"html/template"

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
	repos *repo.AllRepos
}

// NewBaseTemplateLogic creates a new base template logic instance
func NewBaseTemplateLogic(repos *repo.AllRepos) *BaseTemplateLogic {
	return &BaseTemplateLogic{
		repos: repos,
	}
}

// Create creates a new base template
func (l *BaseTemplateLogic) Create(ctx context.Context, req dto.BaseTemplateCreateDto) (result *bo.BaseTemplateBo, err error) {
	// check tag uniqueness within the same type
	exists, err := l.ExistsByPurpose(ctx, req.Purpose, req.Type)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errorx.New("Template Tag Already Exists Within The Same Type")
	}

	// Convert DTO to model
	data, err := copierx.Trans[dto.BaseTemplateCreateDto, model.BaseTemplate](&req)
	if err != nil {
		return nil, err
	}
	data.Status = constants.CommonStatusActive
	var created *model.BaseTemplate
	err = l.repos.BaseTemplateRepo.Transaction(ctx, func(tx *gorm.DB) error {
		var innerErr error
		// Create base template
		created, innerErr = l.repos.BaseTemplateRepo.CreateWithTx(ctx, tx, data, false)
		if err != nil {
			return errorx.WithStack(err)
		}
		// Sync to organization templates
		organizationIds, innerErr := l.repos.OrganizationRepo.GetIdsWithTx(ctx, tx, req.Type)
		if innerErr != nil {
			return innerErr
		}
		var organizationTemplates []*model.OrganizationTemplate
		for _, organizationId := range organizationIds {
			organizationTemplate, innerErr := copierx.Trans[dto.BaseTemplateCreateDto, model.OrganizationTemplate](&req)
			if innerErr != nil {
				return innerErr
			}
			organizationTemplate.BaseTemplateId = created.Id
			organizationTemplate.OrganizationId = organizationId
			organizationTemplate.Status = constants.CommonStatusActive
			organizationTemplates = append(organizationTemplates, organizationTemplate)
		}
		_, innerErr = l.repos.OrganizationTemplateRepo.BatchCreateWithTx(ctx, tx, organizationTemplates, false)

		return innerErr
	})
	if err != nil {
		return nil, errorx.Wrap(err, "Failed To Create Template")
	}

	result, err = copierx.Trans[model.BaseTemplate, bo.BaseTemplateBo](created)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Detail retrieves base template details by ID
func (l *BaseTemplateLogic) Detail(ctx context.Context, id int64) (result *bo.BaseTemplateBo, err error) {
	template, err := l.repos.BaseTemplateRepo.GetByID(ctx, id)
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
	existing, err := l.repos.BaseTemplateRepo.GetByID(ctx, req.Id)
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
	existing.Category = updateData.Category
	existing.ExampleVariables = updateData.ExampleVariables
	existing.Description = updateData.Description
	existing.VariableGroupSchema = updateData.VariableGroupSchema

	updated, err := l.repos.BaseTemplateRepo.Update(ctx, existing)
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
	results, err := l.repos.BaseTemplateRepo.Page(ctx, pageReq)
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
	err = l.repos.BaseTemplateRepo.Transaction(ctx, func(tx *gorm.DB) error {
		// Update base template status
		if err = l.repos.BaseTemplateRepo.UpdateStatusWithTx(ctx, tx, id, status); err != nil {
			return err
		}
		// Cascade to organization templates
		return l.repos.OrganizationTemplateRepo.CascadeStatus(ctx, tx, id, status)
	})
	if err != nil {
		return errorx.WithStack(err)
	}
	return nil
}

// Delete deletes a base template by ID (hard delete with cascade)
func (l *BaseTemplateLogic) Delete(ctx context.Context, id int64) (err error) {
	// Use transaction to cascade delete
	err = l.repos.BaseTemplateRepo.Transaction(ctx, func(tx *gorm.DB) error {
		// Cascade delete to organization templates first
		if err = l.repos.OrganizationTemplateRepo.CascadeDelete(ctx, tx, id); err != nil {
			return err
		}
		// Hard delete base template
		return l.repos.BaseTemplateRepo.DeleteWithTx(ctx, tx, id)
	})
	if err != nil {
		return errorx.WithStack(err)
	}
	return nil
}

// ExistsByTag checks if a template tag exists within the same type and category
func (l *BaseTemplateLogic) ExistsByPurpose(ctx context.Context, purpose string, templateType constants.UserType) (exists bool, err error) {
	return l.repos.BaseTemplateRepo.ExistsByPurpose(ctx, purpose, templateType)
}

// Preview retrieves a preview of a base template
func (l *BaseTemplateLogic) Preview(ctx context.Context, req dto.BaseTemplatePreviewDto) (result string, err error) {
	tpl, err := template.New(req.Name).Parse(req.Content)
	if err != nil {
		return "", errorx.Wrap(err, "Failed To Parse Template")
	}
	data := map[string]interface{}{}
	if len(req.ExampleVariables) > 0 {
		var exampleData map[string]interface{}
		err = json.Unmarshal([]byte(req.ExampleVariables), &exampleData)
		if err != nil {
			return
		}
		for k, v := range exampleData {
			data[k] = v
		}
	}
	for _, groupDto := range req.VariableGroupSchema {

		for _, variableDto := range groupDto.Variables {
			if variableDto.Source == constants.TemplateVariableSourceFixed || variableDto.Source == constants.TemplateVariableSourceTemplate {
				data[variableDto.Key] = variableDto.FixedValue
			}
		}
	}

	buf := bytes.NewBuffer(nil)
	err = tpl.Execute(buf, data)
	if err != nil {
		return "", errorx.Wrap(err, "Failed To Render Template")
	}
	result = buf.String()

	return result, nil
}
