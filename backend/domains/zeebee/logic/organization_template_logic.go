package logic

import (
	"bytes"
	"context"
	"text/template"
	"time"

	"github.com/duke-git/lancet/v2/slice"
	"gorm.io/gorm"

	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
)

// OrganizationTemplateLogic provides organization template business logic
type OrganizationTemplateLogic struct {
	*BaseLogic
}

// NewOrganizationTemplateLogic creates a new organization template logic instance
func NewOrganizationTemplateLogic(baseLogic *BaseLogic) *OrganizationTemplateLogic {
	return &OrganizationTemplateLogic{
		BaseLogic: baseLogic,
	}
}

// Detail retrieves detailed information about an organization template by its ID.
// It fetches the template from the repository and converts it to a business object.
func (l *OrganizationTemplateLogic) Detail(ctx context.Context, id int64) (result *dto.OrganizationTemplateEntityDto, err error) {
	organizationTemplate, err := l.repos.OrganizationTemplateRepo.Detail(ctx, id, true)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	organizationTemplate, err = l.renderOrganizationSystemInfo(organizationTemplate)
	if err != nil {
		return nil, err
	}
	result, err = copierx.Trans[model.OrganizationTemplate, dto.OrganizationTemplateEntityDto](organizationTemplate)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// DetailByPurpose retrieves an organization template by its purpose and organization ID.
func (l *OrganizationTemplateLogic) DetailByPurpose(ctx context.Context, organizationId int64, purpose string) (result *dto.OrganizationTemplateEntityDto, err error) {
	organizationTemplate, err := l.repos.OrganizationTemplateRepo.DetailByPurpose(ctx, organizationId, purpose, true)
	if err != nil {
		return nil, err
	}

	organizationTemplate, err = l.renderOrganizationSystemInfo(organizationTemplate)
	if err != nil {
		return nil, err
	}

	result, err = copierx.Trans[model.OrganizationTemplate, dto.OrganizationTemplateEntityDto](organizationTemplate)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (l *OrganizationTemplateLogic) UpdateVariableGroup(ctx context.Context, req dto.OrganizationTemplateUpdateVariableGroupDto) (result *dto.OrganizationTemplateEntityDto, err error) {
	// Get existing template
	existing, err := l.repos.OrganizationTemplateRepo.Detail(ctx, req.Id, false)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	if existing.OrganizationId != req.OrganizationId {
		return nil, errorx.New("Template does not belong to organization")
	}

	oldGroups := existing.VariableGroupSchema
	groups := req.VariableGroupSchema
	if len(groups) != len(oldGroups) {
		return nil, errorx.Errorf("Variable group number not match: expected %d, got %d", len(oldGroups), len(groups))
	}
	// Check variable group
	var checkVariableFn func(oldVariable *data.TemplateVariable, newVariable *data.TemplateVariable) (err error)
	checkVariableFn = func(oldVariable *data.TemplateVariable, newVariable *data.TemplateVariable) (err error) {
		if oldVariable == nil {
			return nil
		}
		if newVariable == nil {
			return errorx.Errorf("Variable %s not Found", oldVariable.Label)
		}

		if oldVariable.Source == constants.TemplateVariableSourceTemplate && newVariable.Source != constants.TemplateVariableSourceTemplate {
			return errorx.Errorf("%s source type template can not be changed To %v", oldVariable.Label, newVariable.Source)
		}

		if oldVariable.Type == constants.TemplateVariableTypeObject {
			if len(oldVariable.Properties) != len(newVariable.Properties) {
				return errorx.Errorf("Variable %s Property Number Not Match", oldVariable.Label)
			}

			for key, oldProp := range oldVariable.Properties {
				newProp, ok := newVariable.Properties[key]
				if !ok || newProp == nil {
					return errorx.Errorf("Variable Property %s Not Found", oldProp.Key)
				}
				if oldProp.Key != newProp.Key ||
					oldProp.Label != newProp.Label ||
					oldProp.IsRequired != newProp.IsRequired ||
					oldProp.IsEditable != newProp.IsEditable ||
					oldProp.Type != newProp.Type {
					return errorx.Errorf("%s Property Type Not Match", oldProp.Label)
				}

				err = checkVariableFn(oldProp, newProp)
				if err != nil {
					return err
				}
			}
		}
		if oldVariable.Type == constants.TemplateVariableTypeArray {
			if len(oldVariable.Items) != len(newVariable.Items) {
				return errorx.Errorf("Variable %s Item Number Not Match", oldVariable.Label)
			}
			for _, oldItem := range oldVariable.Items {
				newItem, ok := slice.FindBy(newVariable.Items, func(index int, item *data.TemplateVariable) bool {
					return oldItem.Key == item.Key &&
						oldItem.Label == item.Label &&
						oldItem.IsRequired == item.IsRequired &&
						oldItem.IsEditable == item.IsEditable &&
						oldItem.Type == item.Type
				})
				if !ok {
					return errorx.Errorf("Variable Item %s Not Found", oldItem.Key)
				}
				err = checkVariableFn(oldItem, newItem)
				if err != nil {
					return err
				}
			}
		}

		return newVariable.ValidateRequiredField()
	}
	for _, oldGroup := range oldGroups {
		newGroup, ok := slice.FindBy(groups, func(index int, item *data.TemplateVariableGroup) bool {
			return oldGroup.Name == item.Name
		})
		if !ok {
			return nil, errorx.Errorf("Variable Group %s Not Found", oldGroup.Name)
		}
		if len(newGroup.Variables) != len(oldGroup.Variables) {
			return nil, errorx.Errorf("Variable Group %s Number Not Match", oldGroup.Name)
		}

		for _, oldVariable := range oldGroup.Variables {
			newVariable, ok := slice.FindBy(newGroup.Variables, func(index int, item *data.TemplateVariable) bool {
				return oldVariable.Key == item.Key &&
					oldVariable.Label == item.Label &&
					oldVariable.IsRequired == item.IsRequired &&
					oldVariable.IsEditable == item.IsEditable &&
					oldVariable.Type == item.Type
			})
			if !ok {
				return nil, errorx.Errorf("Variable %s Not Matched", oldVariable.Label)
			}
			err = checkVariableFn(oldVariable, newVariable)
			if err != nil {
				return nil, err
			}
		}
	}

	existing.VariableGroupSchema = groups
	existing.HasConfigured = true
	updated, err := l.repos.OrganizationTemplateRepo.UpdateOneById(ctx, req.Id, existing)
	if err != nil {
		return nil, err
	}

	organization, err := l.repos.OrganizationRepo.FindOneById(ctx, updated.OrganizationId)
	if err != nil {
		return nil, err
	}
	updated.Organization = organization
	updated, err = l.renderOrganizationSystemInfo(updated)
	if err != nil {
		return nil, err
	}

	result, err = copierx.Trans[model.OrganizationTemplate, dto.OrganizationTemplateEntityDto](updated)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (l *OrganizationTemplateLogic) Page(ctx context.Context, pageReq dto.OrganizationTemplatePageRequestDto) (page dto.PageResultDto[dto.OrganizationTemplatePageEntityDto], err error) {
	pageResult, err := l.repos.OrganizationTemplateRepo.Page(ctx, pageReq)
	if err != nil {
		return page, err
	}

	pagePtr, err := copierx.Trans[gormx.PageResponse[*model.OrganizationTemplate], dto.PageResultDto[dto.OrganizationTemplatePageEntityDto]](&pageResult)
	if err != nil {
		return page, err
	}
	return *pagePtr, nil
}

// Reset resets an organization template to its base template settings.
func (l *OrganizationTemplateLogic) Reset(ctx context.Context, id int64, organizationId int64) (result *dto.OrganizationTemplateEntityDto, err error) {
	// Retrieve the template details with full information
	template, err := l.repos.OrganizationTemplateRepo.Detail(ctx, id, true)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	// Verify that the template belongs to the specified organization
	if template.OrganizationId != organizationId {
		return nil, errorx.New("Template does not belong to organization")
	}
	// Get the base template reference
	baseTemplate := template.BaseTemplate
	if baseTemplate == nil {
		return nil, errorx.New("Template has no base template")
	}
	// Reset all template properties to match the base template
	template.Content = baseTemplate.Content
	template.VariableGroupSchema = l.mergeTemplateVariableGroup(template.VariableGroupSchema, baseTemplate.VariableGroupSchema)
	template.ExampleVariables = baseTemplate.ExampleVariables
	template.Name = baseTemplate.Name
	template.Purpose = baseTemplate.Purpose
	template.Description = baseTemplate.Description
	template.Category = baseTemplate.Category
	template.HasConfigured = !baseTemplate.IsConfigurable

	// Update the template in the repository
	updated, err := l.repos.OrganizationTemplateRepo.UpdateOneById(ctx, id, template)
	if err != nil {
		return nil, err
	}

	organization, err := l.repos.OrganizationRepo.FindOneById(ctx, template.OrganizationId)
	if err != nil {
		return nil, err
	}
	updated.Organization = organization

	updated, err = l.renderOrganizationSystemInfo(updated)
	if err != nil {
		return nil, err
	}

	// Convert the updated model to the business object format
	result, err = copierx.Trans[model.OrganizationTemplate, dto.OrganizationTemplateEntityDto](updated)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// ResetByBaseTemplateId
func (l *OrganizationTemplateLogic) ResetByBaseTemplateId(ctx context.Context, baseTemplateId int64) (err error) {
	baseTemplate, err := l.repos.BaseTemplateRepo.Detail(ctx, baseTemplateId, true)
	if err != nil {
		return err
	}
	err = l.repos.OrganizationTemplateRepo.Transaction(ctx, func(tx *gorm.DB) error {
		var innerErr error
		for _, organizationTemplate := range baseTemplate.OrganizationTemplates {
			organizationTemplate := organizationTemplate

			organizationTemplate.Content = baseTemplate.Content
			organizationTemplate.VariableGroupSchema = l.mergeTemplateVariableGroup(organizationTemplate.VariableGroupSchema, baseTemplate.VariableGroupSchema)
			organizationTemplate.ExampleVariables = baseTemplate.ExampleVariables
			organizationTemplate.Name = baseTemplate.Name
			organizationTemplate.Purpose = baseTemplate.Purpose
			organizationTemplate.Description = baseTemplate.Description
			organizationTemplate.Category = baseTemplate.Category
			organizationTemplate.HasConfigured = !baseTemplate.IsConfigurable

			_, innerErr = l.repos.OrganizationTemplateRepo.UpdateOneByIdWithTx(ctx, tx, organizationTemplate.Id, organizationTemplate)
			if innerErr != nil {
				return innerErr
			}
		}

		return innerErr
	})

	return errorx.WithStack(err)
}

// BatchCreateWithTx creates multiple organization templates in a single database transaction.
// It can optionally check for existing templates before creation to avoid duplicates.
func (l *OrganizationTemplateLogic) BatchCreateWithTx(ctx context.Context, tx *gorm.DB, templates []*model.OrganizationTemplate, checkExist bool) (err error) {
	if checkExist {
		var tmpTemplates []*model.OrganizationTemplate
		for index, tpl := range templates {
			exist, err := l.repos.OrganizationTemplateRepo.IsExistedByColumnsWithTx(ctx, tx, map[string]interface{}{
				"organization_id":  tpl.OrganizationId,
				"base_template_id": tpl.BaseTemplateId,
			})
			if err != nil {
				return err
			}
			if !exist {
				tmpTemplates = append(tmpTemplates, templates[index])
			}
		}
		templates = tmpTemplates
	}
	if len(templates) == 0 {
		return nil
	}
	_, err = l.repos.OrganizationTemplateRepo.BatchCreateWithTx(ctx, tx, templates, false)

	return err
}

// SyncByOrganizationIdWithTx synchronizes organization templates based on the organization ID within a transaction.
// It retrieves the organization details and all active base templates matching the organization type,
// then creates corresponding organization templates.
func (l *OrganizationTemplateLogic) SyncByOrganizationIdWithTx(ctx context.Context, tx *gorm.DB, organization *model.Organization) (err error) {
	if organization == nil {
		return errorx.New("Organization is empty")
	}
	// Find the organization by its ID within the transaction
	// organization, err := l.repos.OrganizationRepo.FindOneByIdWithTx(ctx, tx, organizationId)
	// if err != nil {
	// 	return err
	// }
	// Get all base templates that are active and match the organization type
	baseTemplates, err := l.repos.BaseTemplateRepo.FindAllWithTx(ctx, tx, map[string]interface{}{
		"status": constants.CommonStatusActive,
		"type":   organization.Type,
		"scope":  constants.BaseTemplateScopeOrganization,
	})
	if err != nil {
		return err
	}
	// Prepare a slice of organization templates based on base templates
	templates := make([]*model.OrganizationTemplate, 0, len(baseTemplates))
	for _, baseTemplate := range baseTemplates {
		// Create a local copy of the base template to avoid reference issues
		baseTemplate := baseTemplate
		// Copy base template properties to organization template
		tpl, err := copierx.Trans[model.BaseTemplate, model.OrganizationTemplate](baseTemplate)
		if err != nil {
			return err
		}
		// Reset the ID and set organization-specific fields
		tpl.Id = 0
		tpl.OrganizationId = organization.Id
		tpl.BaseTemplateId = baseTemplate.Id
		tpl.HasConfigured = !baseTemplate.IsConfigurable
		templates = append(templates, tpl)
	}

	// Batch create all organization templates within the transaction
	return l.BatchCreateWithTx(ctx, tx, templates, true)
}

func (l *OrganizationTemplateLogic) renderOrganizationSystemInfo(organizationTemplate *model.OrganizationTemplate) (result *model.OrganizationTemplate, err error) {
	if organizationTemplate == nil {
		return nil, errorx.New("Organization template not found")
	}
	organization := organizationTemplate.Organization
	if organization != nil {
		tpl, err := template.New(organizationTemplate.Name).
			Delims("[[", "]]").
			Parse(organizationTemplate.Content)
		if err != nil {
			return nil, errorx.Wrap(err, "Failed to parse template")
		}
		now := time.Now()
		renderData := map[string]interface{}{
			"system": map[string]interface{}{
				"organization": organization,
				"now_format1":  now.Format("02/01/2006"),
				"now_format2":  now.Month().String() + now.Format(" 02, 2006"),
			},
		}
		var buf bytes.Buffer
		err = tpl.Execute(&buf, renderData)
		if err != nil {
			return nil, errorx.Wrap(err, "Failed to render template")
		}
		organizationTemplate.Content = buf.String()
	}

	return organizationTemplate, nil
}

// mergeTemplateVariables merges template variables from different sources
func (l *OrganizationTemplateLogic) mergeTemplateVariableGroup(oldVariableGroup, newVariableGroup []*data.TemplateVariableGroup) (result []*data.TemplateVariableGroup) {
	if len(oldVariableGroup) == 0 {
		return newVariableGroup
	}
	result = make([]*data.TemplateVariableGroup, 0, len(newVariableGroup))
	for _, newGroup := range newVariableGroup {
		newGroup := newGroup

		oldGroup, ok := slice.FindBy(oldVariableGroup, func(index int, item *data.TemplateVariableGroup) bool {
			return newGroup.Name == item.Name
		})
		if !ok || oldGroup == nil || oldGroup.Variables == nil {
			result = append(result, newGroup)
			continue
		}

		mergedVariables := make([]*data.TemplateVariable, 0, len(newGroup.Variables))
		for _, newVariable := range newGroup.Variables {
			newVariable := newVariable
			oldVar, ok := slice.FindBy(oldGroup.Variables, func(index int, item *data.TemplateVariable) bool {
				return newVariable.Key == item.Key &&
					newVariable.Type == item.Type &&
					newVariable.Format == item.Format
			})
			if !ok || oldVar == nil {
				mergedVariables = append(mergedVariables, newVariable)
				continue
			}
			newVariable.FixedValue = oldVar.FixedValue
			mergedVar := l.mergeVariable(newVariable, oldVar)
			mergedVariables = append(mergedVariables, mergedVar)
		}

		mergedGroup := &data.TemplateVariableGroup{
			Name:      newGroup.Name,
			Variables: mergedVariables,
		}
		result = append(result, mergedGroup)
	}

	return result
}

func (l *OrganizationTemplateLogic) mergeVariable(newVar *data.TemplateVariable, oldVar *data.TemplateVariable) (result *data.TemplateVariable) {
	if oldVar == nil || newVar == nil {
		return newVar
	}

	result, _ = copierx.Trans[data.TemplateVariable, data.TemplateVariable](newVar)
	if !newVar.IsEditable {
		return result
	}

	if (newVar.Source == constants.TemplateVariableSourceTemplate && oldVar.Source != constants.TemplateVariableSourceTemplate) ||
		(newVar.Source != constants.TemplateVariableSourceTemplate && oldVar.Source == constants.TemplateVariableSourceTemplate) {
		return result
	}
	result.FixedValue = oldVar.FixedValue

	if result.Type == constants.TemplateVariableTypeObject {
		if oldVar.Properties == nil {
			oldVar.Properties = make(map[string]*data.TemplateVariable)
		}

		for key := range result.Properties {
			result.Properties[key] = l.mergeVariable(result.Properties[key], oldVar.Properties[key])
		}
	}

	if result.Type == constants.TemplateVariableTypeArray {
		for key := range result.Items {
			subOldVar, ok := slice.FindBy(oldVar.Items, func(index int, item *data.TemplateVariable) bool {
				return result.Items[key].Key == item.Key &&
					result.Items[key].Type == item.Type &&
					result.Items[key].Format == item.Format
			})
			if ok {
				result.Items[key].FixedValue = subOldVar.FixedValue
				result.Items[key] = l.mergeVariable(result.Items[key], subOldVar)
			}
		}
	}

	return result
}
