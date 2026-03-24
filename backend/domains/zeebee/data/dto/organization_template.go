package dto

import (
	"time"

	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data"
)

type OrganizationTemplateEntityDto struct {
	Id                  int64                         `json:"id"`
	Name                string                        `json:"name"`
	Type                constants.UserType            `json:"type"`
	Category            constants.TemplateCategory    `json:"category"`
	Purpose             string                        `json:"purpose"`
	Content             string                        `json:"content"`
	VariableGroupSchema []*data.TemplateVariableGroup `json:"variable_group_schema"`
	Status              constants.CommonStatus        `json:"status"`
	HasConfigured       bool                          `json:"has_configured"`
	Description         string                        `json:"description"`
	ExampleVariables    string                        `json:"example_variables"`
	CreatedAt           time.Time                     `json:"created_at"`
	UpdatedAt           time.Time                     `json:"updated_at"`

	BaseTemplate *BaseTemplateEntityDto `json:"base_template,omitempty"`
}

// OrganizationTemplatePageBo represents a base template in page list
type OrganizationTemplatePageEntityDto struct {
	Id            int64                      `json:"id"`
	Name          string                     `json:"name"`
	Type          constants.UserType         `json:"type"`
	Category      constants.TemplateCategory `json:"category"`
	Purpose       string                     `json:"purpose"`
	Status        constants.CommonStatus     `json:"status"`
	HasConfigured bool                       `json:"has_configured"`
	Description   string                     `json:"description"`
	CreatedAt     time.Time                  `json:"created_at"`
	UpdatedAt     time.Time                  `json:"updated_at"`
}

// OrganizationTemplatePageRequestDto represents the request for paginated base template list
type OrganizationTemplatePageRequestDto struct {
	OrganizationPageRequestDto
	Sorter BasePageSorter
	Filter OrganizationTemplatePageFilter
}

// OrganizationTemplatePageFilter represents filter criteria for base template page queries
type OrganizationTemplatePageFilter struct {
	Keyword        string                     `json:"keyword"`
	Type           constants.UserType         `json:"type"`
	Category       constants.TemplateCategory `json:"category"`
	Status         constants.CommonStatus     `json:"status"`
	OrganizationId int64                      `json:"organization_id"`
}

// OrganizationTemplateUpdateDto represents the request for updating a base template
type OrganizationTemplateUpdateVariableGroupDto struct {
	Id                  int64                         `json:"id"`
	OrganizationId      int64                         `json:"organization_id"`
	VariableGroupSchema []*data.TemplateVariableGroup `json:"variable_group_schema"`
}

type OrganizationTemplatePreviewDto struct {
	Name                string                        `json:"name"`
	Content             string                        `json:"content"`
	VariableGroupSchema []*data.TemplateVariableGroup `json:"variable_group_schema"`
	ExampleVariables    string                        `json:"example_variables"`
}
