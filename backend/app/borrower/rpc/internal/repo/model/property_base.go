package model

import (
	"zeebee-backend/common/types"
)

type PropertyBase struct {
	Base

	BorrowerId  int64         `gorm:"column:borrower_id;not null;default:0" json:"borrower_id"`
	Unit        int32         `gorm:"column:unit;not null;default:0" json:"unit"`
	Address     types.Address `gorm:"column:address;type:jsonb;serializer:json" json:"address"`
	TotalAmount float64       `gorm:"column:total_amount;not null;default:0" json:"total_amount"`

	Borrower       *Borrower       `json:"borrower" gorm:"foreignKey:BorrowerId;references:Id"`
	Loans          []*Loan         `json:"loans" gorm:"-"`
	Applicants     []*Applicant    `json:"applicants" gorm:"foreignKey:PropertyId;references:Id"`
	LoanProperties []*LoanProperty `json:"loan_properties" gorm:"foreignKey:PropertyBaseId;references:Id"`
}

func (m PropertyBase) TableName() string {
	return "property_bases"
}
