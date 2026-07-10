package loader

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/logic"
)

var ErrTemplateNotActive = fmt.Errorf("Template is not active")

type TemplateLoader interface {
	LoadSystemTemplate(ctx context.Context, purpose string, organizationType string) (string, map[string]interface{}, error)
	LoadOrganizationTemplate(ctx context.Context, purpose string, organizationId int64) (string, map[string]interface{}, error)
}

type DefaultTemplateLoader struct {
	baseTemplateLogic         *logic.BaseTemplateLogic
	organizationTemplateLogic *logic.OrganizationTemplateLogic
}

func NewDefaultTemplateLoader(baseTemplateLogic *logic.BaseTemplateLogic, organizationTemplateLogic *logic.OrganizationTemplateLogic) *DefaultTemplateLoader {
	return &DefaultTemplateLoader{
		baseTemplateLogic:         baseTemplateLogic,
		organizationTemplateLogic: organizationTemplateLogic,
	}
}

func (l *DefaultTemplateLoader) LoadSystemTemplate(ctx context.Context, purpose string, organizationType string) (content string, templateRenderData map[string]interface{}, err error) {
	tpl, err := l.baseTemplateLogic.DetailByPurpose(ctx, purpose, constants.UserType(organizationType))
	if err != nil {
		return "", nil, err
	}
	if tpl.Status != constants.CommonStatusActive {
		return "", nil, ErrTemplateNotActive
	}
	content = tpl.Content
	templateRenderData = getTemplateRenderData(tpl.VariableGroupSchema)

	return content, templateRenderData, err
}

func (l *DefaultTemplateLoader) LoadOrganizationTemplate(ctx context.Context, purpose string, organizationId int64) (content string, templateRenderData map[string]interface{}, err error) {
	tpl, err := l.organizationTemplateLogic.DetailByPurpose(ctx, organizationId, purpose)
	if err != nil {
		return "", nil, err
	}
	if tpl.Status != constants.CommonStatusActive {
		return "", nil, ErrTemplateNotActive
	}
	if !tpl.HasConfigured {
		return "", nil, errors.New("Template is not configured")
	}

	content = tpl.Content
	templateRenderData = getTemplateRenderData(tpl.VariableGroupSchema)
	return content, templateRenderData, err
}

func getTemplateRenderData(variableGroups []*dto.TemplateVariableGroup) (templateRenderData map[string]interface{}) {
	templateRenderData = make(map[string]interface{})

	for _, variableGroup := range variableGroups {
		for _, variable := range variableGroup.Variables {
			if variable.Source == constants.TemplateVariableSourceFixed || variable.Source == constants.TemplateVariableSourceTemplate {
				k := strings.Trim(variable.Key, ".")
				templateRenderData[k] = variable.FixedValue
			}
		}
	}
	return templateRenderData
}
