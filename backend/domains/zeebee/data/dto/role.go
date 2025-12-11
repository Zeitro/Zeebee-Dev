package dto

import "zeebee-backend/common/constants"

type RolePageRequestDto struct {
	BasePageRequestDto
	Filter RolePageRequestFilter `json:"filter"`
	Sorter BasePageSorter        `json:"sorter"`
}

type RolePageRequestFilter struct {
	UserId         int64              `json:"user_id"`
	OrganizationId int64              `json:"organization_id"`
	Keyword        string             `json:"keyword"`
	RoleType       constants.UserType `json:"role_type"`
}

type RoleCreateRequestDto struct {
	Name             string                 `json:"name"`
	Type             constants.UserType     `json:"type"`
	Description      string                 `json:"description"`
	OrganizationId   int64                  `json:"organization_id"`
	RoleAndResources []RoleAndResourceDto   `json:"role_and_resources"`
	Status           constants.CommonStatus `json:"status"`
}

type RoleAndResourceDto struct {
	Id         int64               `json:"id"`
	RoleId     int64               `json:"role_id"`
	ResourceId int64               `json:"resource_id"`
	DataScope  constants.DataScope `json:"data_scope"`
}
