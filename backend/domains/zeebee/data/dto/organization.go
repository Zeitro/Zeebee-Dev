package dto

import "zeebee-backend/domains/zeebee/constants"

type OrganizationSaveEmailRequestDto struct {
	OrganizationId    int64                       `json:"organization_id"`
	UserId            int64                       `json:"user_id"`
	Email             string                      `json:"email"`
	EmailScope        constants.EmailScope        `json:"email_scope"`
	EmailAuthStrategy constants.EmailAuthStrategy `json:"email_auth_strategy"`
	OrganizationType  constants.UserType          `json:"organization_type"`
}
