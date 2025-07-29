package model

import (
	"zeebee-backend/common/constants"
	"zeebee-backend/common/types"
)

type Property struct {
	Base

	BorrowerId   int64                  `gorm:"column:borrower_id;not null;default:0" json:"borrower_id"`
	PropertyType constants.PropertyType `gorm:"column:property_type;size:30;not null;default:''" json:"property_type"`
	Unit         int32                  `gorm:"column:unit;not null;default:0" json:"unit"`
	Address      types.Address          `gorm:"column:address;type:jsonb;serializer:json" json:"address"`

	Borrower *Borrower `json:"borrower" gorm:"foreignKey:BorrowerId;references:Id"`
}

func (p *Property) TableName() string {
	return "properties"
}
