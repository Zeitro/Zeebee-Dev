package logic

import (
	"context"

	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/bo"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/repo"
)

// OrganizationTemplateLogic provides base template business logic
type OrganizationTemplateLogic struct {
	repos *repo.AllRepos
}

// NewOrganizationTemplateLogic creates a new base template logic instance
func NewOrganizationTemplateLogic(repos *repo.AllRepos) *OrganizationTemplateLogic {
	return &OrganizationTemplateLogic{
		repos: repos,
	}
}

// Detail retrieves detailed information about an organization template by its ID.
// It fetches the template from the repository and converts it to a business object.
func (l *OrganizationTemplateLogic) Detail(ctx context.Context, id int64) (result *bo.OrganizationTemplateBo, err error) {
	template, err := l.repos.OrganizationTemplateRepo.Detail(ctx, id, true)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	result, err = copierx.Trans[model.OrganizationTemplate, bo.OrganizationTemplateBo](template)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// DetailByPurpose retrieves an organization template by its purpose and organization ID.
func (l *OrganizationTemplateLogic) DetailByPurpose(ctx context.Context, organizationId int64, purpose string) (result *bo.OrganizationTemplateBo, err error) {
	template, err := l.repos.OrganizationTemplateRepo.FindOneByColumns(ctx, map[string]interface{}{
		"organization_id": organizationId,
		"purpose":         purpose,
	})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	result, err = copierx.Trans[model.OrganizationTemplate, bo.OrganizationTemplateBo](template)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (l *OrganizationTemplateLogic) UpdateVariableGroup(ctx context.Context, req dto.OrganizationTemplateUpdateVariableGroupDto) (result *bo.OrganizationTemplateBo, err error) {
	// Get existing template
	existing, err := l.repos.OrganizationTemplateRepo.Detail(ctx, req.Id, false)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	if existing.OrganizationId != req.OrganizationId {
		return nil, errorx.New("Template Not Belong To Organization")
	}

	oldGroups := existing.VariableGroupSchema
	groupsPtr, err := copierx.Trans[[]*dto.TemplateVariableGroupDto, []*model.TemplateVariableGroup](&req.VariableGroupSchema)
	if err != nil {
		return nil, err
	}
	groups := *groupsPtr
	if len(groups) != len(oldGroups) {
		return nil, errorx.New("Variable Group Number Not Match")
	}
	// Check variable group
	for index, newGroup := range groups {
		if newGroup.Name != oldGroups[index].Name {
			return nil, errorx.Errorf("Variable Name Not Match %s", oldGroups[index].Name)
		}
		if len(newGroup.Variables) != len(oldGroups[index].Variables) {
			return nil, errorx.New("Variable Number Not Match")
		}
		for subIndex, newVariable := range newGroup.Variables {
			if oldGroups[index].Variables[subIndex].Source == constants.TemplateVariableSourceTemplate && newVariable.Source != constants.TemplateVariableSourceTemplate {
				return nil, errorx.Errorf("Source Type Template Can Not Be Changed To %v", newVariable.Source)
			}
		}
	}
	existing.VariableGroupSchema = groups
	updated, err := l.repos.OrganizationTemplateRepo.UpdateOneById(ctx, req.Id, existing)
	if err != nil {
		return nil, err
	}

	result, err = copierx.Trans[model.OrganizationTemplate, bo.OrganizationTemplateBo](updated)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (l *OrganizationTemplateLogic) Page(ctx context.Context, pageReq dto.OrganizationTemplatePageRequestDto) (page bo.PageResultBo[bo.OrganizationTemplatePageBo], err error) {
	pageResult, err := l.repos.OrganizationTemplateRepo.Page(ctx, pageReq)
	if err != nil {
		return page, err
	}

	pagePtr, err := copierx.Trans[gormx.PageResponse[*model.OrganizationTemplate], bo.PageResultBo[bo.OrganizationTemplatePageBo]](&pageResult)
	if err != nil {
		return page, err
	}
	return *pagePtr, nil
}

// Reset resets an organization template to its base template settings.
func (l *OrganizationTemplateLogic) Reset(ctx context.Context, id int64, organizationId int64) (result *bo.OrganizationTemplateBo, err error) {
	// Retrieve the template details with full information
	template, err := l.repos.OrganizationTemplateRepo.Detail(ctx, id, true)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	// Verify that the template belongs to the specified organization
	if template.OrganizationId != organizationId {
		return nil, errorx.New("Template Not Belong To Organization")
	}
	// Get the base template reference
	baseTemplate := template.BaseTemplate
	if baseTemplate == nil {
		return nil, errorx.New("Template Not Base Template")
	}
	// Reset all template properties to match the base template
	template.Content = baseTemplate.Content
	template.VariableGroupSchema = baseTemplate.VariableGroupSchema
	template.ExampleVariables = baseTemplate.ExampleVariables
	template.Name = baseTemplate.Name
	template.Purpose = baseTemplate.Purpose
	template.Description = baseTemplate.Description
	template.Category = baseTemplate.Category

	// Update the template in the repository
	updated, err := l.repos.OrganizationTemplateRepo.UpdateOneById(ctx, id, template)
	if err != nil {
		return nil, err
	}
	// Convert the updated model to the business object format
	result, err = copierx.Trans[model.OrganizationTemplate, bo.OrganizationTemplateBo](updated)
	if err != nil {
		return nil, err
	}

	return result, nil
}
