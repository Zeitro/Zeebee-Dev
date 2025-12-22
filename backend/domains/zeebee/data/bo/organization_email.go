package bo

import (
	"time"

	"zeebee-backend/domains/zeebee/constants"
)

type OrganizationEmailBo struct {
	Id             int64                       `json:"id"`
	OrganizationId int64                       `json:"organization_id"`
	AuthStrategy   constants.EmailAuthStrategy `json:"auth_strategy"`
	Email          string                      `json:"email"`
	Purpose        constants.EmailPurpose      `json:"purpose"`
	CreatedAt      time.Time                   `json:"created_at"`
	UpdatedAt      time.Time                   `json:"updated_at"`
	IsAuthorized   bool                        `json:"is_authorized"`
}
