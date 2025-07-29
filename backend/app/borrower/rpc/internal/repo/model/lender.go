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

func (l *Lender) TableName() string {
	return "lenders"
}
