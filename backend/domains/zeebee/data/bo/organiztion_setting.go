package bo

import (
	"time"
)

// OrganizationSettingBo represents an organization setting in the system
type OrganizationSettingBo struct {
	Id             int64                  `json:"id"`
	Purpose        string                 `json:"purpose"`
	OrganizationId int64                  `json:"organization_id"`
	Enabled        bool                   `json:"enabled"`
	Setting        map[string]interface{} `json:"setting"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`

	Organization *OrganizationBo `json:"organization,omitempty"`
}
