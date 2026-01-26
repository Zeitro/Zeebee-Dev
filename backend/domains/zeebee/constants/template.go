package constants

import (
	"fmt"

	"zeebee-backend/common/enums"
)

type TemplateCategory string

var TemplateCategoryName = "template_category"

const (
	TemplateCategoryEmail    TemplateCategory = "Email"
	TemplateCategoryContract TemplateCategory = "Contract"
)

func (c TemplateCategory) Validate() error {
	if c != TemplateCategoryEmail && c != TemplateCategoryContract {
		return fmt.Errorf("Invalid Template Category:%v", c)
	}
	return nil
}

type TemplateVariableSource string

var TemplateVariableSourceName = "template_variable_source"

const (
	TemplateVariableSourceFixed    TemplateVariableSource = "Fixed"
	TemplateVariableSourceDynamic  TemplateVariableSource = "Dynamic"
	TemplateVariableSourceTemplate TemplateVariableSource = "Template"
)

func (s TemplateVariableSource) Validate() error {
	if s != TemplateVariableSourceFixed && s != TemplateVariableSourceDynamic && s != TemplateVariableSourceTemplate {
		return fmt.Errorf("Invalid Template Variable Source:%v", s)
	}
	return nil
}

type TemplateVariableFormat string

var TemplateVariableFormatName = "template_variable_format"

const (
	TemplateVariableFormatString  TemplateVariableFormat = "Email"
	TemplateVariableFormatNumber  TemplateVariableFormat = "Phone"
	TemplateVariableFormatBoolean TemplateVariableFormat = "Currency"
)

func (f TemplateVariableFormat) Validate() error {
	if f != TemplateVariableFormatString && f != TemplateVariableFormatNumber && f != TemplateVariableFormatBoolean {
		return fmt.Errorf("Invalid Template Variable Format:%v", f)
	}
	return nil
}

func init() {
	// register template category
	enums.Register(TemplateCategoryName, []enums.EnumOpt{
		{
			Label: string(TemplateCategoryEmail),
			Value: string(TemplateCategoryEmail),
		},
		{
			Label: string(TemplateCategoryContract),
			Value: string(TemplateCategoryContract),
		},
	})
	// register template variable source
	enums.Register(TemplateVariableSourceName, []enums.EnumOpt{
		{
			Label: string(TemplateVariableSourceFixed),
			Value: string(TemplateVariableSourceFixed),
		},
		{
			Label: string(TemplateVariableSourceDynamic),
			Value: string(TemplateVariableSourceDynamic),
		},
		{
			Label: string(TemplateVariableSourceTemplate),
			Value: string(TemplateVariableSourceTemplate),
		},
	})
	enums.Register(TemplateVariableFormatName, []enums.EnumOpt{
		{
			Label: string(TemplateVariableFormatString),
			Value: string(TemplateVariableFormatString),
		},
		{
			Label: string(TemplateVariableFormatNumber),
			Value: string(TemplateVariableFormatNumber),
		},
		{
			Label: string(TemplateVariableFormatBoolean),
			Value: string(TemplateVariableFormatBoolean),
		},
	})
}
