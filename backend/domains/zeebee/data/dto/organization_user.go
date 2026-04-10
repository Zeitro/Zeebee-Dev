package dto

import (
	"time"

	"zeebee-backend/domains/zeebee/constants"
)

type OrganizationUserEntityDto struct {
	Id             int64                            `json:"id"`
	OrganizationId int64                            `json:"organization_id"`
	UserId         int64                            `json:"user_id"`
	Status         constants.OrganizationUserStatus `json:"status"`
	FirstName      string                           `json:"first_name"`
	MidName        string                           `json:"mid_name"`
	LastName       string                           `json:"last_name"`
	FullName       string                           `json:"full_name"`
	Phone          string                           `json:"phone"`
	Remark         string                           `json:"remark"`
	CreatedAt      time.Time                        `json:"created_at"`
	UpdatedAt      time.Time                        `json:"updated_at"`
	Email          string                           `json:"email"`
	IsOwner        bool                             `json:"is_owner"`
	Roles          []*RoleEntityDto                 `json:"roles"`
}

type OrganizationUserCreateRequestDto struct {
	Email          string                    `json:"email"`
	Phone          string                    `json:"phone"`
	Remark         string                    `json:"remark"`
	RoleIds        []int64                   `json:"role_ids"`
	UserType       constants.UserType        `json:"type"`
	OrganizationId int64                     `json:"organization_id"`
	EmailConfig    SendInviteUserEmailConfig `json:"email_config"`
	Operator       string                    `json:"operator"`
}

type OrganizationUserCreateAdminRequestDto struct {
	Email          string             `json:"email"`
	Phone          string             `json:"phone"`
	FirstName      string             `json:"first_name"`
	LastName       string             `json:"last_name"`
	Password       string             `json:"password"`
	Remark         string             `json:"remark"`
	RoleIds        []int64            `json:"role_ids"`
	UserType       constants.UserType `json:"type"`
	OrganizationId int64              `json:"organization_id"`
}

type OrganizationUserUpdateRequestDto struct {
	Id             int64              `json:"id"`
	Phone          string             `json:"phone"`
	FirstName      string             `json:"first_name"`
	MidName        string             `json:"mid_name"`
	LastName       string             `json:"last_name"`
	Remark         string             `json:"remark"`
	RoleIds        []int64            `json:"role_ids"`
	UserType       constants.UserType `json:"type"`
	OrganizationId int64              `json:"organization_id"`
}

type OrganizationUserSendInviteEmailRequestDto struct {
	Id             int64                     `json:"id"`
	UserType       constants.UserType        `json:"type"`
	OrganizationId int64                     `json:"organization_id"`
	EmailConfig    SendInviteUserEmailConfig `json:"email_config"`
	Operator       string                    `json:"operator"`
}

type OrganizationUserPageRequestDto struct {
	BasePageRequestDto
	Filter OrganizationUserPageRequestFilter
	Sorter BasePageSorter
}

type OrganizationUserPageRequestFilter struct {
	BasePageFilter
	OrganizationId int64              `json:"organization_id"`
	UserId         int64              `json:"user_id"`
	UserType       constants.UserType `json:"user_type"`
	RoleId         int64              `json:"role_id"`
}

type SendInviteUserEmailConfig struct {
	HostEmail         string
	TplName           string
	ExpireSeconds     int
	SendPeriodSeconds int
	Url               string
}

type OrganizationUserAcceptInvitationRequestDto struct {
	Code           string             `json:"code"`
	UserType       constants.UserType `json:"type"`
	Email          string             `json:"email"`
	OrganizationId int64              `json:"organization_id"`
	FirstName      string             `json:"first_name"`
	LastName       string             `json:"last_name"`
}

type OrganizationUserProfileEditRequestDto struct {
	OrganizationId int64              `json:"organization_id"`
	UserId         int64              `json:"user_id"`
	UserType       constants.UserType `json:"type"`
	FirstName      string             `json:"first_name"`
	MidName        string             `json:"mid_name"`
	LastName       string             `json:"last_name"`
	Phone          string             `json:"phone"`
}
