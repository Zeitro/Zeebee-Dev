package bo

import (
	"time"

	"zeebee-backend/common/constants"
)

type ResourceBo struct {
	Id                 int64                  `json:"id"`
	ParentId           int64                  `json:"parent_id"`
	Type               constants.UserType     `json:"type"`
	ResourceType       constants.ResourceType `json:"resource_type"`
	Name               string                 `json:"name"`
	Path               string                 `json:"path"`
	Icon               string                 `json:"icon"`
	Tag                string                 `json:"tag"`
	IsShow             bool                   `json:"is_show"`
	Status             constants.CommonStatus `json:"status"`
	SortOrder          int32                  `json:"sort_order"`
	CreatedAt          time.Time              `json:"created_at"`
	UpdatedAt          time.Time              `json:"updated_at"`
	PivotRoleDataScope constants.DataScope    `json:"pivot_role_data_scope"`
}
