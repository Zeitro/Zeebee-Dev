package dto

import (
	"golang.org/x/oauth2"

	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/email"
)

type OrganizationEmailSaveEmailRequestDto struct {
	OrganizationId   int64                       `json:"organization_id"`
	UserId           int64                       `json:"user_id"`
	Email            string                      `json:"email"`
	Purpose          constants.EmailPurpose      `json:"Purpose"`
	AuthStrategy     constants.EmailAuthStrategy `json:"auth_strategy"`
	OrganizationType constants.UserType          `json:"organization_type"`
}

type OrganizationEmailGenerateAuthCodeUrlRequestDto struct {
	OrganizationId    int64                   `json:"organization_id"`
	UserId            int64                   `json:"user_id"`
	EmailProvider     constants.EmailProvider `json:"email_provider"`
	Purpose           constants.EmailPurpose  `json:"Purpose"`
	OrganizationType  constants.UserType      `json:"organization_type"`
	AuthCodeGenerator email.AuthCodeUrlGenerator
}

type OrganizationEmailAuthorizeRequestDto struct {
	Code             string                  `json:"code"`
	State            string                  `json:"state"`
	EmailProvider    constants.EmailProvider `json:"email_provider"`
	Purpose          constants.EmailPurpose  `json:"Purpose"`
	OrganizationId   int64                   `json:"organization_id"`
	UserId           int64                   `json:"user_id"`
	OrganizationType constants.UserType      `json:"organization_type"`
	Oauth2Config     *oauth2.Config
}
