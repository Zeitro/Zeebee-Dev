package dto

import "zeebee-backend/domains/zeebee/constants"

type OrganizationPageRequestDto struct {
	BasePageRequestDto
	Sorter BasePageSorter
	Filter OrganizationPageFilter
}
type OrganizationPageFilter struct {
	Keyword          string                       `json:"keyword"`
	OrganizationType constants.UserType           `json:"organization_type"`
	Status           constants.OrganizationStatus `json:"status"`
}
