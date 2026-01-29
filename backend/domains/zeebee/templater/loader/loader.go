package loader

import (
	"context"

	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/logic"
)

type TemplateLoader interface {
	LoadSystemTemplate(ctx context.Context, purpose string, organizationType string) (string, error)
	LoadOrganizationTemplate(ctx context.Context, purpose string, organizationId int64) (string, error)
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

func (l *DefaultTemplateLoader) LoadSystemTemplate(ctx context.Context, purpose string, organizationType string) (content string, err error) {
	content, err = l.baseTemplateLogic.GetContentByPurpose(ctx, purpose, constants.UserType(organizationType))
	return content, err
}

func (l *DefaultTemplateLoader) LoadOrganizationTemplate(ctx context.Context, purpose string, organizationId int64) (content string, err error) {
	tpl, err := l.organizationTemplateLogic.DetailByPurpose(ctx, organizationId, purpose)
	if err != nil {
		return "", err
	}
	content = tpl.Content
	return content, err
}
