package dto

import "zeebee-backend/domains/zeebee/constants"

// BaseTemplatePageRequestDto represents the request for paginated base template list
type BaseTemplatePageRequestDto struct {
	BasePageRequestDto
	Sorter BasePageSorter
	Filter BaseTemplatePageFilter
}

// BaseTemplatePageFilter represents filter criteria for base template page queries
type BaseTemplatePageFilter struct {
	Keyword  string                     `json:"keyword"`
	Type     constants.UserType         `json:"type"`
	Category constants.TemplateCategory `json:"category"`
	Status   constants.CommonStatus     `json:"status"`
}

// BaseTemplateCreateDto represents the request for creating a base template
type BaseTemplateCreateDto struct {
	Name                string                      `json:"name"`
	Type                constants.UserType          `json:"type"`
	Category            constants.TemplateCategory  `json:"category"`
	Tag                 string                      `json:"tag"`
	Content             string                      `json:"content"`
	VariableGroupSchema []*TemplateVariableGroupDto `json:"variable_group_schema"`
	ExampleVariables    string                      `json:"example_variables"`
	Description         string                      `json:"description"`
}

// BaseTemplateUpdateDto represents the request for updating a base template
type BaseTemplateUpdateDto struct {
	Id                  int64                       `json:"id"`
	Name                string                      `json:"name"`
	Category            constants.TemplateCategory  `json:"category"`
	Content             string                      `json:"content"`
	VariableGroupSchema []*TemplateVariableGroupDto `json:"variable_group_schema"`
	Description         string                      `json:"description"`
	ExampleVariables    string                      `json:"example_variables"`
}

// TemplateVariableGroupDto represents a group of variables in DTO
type TemplateVariableGroupDto struct {
	Name      string                 `json:"name"`
	Variables []*TemplateVariableDto `json:"variables"`
}

// TemplateVariableDto represents a variable definition in DTO
type TemplateVariableDto struct {
	Label      string                           `json:"label"`
	Type       string                           `json:"type"`
	Key        string                           `json:"key"`
	IsEditable bool                             `json:"is_editable"`
	Source     string                           `json:"source"`
	Format     constants.TemplateVariableFormat `json:"format,omitempty"`
	FixedValue string                           `json:"fixed_value"`
	Properties map[string]*TemplateVariableDto  `json:"properties,omitempty"`
	Items      []*TemplateVariableDto           `json:"items,omitempty"`
}
