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
	TotalAmount  float64                `gorm:"column:total_amount;not null;default:0" json:"total_amount"`

	Borrower   *Borrower    `json:"borrower" gorm:"foreignKey:BorrowerId;references:Id"`
	Loans      []*Loan      `json:"loans" gorm:"foreignKey:PropertyId;references:Id"`
	Applicants []*Applicant `json:"applicants" gorm:"foreignKey:PropertyId;references:Id"`
}

func (m Property) TableName() string {
	return "properties"
}
