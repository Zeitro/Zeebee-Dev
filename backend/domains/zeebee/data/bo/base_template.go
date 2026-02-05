package bo

import (
	"time"

	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data"
)

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
	IsConfigurable      bool                          `json:"is_configurable"`
	CreatedAt           time.Time                     `json:"created_at"`
	UpdatedAt           time.Time                     `json:"updated_at"`
}

// BaseTemplatePageBo represents a base template in page list
type BaseTemplatePageBo struct {
	Id             int64                       `json:"id"`
	Name           string                      `json:"name"`
	Type           constants.UserType          `json:"type"`
	Category       constants.TemplateCategory  `json:"category"`
	Purpose        string                      `json:"purpose"`
	Status         constants.CommonStatus      `json:"status"`
	Scope          constants.BaseTemplateScope `json:"scope"`
	IsConfigurable bool                        `json:"is_configurable"`
	Description    string                      `json:"description"`
	CreatedAt      time.Time                   `json:"created_at"`
	UpdatedAt      time.Time                   `json:"updated_at"`
}
