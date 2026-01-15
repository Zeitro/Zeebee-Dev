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
type OrganizationRegisterRequestDto struct {
	Id           int64                           `json:"id"`
	Name         string                          `json:"name"`
	Status       constants.OrganizationStatus    `json:"status"`
	OwnerType    constants.OrganizationOwnerType `json:"owner_type"`
	Type         constants.UserType              `json:"type"`
	ContactEmail string                          `json:"contact_email"`
	ContactPhone string                          `json:"contact_phone"`
	ContactName  string                          `json:"contact_name"`
	Logo         string                          `json:"logo"`
	Website      string                          `json:"website"`
}
