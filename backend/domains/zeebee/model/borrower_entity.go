package model

import (
	"zeebee-backend/common/constants"
	"zeebee-backend/common/types"
)

type BorrowerEntity struct {
	Base
	UserId         int64                        `gorm:"column:user_id;default:0" json:"user_id"`
	OrganizationId int64                        `gorm:"column:organization_id;default:0" json:"organization_id"`
	Type           constants.ApplicantType      `gorm:"column:type;default:''" json:"type"`
	Name           string                       `gorm:"column:name;default:''" json:"name"`
	TinType        constants.BorrowerEntityType `gorm:"column:tin_type;default:''" json:"tin_type"`
	Tin            string                       `gorm:"column:tin;default:''" json:"tin"`
	MailingAddress types.Address                `gorm:"column:mailing_address;default:'';serializer:json" json:"mailing_address"`
	ContactEmail   string                       `gorm:"column:contact_email;default:''" json:"contact_email"`
	ContactName    string                       `gorm:"column:contact_name;default:''" json:"contact_name"`
	ContactPhone   string                       `gorm:"column:contact_phone;default:''" json:"contact_phone"`

	//  relationships
	User         *User         `json:"user" gorm:"foreignKey:UserId;references:Id"`
	Organization *Organization `json:"organization" gorm:"foreignKey:OrganizationId;references:Id"`
}

func (m BorrowerEntity) TableName() string {
	return "borrower_entities"
}

func (m *BorrowerEntity) RemoveRelationships() {
	m.User = nil
	m.Organization = nil
}
