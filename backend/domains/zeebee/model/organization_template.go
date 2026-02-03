package model

import (
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data"
)

// OrganizationTemplate represents an organization's customized template
type OrganizationTemplate struct {
	Base
	OrganizationId      int64                         `json:"organization_id" gorm:"column:organization_id;default:0;index"`
	BaseTemplateId      int64                         `json:"base_template_id" gorm:"column:base_template_id;default:0;index"`
	Name                string                        `json:"name" gorm:"column:name;type:varchar(100);not null;default:''"`
	Type                constants.UserType            `json:"type" gorm:"column:type;type:varchar(20);not null;default:''"`         // Borrower,Lender
	Category            constants.TemplateCategory    `json:"category" gorm:"column:category;type:varchar(20);not null;default:''"` // Email,Contract
	Purpose             string                        `json:"purpose" gorm:"column:purpose;type:varchar(100);not null;default:''"`
	Content             string                        `json:"content" gorm:"column:content;type:text;not null;default:''"`
	VariableGroupSchema []*data.TemplateVariableGroup `json:"variable_group_schema" gorm:"column:variable_group_schema;type:jsonb;not null;default:'[]';serializer:json"`
	Status              constants.CommonStatus        `json:"status" gorm:"column:status;type:varchar(20);not null;default:''"` // Active,Disabled
	Description         string                        `json:"description" gorm:"column:description;type:text;not null;default:''"`
	ExampleVariables    string                        `json:"example_variables,omitempty" gorm:"column:example_variables;type:jsonb;not null;default:'{}'"`

	// Relationships
	BaseTemplate *BaseTemplate `json:"base_template,omitempty" gorm:"foreignKey:BaseTemplateId;references:Id"`
	Organization *Organization `json:"organization,omitempty" gorm:"foreignKey:OrganizationId;references:Id"`
}

// TableName specifies the table name for OrganizationTemplate
func (OrganizationTemplate) TableName() string {
	return "organization_templates"
}

// RemoveRelationships removes relationship data to avoid circular references
func (ot *OrganizationTemplate) RemoveRelationships() {
	ot.BaseTemplate = nil
	ot.Organization = nil
}
