package model

import (
	"zeebee-backend/domains/zeebee/constants"
)

type OrganizationEmail struct {
	Base
	OrganizationId int64                       `gorm:"column:organization_id;default:0" json:"organization_id"`
	AuthStrategy   constants.EmailAuthStrategy `gorm:"column:auth_strategy;default:''" json:"auth_strategy"`
	Email          string                      `gorm:"column:email;default:''" json:"email"`
	Purpose        constants.EmailPurpose      `gorm:"column:purpose;default:''" json:"purpose"`

	EmailAuthorizations []*EmailAuthorization `json:"email_authorizations" gorm:"foreignKey:OrganizationEmailId"`
}

func (m OrganizationEmail) TableName() string {
	return "organization_emails"
}

func (m *OrganizationEmail) RemoveRelationships() {
	m.EmailAuthorizations = nil
}
