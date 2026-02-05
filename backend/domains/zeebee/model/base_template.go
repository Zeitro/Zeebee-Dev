package model

import (
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data"
)

// BaseTemplate represents a base template in the system
type BaseTemplate struct {
	Base
	Name                string                        `json:"name" gorm:"column:name;type:varchar(100);not null;default:''"`
	Type                constants.UserType            `json:"type" gorm:"column:type;type:varchar(20);not null;default:''"`         // Borrower,Lender
	Category            constants.TemplateCategory    `json:"category" gorm:"column:category;type:varchar(20);not null;default:''"` // Email,Contract
	Purpose             string                        `json:"purpose" gorm:"column:purpose;type:varchar(100);not null;default:''"`
	Content             string                        `json:"content" gorm:"column:content;type:text;not null;default:''"`
	VariableGroupSchema []*data.TemplateVariableGroup `json:"variable_group_schema" gorm:"column:variable_group_schema;type:jsonb;not null;default:'[]';serializer:json"`
	Status              constants.CommonStatus        `json:"status" gorm:"column:status;type:varchar(20);not null;default:''"` // Active,Disabled
	Scope               constants.BaseTemplateScope   `json:"scope" gorm:"column:scope;type:varchar(20);not null;default:''"`
	Description         string                        `json:"description" gorm:"column:description;type:text;not null;default:''"`
	ExampleVariables    string                        `json:"example_variables,omitempty" gorm:"column:example_variables;type:jsonb;not null;default:'{}'"`
	IsConfigurable      bool                          `json:"is_configurable" gorm:"column:is_configurable;type:boolean;not null;default:false"`

	// Relationships - derived templates that use this base template
	OrganizationTemplates []*OrganizationTemplate `json:"organization_templates,omitempty" gorm:"foreignKey:BaseTemplateId;references:Id"`
}

// TableName specifies the table name for BaseTemplate
func (BaseTemplate) TableName() string {
	return "base_templates"
}

// BeforeValidate runs validation before saving
func (bt *BaseTemplate) BeforeValidate() error {
	// Validate variable group schema
	for _, group := range bt.VariableGroupSchema {
		for _, variable := range group.Variables {
			if err := variable.Validate(); err != nil {
				return err
			}
		}
	}
	return nil
}

// RemoveRelationships removes relationship data to avoid circular references
func (bt *BaseTemplate) RemoveRelationships() {
	bt.OrganizationTemplates = nil
}
