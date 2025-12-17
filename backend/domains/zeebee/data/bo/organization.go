package bo

import (
	"time"

	"zeebee-backend/domains/zeebee/constants"
)

type OrganizationBo struct {
	Id                  int64                           `json:"id"`
	Name                string                          `json:"name"`
	Status              constants.OrganizationStatus    `json:"status"`
	OwnerId             int64                           `json:"owner_id"`
	Address             *AddressBo                      `json:"address"`
	ContactEmail        string                          `json:"contact_email"`
	EmailAuthStrategy   constants.EmailAuthStrategy     `json:"email_auth_strategy"`
	OwnerType           constants.OrganizationOwnerType `json:"owner_type"`
	Type                constants.UserType              `json:"type"`
	EmailAuthorizations []*EmailAuthorizationBo         `json:"email_authorizations"`
	CreatedAt           time.Time                       `json:"created_at"`
	UpdatedAt           time.Time                       `json:"updated_at"`
}
