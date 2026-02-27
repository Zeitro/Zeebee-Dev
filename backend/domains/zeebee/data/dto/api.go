package dto

import (
	"time"

	"zeebee-backend/domains/zeebee/constants"
)

type ApiEntityDto struct {
	Id                 int64                  `json:"id"`
	ParentId           int64                  `json:"parent_id"`
	Name               string                 `json:"name"`
	Path               string                 `json:"path"`
	Method             string                 `json:"method"`
	Status             constants.CommonStatus `json:"status"`
	Description        string                 `json:"description"`
	Type               constants.UserType     `json:"type"`
	SortOrder          int32                  `json:"sort_order"`
	AuthRequired       bool                   `json:"auth_required"`
	PermissionRequired bool                   `json:"permission_required"`
	CreatedAt          time.Time              `json:"created_at"`
	UpdatedAt          time.Time              `json:"updated_at"`
	Children           []*ApiEntityDto        `json:"children,omitempty"`
}

type ApiCreateRequestDto struct {
	ParentId           int64                  `json:"parent_id"`
	Name               string                 `json:"name"`
	Path               string                 `json:"path"`
	Method             string                 `json:"method"`
	Description        string                 `json:"description"`
	Type               constants.UserType     `json:"type"`
	Status             constants.CommonStatus `json:"status"`
	SortOrder          int32                  `json:"sort_order"`
	AuthRequired       bool                   `json:"auth_required"`
	PermissionRequired bool                   `json:"permission_required"`
}

type ApiUpdateRequestDto struct {
	Id                 int64                  `json:"id"`
	ParentId           int64                  `json:"parent_id"`
	Name               string                 `json:"name"`
	Path               string                 `json:"path"`
	Method             string                 `json:"method"`
	Description        string                 `json:"description"`
	Type               constants.UserType     `json:"type"`
	SortOrder          int32                  `json:"sort_order"`
	Status             constants.CommonStatus `json:"status"`
	AuthRequired       bool                   `json:"auth_required"`
	PermissionRequired bool                   `json:"permission_required"`
}
