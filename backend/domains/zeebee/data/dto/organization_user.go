package dto

import (
	"zeebee-backend/common/constants"
)

type OrganizationUserCreateRequestDto struct {
	Email          string                    `json:"email"`
	Phone          string                    `json:"phone"`
	Remark         string                    `json:"remark"`
	RoleIds        []int64                   `json:"role_ids"`
	Type           constants.UserType        `json:"type"`
	OrganizationId int64                     `json:"organization_id"`
	EmailConfig    SendInviteUserEmailConfig `json:"email_config"`
	Operator       string                    `json:"operator"`
}

type SendInviteUserEmailConfig struct {
	HostEmail         string
	TplName           string
	ExpireSeconds     int
	SendPeriodSeconds int
	Url               string
}
