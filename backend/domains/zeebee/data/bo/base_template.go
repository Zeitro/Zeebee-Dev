package bo

import (
	"time"

	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data"
)

// TemplateVariableBo represents a variable definition in the template schema
// type TemplateVariableBo struct {
// 	Label      string                           `json:"label"`
// 	Type       constants.TemplateVariableType   `json:"type"` // String,Array,Object
// 	Key        string                           `json:"key"`
// 	IsEditable bool                             `json:"is_editable"`
// 	IsRequired bool                             `json:"is_required"`
// 	Source     constants.TemplateVariableSource `json:"source"` // Fixed/Dynamic
// 	FixedValue string                           `json:"fixed_value"`
// 	Format     constants.TemplateVariableFormat `json:"format,omitempty"`
// 	Properties map[string]*TemplateVariableBo   `json:"properties,omitempty"` // for Object type
// 	Items      []*TemplateVariableBo            `json:"items,omitempty"`      // for Array type
//
// }

// TemplateVariableGroupBo represents a group of variables
// type TemplateVariableGroupBo struct {
// 	// Name      string                `json:"name"`
// 	// Variables []*TemplateVariableBo `json:"variables"`
// 	data.TemplateVariableGroup
// }

// BaseTemplateBo represents a base template in the system
type BaseTemplateBo struct {
	Id                  int64                         `json:"id"`
	Name                string                        `json:"name"`
	Type                constants.UserType            `json:"type"`
	Category            constants.TemplateCategory    `json:"category"`
	Purpose             string                        `json:"purpose"`
	Content             string                        `json:"content"`
	VariableGroupSchema []*data.TemplateVariableGroup `json:"variable_group_schema"`
	Status              constants.CommonStatus        `json:"status"`
	Scope               constants.BaseTemplateScope   `json:"scope"`
	Description         string                        `json:"description"`
	ExampleVariables    string                        `json:"example_variables"`
	CreatedAt           time.Time                     `json:"created_at"`
	UpdatedAt           time.Time                     `json:"updated_at"`
}

// BaseTemplatePageBo represents a base template in page list
type BaseTemplatePageBo struct {
	Id          int64                       `json:"id"`
	Name        string                      `json:"name"`
	Type        constants.UserType          `json:"type"`
	Category    constants.TemplateCategory  `json:"category"`
	Purpose     string                      `json:"purpose"`
	Status      constants.CommonStatus      `json:"status"`
	Scope       constants.BaseTemplateScope `json:"scope"`
	Description string                      `json:"description"`
	CreatedAt   time.Time                   `json:"created_at"`
	UpdatedAt   time.Time                   `json:"updated_at"`
}
