package model

import (
	"golang.org/x/oauth2"

	"zeebee-backend/domains/zeebee/constants"
)

type EmailToken struct {
	oauth2.Token
}

type EmailAuthorization struct {
	Base
	OrganizationId    int64                       `gorm:"column:organization_id;default:0" json:"organization_id"`
	UserId            int64                       `gorm:"column:user_id;default:0" json:"user_id"`
	EmailAuthStrategy constants.EmailAuthStrategy `gorm:"column:email_auth_stategy;default:''" json:"email_auth_stategy"`
	Email             string                      `gorm:"column:email;default:''" json:"email"`
	EmailProvider     constants.EmailProvider     `gorm:"column:email_provider;default:''" json:"email_provider"`
	IsAuthorized      bool                        `gorm:"column:is_authorized;default:false" json:"is_authorized"`
	Token             *EmailToken                 `gorm:"column:token;type:jsonb;serializer:json" json:"token"`
	EmailScope        constants.EmailScope        `gorm:"column:email_scope;default:''" json:"email_scope"`
}

func (m EmailAuthorization) TableName() string {
	return "email_authorizations"
}
