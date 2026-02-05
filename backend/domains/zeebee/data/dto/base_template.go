package dto

import (
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data"
)

// BaseTemplatePageRequestDto represents the request for paginated base template list
type BaseTemplatePageRequestDto struct {
	BasePageRequestDto
	Sorter BasePageSorter
	Filter BaseTemplatePageFilter
}

// BaseTemplatePageFilter represents filter criteria for base template page queries
type BaseTemplatePageFilter struct {
	Keyword  string                      `json:"keyword"`
	Type     constants.UserType          `json:"type"`
	Category constants.TemplateCategory  `json:"category"`
	Status   constants.CommonStatus      `json:"status"`
	Scope    constants.BaseTemplateScope `json:"scope"`
}

// BaseTemplateCreateDto represents the request for creating a base template
type BaseTemplateCreateDto struct {
	Name                string                        `json:"name"`
	Type                constants.UserType            `json:"type"`
	Category            constants.TemplateCategory    `json:"category"`
	Purpose             string                        `json:"purpose"`
	Content             string                        `json:"content"`
	VariableGroupSchema []*data.TemplateVariableGroup `json:"variable_group_schema"`
	ExampleVariables    string                        `json:"example_variables"`
	Description         string                        `json:"description"`
	Scope               constants.BaseTemplateScope   `json:"scope"`
	IsConfigurable      bool                          `json:"is_configurable"`
}

// BaseTemplateUpdateDto represents the request for updating a base template
type BaseTemplateUpdateDto struct {
	Id                  int64                         `json:"id"`
	Name                string                        `json:"name"`
	Category            constants.TemplateCategory    `json:"category"`
	Content             string                        `json:"content"`
	VariableGroupSchema []*data.TemplateVariableGroup `json:"variable_group_schema"`
	Description         string                        `json:"description"`
	ExampleVariables    string                        `json:"example_variables"`
	IsConfigurable      bool                          `json:"is_configurable"`
	Scope               constants.BaseTemplateScope   `json:"scope"`
}

type BaseTemplatePreviewDto struct {
	Name                string                        `json:"name"`
	Content             string                        `json:"content"`
	VariableGroupSchema []*data.TemplateVariableGroup `json:"variable_group_schema"`
	ExampleVariables    string                        `json:"example_variables"`
}
