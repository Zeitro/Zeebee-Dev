package dto

import (
	"time"

	"zeebee-backend/common/redisx"
	"zeebee-backend/domains/zeebee/constants"
)

type UserEntityDto struct {
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

type UserinfoDto struct {
	Id                 int64
	OrganizationUserId int64
	Email              string
	FirstName          string
	LastName           string
	FullName           string
	OrganizationId     int64
	OrganizationName   string
	Status             constants.OrganizationUserStatus
	Organizations      []*UserOrganizationDto
	IsOwner            bool
	IsAdmin            bool
}

type UserOrganizationDto struct {
	Id        int64
	Name      string
	Status    constants.OrganizationStatus
	IsDefault bool
}

type UserJwtTokenDto struct {
	AccessToken        string `json:"access_token,omitempty"`
	RefreshToken       string `json:"refresh_token,omitempty"`
	AccessTokenExpire  int64  `json:"access_token_expire,omitempty"`
	RefreshTokenExpire int64  `json:"refresh_token_expire,omitempty"`
}

type UserChangePasswordRequestDto struct {
	OldPassword string             `json:"old_password"`
	NewPassword string             `json:"new_password"`
	UserId      int64              `json:"user_id"`
	UserType    constants.UserType `json:"user_type"`
}

type UserResetPasswordByEmailRequestDto struct {
	UserType    constants.UserType `json:"user_type"`
	NewPassword string             `json:"new_password"`
	Email       string             `json:"email"`
	Code        string             `json:"code"`
}

type UserLoginRequestDto struct {
	Email    string             `json:"email"`
	Password string             `json:"password"`
	UserType constants.UserType `json:"user_type"`
}

type LoginLimiterDto struct {
	Rate     int
	Period   int
	RedisCli *redisx.Redisx
}

type UserRegisterRequestDto struct {
	Email                 string               `json:"email"`
	PasswordHash          string               `json:"password_hash"`
	Status                constants.UserStatus `json:"status"`
	Type                  constants.UserType   `json:"type"`
	FirstName             string               `json:"first_name"`
	LastName              string               `json:"last_name"`
	FullName              string               `json:"full_name"`
	DefaultOrganizationId int64                `json:"default_organization_id"`
}
