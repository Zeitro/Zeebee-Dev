package templater

import (
	"context"
	"fmt"

	"zeebee-backend/domains/zeebee/templater/convertor"
	"zeebee-backend/domains/zeebee/templater/loader"
)

type Engine struct {
	TemplateLoader loader.TemplateLoader
	Convertor      convertor.Convertor
}

func NewEngine(templateLoader loader.TemplateLoader, convertor convertor.Convertor) *Engine {
	return &Engine{
		TemplateLoader: templateLoader,
		Convertor:      convertor,
	}
}

func (e *Engine) Export(ctx context.Context, conf TemplateConf, data map[string]interface{}) (content []byte, contentType string, err error) {
	var template string
	switch conf.Scope {
	case ScopeSystem:
		template, err = e.TemplateLoader.LoadSystemTemplate(ctx, conf.Purpose, conf.OrganizationType)
	case ScopeOrganization:
		template, err = e.TemplateLoader.LoadOrganizationTemplate(ctx, conf.Purpose, conf.OrganizationId)
	default:
		err = fmt.Errorf("Unknown Scope: %s", conf.Scope)
	}
	if err != nil {
		return
	}
	if len(template) == 0 {
		err = fmt.Errorf("Template Is Empty")
		return
	}

	switch conf.FileType {
	case FileTypeHTML:
		return e.Convertor.ConvertToHtml(ctx, conf.Purpose, template, data)
	case FileTypePDF:
		return e.Convertor.ConvertToPdf(ctx, conf.Purpose, template, data)
	default:
		return nil, "", fmt.Errorf("Unknown File Type: %s", conf.FileType)
	}
}
