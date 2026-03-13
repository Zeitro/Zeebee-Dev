package dto

import (
	"time"

	"zeebee-backend/domains/zeebee/constants"
)

type ResourceTreeDto struct {
	Id              int64                  `json:"id"`
	ParentId        int64                  `json:"parent_id"`
	Name            string                 `json:"name"`
	Path            string                 `json:"path"`
	ResourceType    constants.ResourceType `json:"resource_type"`
	Status          constants.CommonStatus `json:"status"`
	Description     string                 `json:"description"`
	Type            constants.UserType     `json:"type"`
	Icon            string                 `json:"icon"`
	Tag             string                 `json:"tag"`
	SortOrder       int32                  `json:"sort_order"`
	IsShow          bool                   `json:"is_show"`
	EnableDataScope bool                   `json:"enable_data_scope"`
	Apis            []*ApiEntityDto        `json:"apis"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
	Children        []*ResourceTreeDto     `json:"children,omitempty"`
}

type ResourceEntityDto struct {
	Id              int64                  `json:"id"`
	ParentId        int64                  `json:"parent_id"`
	Name            string                 `json:"name"`
	Path            string                 `json:"path"`
	ResourceType    constants.ResourceType `json:"resource_type"`
	Status          constants.CommonStatus `json:"status"`
	Description     string                 `json:"description"`
	Type            constants.UserType     `json:"type"`
	Icon            string                 `json:"icon"`
	Tag             string                 `json:"tag"`
	SortOrder       int32                  `json:"sort_order"`
	IsShow          bool                   `json:"is_show"`
	DataScope       constants.DataScope    `json:"data_scope"`
	EnableDataScope bool                   `json:"enable_data_scope"`
	Apis            []*ApiEntityDto        `json:"apis"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
}

type ResourceCreateRequestDto struct {
	ParentId        int64                  `json:"parent_id"`
	Type            constants.UserType     `json:"type"`
	ResourceType    constants.ResourceType `json:"resource_type"`
	Name            string                 `json:"name"`
	Path            string                 `json:"path"`
	Icon            string                 `json:"icon"`
	Tag             string                 `json:"tag"`
	IsShow          bool                   `json:"is_show"`
	Status          constants.CommonStatus `json:"status"`
	SortOrder       int32                  `json:"sort_order"`
	ApiIds          []int64                `json:"api_ids"`
	EnableDataScope bool                   `json:"enable_data_scope"`
}

type ResourceUpdateRequestDto struct {
	Id int64 `json:"id"`
	ResourceCreateRequestDto
}

type ResourceUserMenuRequestDto struct {
	OrganizationId     int64              `json:"organization_id"`
	OrganizationUserId int64              `json:"organization_user_id"`
	Type               constants.UserType `json:"type"`
}
