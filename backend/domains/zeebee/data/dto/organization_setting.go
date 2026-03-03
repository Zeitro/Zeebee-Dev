package dto

import "time"

type OrganizationSettingEntityDto struct {
	Id             int64                  `json:"id"`
	Purpose        string                 `json:"purpose"`
	OrganizationId int64                  `json:"organization_id"`
	Enabled        bool                   `json:"enabled"`
	Setting        map[string]interface{} `json:"setting"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`

	Organization *OrganizationEntityDto `json:"organization,omitempty"`
}
