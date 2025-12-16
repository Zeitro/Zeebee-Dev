package bo

import (
	"time"

	"zeebee-backend/domains/zeebee/constants"
)

type UserBo struct {
	Id                    int64                `json:"id"`
	CreatedAt             time.Time            `json:"created_at"`
	UpdatedAt             time.Time            `json:"updated_at"`
	Type                  constants.UserType   `json:"type"`
	Email                 string               `json:"email"`
	Status                constants.UserStatus `json:"status"`
	FirstName             string               `json:"first_name"`
	MidName               string               `json:"mid_name" `
	LastName              string               `json:"last_name" `
	FullName              string               `json:"full_name" `
	Phone                 string               `json:"phone"`
	DefaultOrganizationId int64                `json:"default_organization_id"`
}

type UserinfoBo struct {
	Id                 int64
	OrganizationUserId int64
	Email              string
	FirstName          string
	LastName           string
	FullName           string
	OrganizationId     int64
	OrganizationName   string
	Status             constants.OrganizationUserStatus
	Organizations      []*UserOrganizationBo
	IsOwner            bool
	IsAdmin            bool
}

type UserOrganizationBo struct {
	Id        int64
	Name      string
	Status    constants.OrganizationStatus
	IsDefault bool
}

type UserJwtTokenBo struct {
	AccessToken        string `json:"access_token,omitempty"`
	RefreshToken       string `json:"refresh_token,omitempty"`
	AccessTokenExpire  int64  `json:"access_token_expire,omitempty"`
	RefreshTokenExpire int64  `json:"refresh_token_expire,omitempty"`
}
