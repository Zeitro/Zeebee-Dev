package logic

import (
	"bytes"
	"context"
	"encoding/json"
	"html/template"
	"strings"

	"github.com/duke-git/lancet/v2/slice"
	"gorm.io/gorm"

	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/templater/convertor"
)

// BaseTemplateLogic provides base template business logic
type BaseTemplateLogic struct {
	*BaseLogic
}

// NewBaseTemplateLogic creates a new base template logic instance
func NewBaseTemplateLogic(baseLogic *BaseLogic) *BaseTemplateLogic {
	return &BaseTemplateLogic{
		BaseLogic: baseLogic,
	}
}

// Create creates a new base template
func (l *BaseTemplateLogic) Create(ctx context.Context, req dto.BaseTemplateCreateDto) (result *dto.BaseTemplateEntityDto, err error) {
	// check tag uniqueness within the same type
	exists, err := l.ExistsByPurpose(ctx, req.Purpose, req.Type)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errorx.New("Template purpose already exists within the same type")
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

		if innerErr != nil {
			return errorx.WithStack(innerErr)
		}
		// Sync to organization templates
		if req.Scope == constants.BaseTemplateScopeOrganization {
			innerErr = l.syncOrganizationTemplates(ctx, tx, created, created.Type, false)
		}

		return innerErr
	})
	if err != nil {
		return nil, errorx.Wrap(err, "fail to create template")
	}

	result, err = copierx.Trans[model.BaseTemplate, dto.BaseTemplateEntityDto](created)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Detail retrieves base template details by ID
func (l *BaseTemplateLogic) Detail(ctx context.Context, id int64) (result *dto.BaseTemplateEntityDto, err error) {
	template, err := l.repos.BaseTemplateRepo.Detail(ctx, id, false)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	result, err = copierx.Trans[model.BaseTemplate, dto.BaseTemplateEntityDto](template)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// DetailByPurpose retrieves base template's content by purpose and type
func (l *BaseTemplateLogic) DetailByPurpose(ctx context.Context, purpose string, tplType constants.UserType) (result *dto.BaseTemplateEntityDto, err error) {
	template, err := l.repos.BaseTemplateRepo.FindOneByColumns(ctx, map[string]interface{}{
		"purpose": purpose,
		"type":    tplType,
	})
	if err != nil {
		return nil, err
	}
	result, err = copierx.Trans[model.BaseTemplate, dto.BaseTemplateEntityDto](template)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Update updates an existing base template
func (l *BaseTemplateLogic) Update(ctx context.Context, req dto.BaseTemplateUpdateDto) (result *dto.BaseTemplateEntityDto, err error) {
	// Get existing template
	existing, err := l.repos.BaseTemplateRepo.Detail(ctx, req.Id, false)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	oldScope := existing.Scope

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
	existing.Scope = updateData.Scope
	existing.IsConfigurable = updateData.IsConfigurable
	var updated *model.BaseTemplate
	err = l.repos.BaseTemplateRepo.Transaction(ctx, func(tx *gorm.DB) error {
		var innerErr error
		updated, innerErr = l.repos.BaseTemplateRepo.UpdateOneByIdWithTx(ctx, tx, existing.Id, existing)
		if innerErr != nil {
			return innerErr
		}
		if oldScope == req.Scope {
			return innerErr
		}
		switch req.Scope {
		case constants.BaseTemplateScopeOrganization:
			innerErr = l.syncOrganizationTemplates(ctx, tx, updated, updated.Type, true)
		case constants.BaseTemplateScopeSystem:
			innerErr = l.repos.OrganizationTemplateRepo.CascadeStatus(ctx, tx, updated.Id, constants.CommonStatusActive)
		default:
			innerErr = errorx.Errorf("Invalid scope:%v", req.Scope)
		}

		return innerErr
	})
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	result, err = copierx.Trans[model.BaseTemplate, dto.BaseTemplateEntityDto](updated)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Page retrieves paginated base template list
func (l *BaseTemplateLogic) Page(ctx context.Context, pageReq dto.BaseTemplatePageRequestDto) (page dto.PageResultDto[dto.BaseTemplatePageEntityDto], err error) {
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
		itemBo, err := copierx.Trans[model.BaseTemplate, dto.BaseTemplatePageEntityDto](item)
		if err != nil {
			return page, err
		}
		page.List = append(page.List, *itemBo)
	}

	return page, nil
}

// SwitchStatus switches the status of a base template
func (l *BaseTemplateLogic) SwitchStatus(ctx context.Context, id int64, status constants.CommonStatus) (err error) {
	tpl, err := l.repos.BaseTemplateRepo.FindOneById(ctx, id)
	if err != nil {
		return err
	}
	// Use transaction to cascade status changes
	err = l.repos.BaseTemplateRepo.Transaction(ctx, func(tx *gorm.DB) error {
		var innerErr error
		// Update base template status
		if innerErr = l.repos.BaseTemplateRepo.UpdateStatusWithTx(ctx, tx, id, status); innerErr != nil {
			return innerErr
		}
		// Cascade to organization templates
		innerErr = l.repos.OrganizationTemplateRepo.CascadeStatus(ctx, tx, id, status)
		if innerErr != nil {
			return innerErr
		}
		// Create missing organization templates
		if status == constants.CommonStatusActive {
			var allOrganizationIds []int64
			var tplOrganizationIds []int64
			// Get all active organization IDs
			allOrganizationIds, innerErr = l.repos.OrganizationRepo.GetActiveIdsWithTx(ctx, tx, tpl.Type)
			if innerErr != nil {
				return innerErr
			}
			// Get organization IDs of the template
			tplOrganizationIds, innerErr = l.repos.OrganizationTemplateRepo.GetOrganizationIdsWithTx(ctx, tx, id)
			if innerErr != nil {
				return innerErr
			}
			diffOrganizationIds := slice.Difference(allOrganizationIds, tplOrganizationIds)
			var organizationTemplates []*model.OrganizationTemplate
			for _, organizationId := range diffOrganizationIds {
				organizationTemplate, innerErr := copierx.Trans[model.BaseTemplate, model.OrganizationTemplate](tpl)
				if innerErr != nil {
					return innerErr
				}
				organizationTemplate.Id = 0
				organizationTemplate.BaseTemplateId = tpl.Id
				organizationTemplate.OrganizationId = organizationId
				organizationTemplate.Status = constants.CommonStatusActive
				organizationTemplate.HasConfigured = !tpl.IsConfigurable
				organizationTemplates = append(organizationTemplates, organizationTemplate)
			}
			// Create organization templates
			organizationTemplateLogic := NewOrganizationTemplateLogic(l.BaseLogic)
			innerErr = organizationTemplateLogic.BatchCreateWithTx(ctx, tx, organizationTemplates, false)
		}
		return innerErr
	})
	return errorx.WithStack(err)
}

// Delete deletes a base template by ID (hard delete with cascade)
func (l *BaseTemplateLogic) Delete(ctx context.Context, id int64) (err error) {
	// Use transaction to cascade delete
	err = l.repos.BaseTemplateRepo.Transaction(ctx, func(tx *gorm.DB) error {
		var innerErr error
		// Cascade delete to organization templates first
		if innerErr = l.repos.OrganizationTemplateRepo.CascadeDelete(ctx, tx, id); innerErr != nil {
			return innerErr
		}
		// Hard delete base template
		return l.repos.BaseTemplateRepo.DeleteWithTx(ctx, tx, id)
	})
	if err != nil {
		return errorx.WithStack(err)
	}
	return nil
}

// ExistsByPurpose checks if a template tag exists within the same type and category
func (l *BaseTemplateLogic) ExistsByPurpose(ctx context.Context, purpose string, templateType constants.UserType) (exists bool, err error) {
	return l.repos.BaseTemplateRepo.ExistsByPurpose(ctx, purpose, templateType)
}

// Preview retrieves a preview of a base template
func (l *BaseTemplateLogic) Preview(_ context.Context, req dto.BaseTemplatePreviewDto) (result string, err error) {
	tpl, err := template.New(req.Name).
		Funcs(convertor.TplFunc()).
		Parse(req.Content)
	if err != nil {
		return "", errorx.Wrap(err, "Failed to parse template")
	}
	data := map[string]interface{}{}
	if len(req.ExampleVariables) > 0 {
		var exampleData map[string]interface{}
		err = json.Unmarshal([]byte(req.ExampleVariables), &exampleData)
		if err != nil {
			return
		}
		for k, v := range exampleData {
			key := strings.Trim(k, ".")
			data[key] = v
		}
	}
	for _, groupDto := range req.VariableGroupSchema {
		for _, variableDto := range groupDto.Variables {
			if variableDto.Source == constants.TemplateVariableSourceFixed || variableDto.Source == constants.TemplateVariableSourceTemplate {
				key := strings.Trim(variableDto.Key, ".")
				data[key] = variableDto.FixedValue
			}
		}
	}

	buf := bytes.NewBuffer(nil)
	err = tpl.Execute(buf, data)
	if err != nil {
		return "", errorx.Wrap(err, "Failed to render template")
	}
	result = buf.String()

	return result, nil
}

func (l *BaseTemplateLogic) syncOrganizationTemplates(ctx context.Context, tx *gorm.DB, baseTemplate *model.BaseTemplate, organizationType constants.UserType, checkExisted bool) (err error) {
	if baseTemplate == nil {
		err = errorx.New("Base template must exist")
		return
	}
	if baseTemplate.Scope != constants.BaseTemplateScopeOrganization {
		err = errorx.New("Template scope must be organization")
		return
	}
	// Sync to organization templates
	organizationIds, err := l.repos.OrganizationRepo.GetActiveIdsWithTx(ctx, tx, organizationType)
	if err != nil {
		return err
	}
	var organizationTemplates []*model.OrganizationTemplate
	for _, organizationId := range organizationIds {
		organizationTemplate, err := copierx.Trans[model.BaseTemplate, model.OrganizationTemplate](baseTemplate)
		if err != nil {
			return err
		}
		organizationTemplate.Id = 0
		organizationTemplate.BaseTemplateId = baseTemplate.Id
		organizationTemplate.OrganizationId = organizationId
		organizationTemplate.Status = constants.CommonStatusActive
		organizationTemplate.HasConfigured = !baseTemplate.IsConfigurable
		organizationTemplates = append(organizationTemplates, organizationTemplate)
	}
	// Create organization templates
	organizationTemplateLogic := NewOrganizationTemplateLogic(l.BaseLogic)
	err = organizationTemplateLogic.BatchCreateWithTx(ctx, tx, organizationTemplates, checkExisted)
	return err
}
