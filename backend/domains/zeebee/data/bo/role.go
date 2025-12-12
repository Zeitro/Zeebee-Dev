package bo

import (
	"time"

	"zeebee-backend/domains/zeebee/constants"
)

type RoleBo struct {
	Id             int64                  `json:"id"`
	Name           string                 `json:"name"`
	Description    string                 `json:"description"`
	Type           constants.UserType     `json:"type"`
	Status         constants.CommonStatus `json:"status"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
	OrganizationId int64                  `json:"organization_id"`
	Resources      []*ResourceBo          `json:"resources"`
}
