package dto

import "zeebee-backend/domains/zeebee/constants"

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
	Id                  int64                       `json:"id"`
	OrganizationId      int64                       `json:"organization_id"`
	VariableGroupSchema []*TemplateVariableGroupDto `json:"variable_group_schema"`
}

type OrganizationTemplatePreviewDto struct {
	Name                string                      `json:"name"`
	Content             string                      `json:"content"`
	VariableGroupSchema []*TemplateVariableGroupDto `json:"variable_group_schema"`
	ExampleVariables    string                      `json:"example_variables"`
}
