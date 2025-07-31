package model

import (
	"zeebee-backend/common/types"
)

type Lender struct {
	Base

	Name    string        `gorm:"column:name;size:100;not null;default:''" json:"name"`
	Address types.Address `gorm:"column:address;type:jsonb;serializer:json" json:"address"`
	Status  string        `gorm:"column:status;size:10;not null;default:'active'" json:"status"`

	Loans []Loan `json:"loans" gorm:"foreignKey:LenderId;references:Id"`
}

func (m *Lender) TableName() string {
	return "lenders"
}

func (m Lender) FieldName() string {
	return "name"
}

func (m Lender) FieldAddress() string {
	return "address"
}

func (m Lender) FieldStatus() string {
	return "status"
}
