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
