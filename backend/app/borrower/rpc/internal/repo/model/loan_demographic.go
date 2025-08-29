package model

import "zeebee-backend/app/borrower/rpc/internal/types"

// LoanTerm 借款人信息
type LoanDemographic struct {
	Base
	LoanId    int64                          `gorm:"column:loan_id;not null;default:0" json:"loan_id"`
	Property  types.LoanDemographicProperty  `gorm:"column:property;type:jsonb;serializer:json" json:"property"`
	Ethnicity types.LoanDemographicEthnicity `gorm:"column:ethnicity;type:jsonb;serializer:json" json:"ethnicity"`
	Race      types.LoanDemographicRace      `gorm:"column:race;type:jsonb;serializer:json" json:"race"`
	Sex       string                         `gorm:"column:sex;size:10;not null;default:'Male'" json:"sex"`
	// ExtraDemographic types.LoanDemographicExtraDemographic `gorm:"column:extra_demographic;type:jsonb;serializer:json" json:"extra_demographic"`
	// relationships
	Loan *Loan `json:"loan" gorm:"foreignKey:LoanId;references:Id"`
}

func (m LoanDemographic) TableName() string {
	return "loan_demographics"
}
