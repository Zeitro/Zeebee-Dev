package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"zeebee-backend/domains/zeebee/constants"
)

// BaseTemplate represents a base template in the system
type BaseTemplate struct {
	Base
	Name                string                     `json:"name" gorm:"column:name;type:varchar(100);not null;default:''"`
	Type                constants.UserType         `json:"type" gorm:"column:type;type:varchar(20);not null;default:''"`         // Borrower,Lender
	Category            constants.TemplateCategory `json:"category" gorm:"column:category;type:varchar(20);not null;default:''"` // Email,Contract
	Tag                 string                     `json:"tag" gorm:"column:tag;type:varchar(100);not null;default:''"`
	Content             string                     `json:"content" gorm:"column:content;type:text;not null;default:''"`
	VariableGroupSchema []*TemplateVariableGroup   `json:"variable_group_schema" gorm:"column:variable_group_schema;type:jsonb;not null;default:'[]';serializer:json"`
	Status              constants.CommonStatus     `json:"status" gorm:"column:status;type:varchar(20);not null;default:''"` // Active,Disabled
	Description         string                     `json:"description" gorm:"column:description;type:text;not null;default:''"`
	ExampleVariables    string                     `json:"example_variables,omitempty" gorm:"column:example_variables;type:jsonb;not null;default:'{}'"`

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

// TemplateVariable represents a variable definition in the template schema
type TemplateVariable struct {
	Label      string                           `json:"label"`
	Type       string                           `json:"type"` // string,number,boolean,array,object
	Key        string                           `json:"key"`
	IsEditable bool                             `json:"is_editable"`
	Source     constants.TemplateVariableSource `json:"source"` // Fixed,Dynamic,Template
	FixedValue string                           `json:"fixed_value"`
	Format     constants.TemplateVariableFormat `json:"format,omitempty"` // Email,Phone,Currency
	// ExampleValue string                           `json:"example_value,omitempty"`
	Properties map[string]*TemplateVariable `json:"properties,omitempty"` // for object type
	Items      []*TemplateVariable          `json:"items,omitempty"`      // for array type
}

// Scan implements sql.Scanner interface for TemplateVariable
func (tv *TemplateVariable) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, tv)
}

// Value implements driver.Valuer interface for TemplateVariable
func (tv TemplateVariable) Value() (driver.Value, error) {
	return json.Marshal(tv)
}

// Validate checks if the TemplateVariable has valid configuration
func (tv *TemplateVariable) Validate() error {
	validTypes := map[string]bool{
		"string": true, "number": true, "boolean": true,
		"array": true, "object": true,
	}
	if !validTypes[tv.Type] {
		return errors.New("invalid variable type: " + tv.Type)
	}

	// Validate nested structures based on type
	if tv.Type == "object" && tv.Properties != nil {
		for _, prop := range tv.Properties {
			if err := prop.Validate(); err != nil {
				return err
			}
		}
	} else if tv.Type == "array" && tv.Items != nil {
		for _, item := range tv.Items {
			if err := item.Validate(); err != nil {
				return err
			}
		}
	}

	return nil
}

// TemplateVariableGroup represents a group of variables
type TemplateVariableGroup struct {
	Name      string              `json:"name"`
	Variables []*TemplateVariable `json:"variables"`
}

// Scan implements sql.Scanner interface for TemplateVariableGroup
func (tvg *TemplateVariableGroup) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, tvg)
}

// Value implements driver.Valuer interface for TemplateVariableGroup
func (tvg TemplateVariableGroup) Value() (driver.Value, error) {
	return json.Marshal(tvg)
}
