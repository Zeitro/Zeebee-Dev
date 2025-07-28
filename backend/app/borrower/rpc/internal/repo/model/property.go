package model

import (
	"time"

	"gorm.io/plugin/soft_delete"

	"zeebee-backend/common/constants"
	"zeebee-backend/common/types"
)

type Property struct {
	ID         int64                  `gorm:"primaryKey" json:"id"`
	BorrowerID int64                  `gorm:"column:borrower_id;not null;default:0" json:"borrower_id"`
	Type       constants.PropertyType `gorm:"column:type;size:30;not null;default:''" json:"type"`
	Unit       int32                  `gorm:"column:unit;not null;default:0" json:"unit"`
	Address    types.Address          `gorm:"column:address;type:jsonb;serializer:json" json:"address"`
	CreatedAt  time.Time              `gorm:"column:created_at;default:now()" json:"created_at"`
	UpdatedAt  time.Time              `gorm:"column:updated_at;default:now()" json:"updated_at"`
	IsDel      soft_delete.DeletedAt  `gorm:"softDelete:flag" json:"-"`

	Borrower *Borrower `json:"borrower" gorm:"foreignKey:BorrowerID;references:ID"`
}

func (p *Property) TableName() string {
	return "properties"
}
