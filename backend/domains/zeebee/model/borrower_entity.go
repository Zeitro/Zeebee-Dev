package model

import (
	"zeebee-backend/common/constants"
)

type BorrowerEntity struct {
	Base
	UserId         int64                        `gorm:"column:user_id;type:int4;default:0;not null" json:"user_id"`
	OrganizationId int64                        `gorm:"column:organization_id;type:int4;default:0;not null" json:"organization_id"`
	Type           constants.ApplicantType      `gorm:"column:type;type:varchar(20);default:'';not null" json:"type"`
	Name           string                       `gorm:"column:name;type:varchar(100);default:'';not null" json:"name"`
	TinType        constants.BorrowerEntityType `gorm:"column:tin_type;type:varchar(20);default:'';not null" json:"tin_type"`
	Tin            string                       `gorm:"column:tin;type:varchar(20);default:'';not null" json:"tin"`
	MailingAddress Address                      `gorm:"column:mailing_address;type:jsonb;default:'{}';serializer:json" json:"mailing_address"`
	ContactEmail   string                       `gorm:"column:contact_email;type:varchar(100);default:'';not null" json:"contact_email"`
	ContactName    string                       `gorm:"column:contact_name;type:varchar(100);default:'';not null" json:"contact_name"`
	ContactPhone   string                       `gorm:"column:contact_phone;type:varchar(30);default:'';not null" json:"contact_phone"`

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
