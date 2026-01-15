package model

import (
	"zeebee-backend/domains/zeebee/constants"
)

type BorrowerEntity struct {
	Base
	CreatorId      int64                        `gorm:"column:creator_id;type:int4;default:0;not null" json:"creator_id"`
	CreatorName    string                       `gorm:"column:creator_name;type:varchar(100);default:'';not null" json:"creator_name"`
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
	Creator      *User         `json:"creator" gorm:"foreignKey:CreatorId;references:Id"`
	Organization *Organization `json:"organization" gorm:"foreignKey:OrganizationId;references:Id"`
}

func (m BorrowerEntity) TableName() string {
	return "borrower_entities"
}

func (m *BorrowerEntity) RemoveRelationships() {
	m.Creator = nil
	m.Organization = nil
}
