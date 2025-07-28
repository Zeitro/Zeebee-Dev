package model

import (
	"time"

	"gorm.io/plugin/soft_delete"

	"zeebee-backend/common/types"
)

type Lender struct {
	ID        int64                 `gorm:"primaryKey" json:"id"`
	Name      string                `gorm:"column:name;size:100;not null;default:''" json:"name"`
	Address   types.Address         `gorm:"column:address;type:jsonb;serializer:json" json:"address"`
	Status    string                `gorm:"column:status;size:10;not null;default:'active'" json:"status"`
	CreatedAt time.Time             `gorm:"column:created_at;default:now()" json:"created_at"`
	UpdatedAt time.Time             `gorm:"column:updated_at;default:now()" json:"updated_at"`
	IsDel     soft_delete.DeletedAt `gorm:"softDelete:flag" json:"-"`

	Loans []Loan `json:"loans" gorm:"foreignKey:LenderID;references:ID"`
}

func (l *Lender) TableName() string {
	return "lenders"
}
