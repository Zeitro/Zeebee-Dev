package bo

import (
	"time"

	"zeebee-backend/domains/zeebee/constants"
)

type EmailAuthorizationBo struct {
	Id                int64                       `json:"id"`
	OrganizationId    int64                       `json:"organization_id"`
	Email             string                      `json:"email"`
	EmailAuthStrategy constants.EmailAuthStrategy `json:"strategy"`
	EmailProvider     constants.EmailProvider     `json:"email_provider"`
	IsAuthorized      bool                        `json:"is_authorized"`
	Purpose           constants.EmailPurpose      `json:"purpose"`
	CreatedAt         time.Time                   `json:"created_at"`
	UpdatedAt         time.Time                   `json:"updated_at"`
}
