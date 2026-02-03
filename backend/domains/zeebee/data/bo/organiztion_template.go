package bo

import (
	"time"

	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data"
)

// OrganizationTemplateBo represents a base template in the system
type OrganizationTemplateBo struct {
	Id                  int64                         `json:"id"`
	Name                string                        `json:"name"`
	Type                constants.UserType            `json:"type"`
	Category            constants.TemplateCategory    `json:"category"`
	Purpose             string                        `json:"purpose"`
	Content             string                        `json:"content"`
	VariableGroupSchema []*data.TemplateVariableGroup `json:"variable_group_schema"`
	Status              constants.CommonStatus        `json:"status"`
	Description         string                        `json:"description"`
	ExampleVariables    string                        `json:"example_variables"`
	CreatedAt           time.Time                     `json:"created_at"`
	UpdatedAt           time.Time                     `json:"updated_at"`

	BaseTemplate *BaseTemplateBo `json:"base_template,omitempty"`
}

// OrganizationTemplatePageBo represents a base template in page list
type OrganizationTemplatePageBo struct {
	Id          int64                      `json:"id"`
	Name        string                     `json:"name"`
	Type        constants.UserType         `json:"type"`
	Category    constants.TemplateCategory `json:"category"`
	Purpose     string                     `json:"purpose"`
	Status      constants.CommonStatus     `json:"status"`
	Description string                     `json:"description"`
	CreatedAt   time.Time                  `json:"created_at"`
	UpdatedAt   time.Time                  `json:"updated_at"`
}
