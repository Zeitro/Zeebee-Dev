package constants

import (
	"fmt"

	"zeebee-backend/common/enums"
)

const BorrowerPreapprovalLetterPdfName = "PreApprovalLetterPdf"
const BorrowerPreapprovalLetterEmailName = "PreApprovalLetterEmail"

const (
	OrganizationPurposePreapprovalLetter = "preApprovedLetter"
)

type TemplateCategory string

var TemplateCategoryName = "template_category"

const (
	TemplateCategoryEmail    TemplateCategory = "Email"
	TemplateCategoryContract TemplateCategory = "Contract"
)

func (c TemplateCategory) Validate() error {
	switch c {
	case TemplateCategoryEmail,
		TemplateCategoryContract:
		return nil
	}
	return fmt.Errorf("Invalid Template Category:%v", c)
}

type TemplateVariableSource string

var TemplateVariableSourceName = "template_variable_source"

const (
	TemplateVariableSourceFixed    TemplateVariableSource = "Fixed"
	TemplateVariableSourceDynamic  TemplateVariableSource = "Dynamic"
	TemplateVariableSourceTemplate TemplateVariableSource = "Template"
)

func (s TemplateVariableSource) Validate() error {
	switch s {
	case TemplateVariableSourceFixed,
		TemplateVariableSourceDynamic,
		TemplateVariableSourceTemplate:
		return nil
	}
	return fmt.Errorf("Invalid Template Variable Source:%v", s)
}

type TemplateVariableType string

var TemplateVariableTypeName = "template_variable_type"

const (
	TemplateVariableTypeString TemplateVariableType = "String"
	TemplateVariableTypeArray  TemplateVariableType = "Array"
	TemplateVariableTypeObject TemplateVariableType = "Object"
)

func (t TemplateVariableType) Validate() error {
	switch t {
	case TemplateVariableTypeString,
		TemplateVariableTypeArray,
		TemplateVariableTypeObject:
		return nil
	}
	return fmt.Errorf("Invalid Template Variable Type:%v", t)
}

type TemplateVariableFormat string

var TemplateVariableFormatName = "template_variable_format"

const (
	TemplateVariableFormatString   TemplateVariableFormat = "Email"
	TemplateVariableFormatNumber   TemplateVariableFormat = "Phone"
	TemplateVariableFormatCurrency TemplateVariableFormat = "Currency"
	TemplateVariableFormatInput    TemplateVariableFormat = "Input"
	TemplateVariableFormatTextArea TemplateVariableFormat = "TextArea"
)

func (f TemplateVariableFormat) Validate() error {
	switch f {
	case TemplateVariableFormatString,
		TemplateVariableFormatNumber,
		TemplateVariableFormatCurrency,
		TemplateVariableFormatInput,
		TemplateVariableFormatTextArea:
		return nil
	}
	return fmt.Errorf("Invalid Template Variable Format:%v", f)
}

type BaseTemplateScope string

var BaseTemplateScopeName = "base_template_scope"

const (
	BaseTemplateScopeSystem       BaseTemplateScope = "System"
	BaseTemplateScopeOrganization BaseTemplateScope = "Organization"
)

func (s BaseTemplateScope) Validate() error {
	switch s {
	case BaseTemplateScopeSystem,
		BaseTemplateScopeOrganization:
		return nil
	}
	return fmt.Errorf("Invalid Base Template Scope:%v", s)
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
	// register template variable format
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
			Label: string(TemplateVariableFormatCurrency),
			Value: string(TemplateVariableFormatCurrency),
		},
		{
			Label: string(TemplateVariableFormatInput),
			Value: string(TemplateVariableFormatInput),
		},
		{
			Label: string(TemplateVariableFormatTextArea),
			Value: string(TemplateVariableFormatTextArea),
		},
	})
	// register template scope
	enums.Register(BaseTemplateScopeName, []enums.EnumOpt{
		{
			Label: string(BaseTemplateScopeSystem),
			Value: string(BaseTemplateScopeSystem),
		},
		{
			Label: string(BaseTemplateScopeOrganization),
			Value: string(BaseTemplateScopeOrganization),
		},
	})
	// register template variable type
	enums.Register(TemplateVariableTypeName, []enums.EnumOpt{
		{
			Label: string(TemplateVariableTypeString),
			Value: string(TemplateVariableTypeString),
		},
		{
			Label: string(TemplateVariableTypeArray),
			Value: string(TemplateVariableTypeArray),
		},
		{
			Label: string(TemplateVariableTypeObject),
			Value: string(TemplateVariableTypeObject),
		},
	})
}
