package model

import (
	"zeebee-backend/common/constants"
	"zeebee-backend/common/types"
)

type PropertyBase struct {
	Base

	BorrowerId     int64                  `gorm:"column:borrower_id;not null;default:0" json:"borrower_id"`
	Unit           int32                  `gorm:"column:unit;not null;default:0" json:"unit"`
	PropertyType   constants.PropertyType `gorm:"column:property_type;size:30;not null;default:''" json:"property_type"`
	LivingAreaSqft float32                `gorm:"column:living_area_sqft;type:numeric(8,2);not null;default:0" json:"living_area_sqft"`
	Bedrooms       int32                  `gorm:"column:bedrooms;not null;default:0" json:"bedrooms"`
	Bathrooms      int32                  `gorm:"column:bathrooms;not null;default:0" json:"bathrooms"`
	Address        types.Address          `gorm:"column:address;type:jsonb;serializer:json" json:"address"`
	TotalAmount    float64                `gorm:"column:total_amount;not null;default:0" json:"total_amount"`
	LoanCount      int32                  `gorm:"column:loan_count;not null;default:0" json:"loan_count"`

	Borrower       *Borrower       `json:"borrower" gorm:"foreignKey:BorrowerId;references:Id"`
	Loans          []*Loan         `json:"loans" gorm:"-"`
	Applicants     []*Applicant    `json:"applicants" gorm:"-"`
	LoanProperties []*LoanProperty `json:"loan_properties" gorm:"foreignKey:PropertyBaseId;references:Id"`
}

func (m PropertyBase) TableName() string {
	return "property_bases"
}
