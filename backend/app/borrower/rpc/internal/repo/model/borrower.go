package model

import (
	"zeebee-backend/common/types"
)

type Borrower struct {
	Base

	UserId      int64         `gorm:"column:user_id;not null;default:0" json:"user_id"`
	FirstName   string        `gorm:"column:first_name;size:100;not null;default:''" json:"first_name"`
	LastName    string        `gorm:"column:last_name;size:100;not null;default:''" json:"last_name"`
	Ssn         string        `gorm:"column:ssn;size:20;not null;default:''" json:"ssn"`
	PhoneNumber string        `gorm:"column:phone_number;size:50;not null;default:''" json:"phone_number"`
	Address     types.Address `gorm:"column:address;type:jsonb;serializer:json" json:"address"`

	User *User `json:"user" gorm:"foreignKey:UserId;references:Id"`
}

func (b *Borrower) TableName() string {
	return "borrowers"
}
