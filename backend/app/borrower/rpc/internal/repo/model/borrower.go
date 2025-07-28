package model

import (
	"time"

	"gorm.io/plugin/soft_delete"

	"zeebee-backend/common/types"
)

type Borrower struct {
	ID          int64                 `gorm:"primaryKey" json:"id"`
	UserID      int64                 `gorm:"column:user_id;not null;default:0" json:"user_id"`
	FirstName   string                `gorm:"column:first_name;size:100;not null;default:''" json:"first_name"`
	LastName    string                `gorm:"column:last_name;size:100;not null;default:''" json:"last_name"`
	Ssn         string                `gorm:"column:ssn;size:20;not null;default:''" json:"ssn"`
	PhoneNumber string                `gorm:"column:phone_number;size:50;not null;default:''" json:"phone_number"`
	Address     types.Address         `gorm:"column:address;type:jsonb;serializer:json" json:"address"`
	CreatedAt   time.Time             `gorm:"column:created_at;default:now()" json:"created_at"`
	UpdatedAt   time.Time             `gorm:"column:updated_at;default:now()" json:"updated_at"`
	IsDel       soft_delete.DeletedAt `gorm:"softDelete:flag" json:"-"`

	User *User `json:"user" gorm:"foreignKey:UserID;references:ID"`
}

func (b *Borrower) TableName() string {
	return "borrowers"
}
