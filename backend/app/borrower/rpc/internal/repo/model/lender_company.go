package model

import (
	"zeebee-backend/common/types"
)

type LenderCompany struct {
	Base

	Name    string        `gorm:"column:name;size:100;not null;default:''" json:"name"`
	Address types.Address `gorm:"column:address;type:jsonb;serializer:json" json:"address"`
	Status  string        `gorm:"column:status;size:10;not null;default:'active'" json:"status"`

	Loans []Loan `json:"loans" gorm:"foreignKey:LenderCompanyId;references:Id"`
}

func (m *LenderCompany) TableName() string {
	return "lender_companies"
}
