package dto

import (
	"time"

	"zeebee-backend/domains/zeebee/constants"
)

type RoleEntityDto struct {
	Id                int64                        `json:"id"`
	Name              string                       `json:"name"`
	Description       string                       `json:"description"`
	Type              constants.UserType           `json:"type"`
	Status            constants.CommonStatus       `json:"status"`
	CreatedAt         time.Time                    `json:"created_at"`
	UpdatedAt         time.Time                    `json:"updated_at"`
	OrganizationId    int64                        `json:"organization_id"`
	Resources         []*ResourceEntityDto         `json:"resources"`
	OrganizationUsers []*OrganizationUserEntityDto `json:"organization_users"`
}

type RolePageRequestDto struct {
	BasePageRequestDto
	Filter RolePageRequestFilter `json:"filter"`
	Sorter BasePageSorter        `json:"sorter"`
}

type RolePageRequestFilter struct {
	// UserId         int64              `json:"user_id"`
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

type RoleUpdateRequestDto struct {
	Id int64 `json:"id"`
	RoleCreateRequestDto
}

type RoleAndResourceDto struct {
	Id         int64               `json:"id"`
	RoleId     int64               `json:"role_id"`
	ResourceId int64               `json:"resource_id"`
	DataScope  constants.DataScope `json:"data_scope"`
}

type RoleAssignOrganizationUsersRequestDto struct {
	Id                  int64              `json:"id"`
	OrganizationId      int64              `json:"organization_id"`
	OrganizationUserIds []int64            `json:"organization_user_ids"`
	Type                constants.UserType `json:"type"`
}
