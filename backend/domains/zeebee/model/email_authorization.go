package model

import (
	"golang.org/x/oauth2"

	"zeebee-backend/domains/zeebee/constants"
)

type EmailToken struct {
	*oauth2.Token
}

type EmailAuthorization struct {
	Base
	OrganizationEmailId int64                   `gorm:"column:organization_email_id;default:0" json:"organization_email_id"`
	UserId              int64                   `gorm:"column:user_id;default:0" json:"user_id"`
	Email               string                  `gorm:"column:email;default:''" json:"email"`
	EmailProvider       constants.EmailProvider `gorm:"column:email_provider;default:''" json:"email_provider"`
	IsAuthorized        bool                    `gorm:"column:is_authorized;default:false" json:"is_authorized"`
	Token               *EmailToken             `gorm:"column:token;type:jsonb;serializer:json;default:'{}'" json:"token"`
	Purpose             constants.EmailPurpose  `gorm:"column:purpose;default:''" json:"purpose"`
}

func (m EmailAuthorization) TableName() string {
	return "email_authorizations"
}
