package model

import (
	"zeebee-backend/common/constants"
	"zeebee-backend/common/types"
)

type BorrowerEntity struct {
	Base
	BorrowerId     int64                        `gorm:"column:borrower_id;default:0" json:"borrower_id"`
	Type           constants.ApplicantType      `gorm:"column:type;default:''" json:"type"`
	Name           string                       `gorm:"column:name;default:''" json:"name"`
	TinType        constants.BorrowerEntityType `gorm:"column:tin_type;default:''" json:"tin_type"`
	Tin            string                       `gorm:"column:tin;default:''" json:"tin"`
	MailingAddress types.Address                `gorm:"column:mailing_address;default:'';serializer:json" json:"mailing_address"`
	ContactEmail   string                       `gorm:"column:contact_email;default:''" json:"contact_email"`
	ContactName    string                       `gorm:"column:contact_name;default:''" json:"contact_name"`
	ContactPhone   string                       `gorm:"column:contact_phone;default:''" json:"contact_phone"`
	OwnPercent     int32                        `gorm:"column:own_percent;default:100" json:"own_percent"`
}

func (m BorrowerEntity) TableName() string {
	return "borrower_entities"
}
