package model

type LoanDemographic struct {
	Base
	LoanId    int64                    `gorm:"column:loan_id;not null;default:0" json:"loan_id"`
	Property  LoanDemographicProperty  `gorm:"column:property;type:jsonb;serializer:json" json:"property"`
	Ethnicity LoanDemographicEthnicity `gorm:"column:ethnicity;type:jsonb;serializer:json" json:"ethnicity"`
	Race      LoanDemographicRace      `gorm:"column:race;type:jsonb;serializer:json" json:"race"`
	Sex       string                   `gorm:"column:sex;size:10;not null;default:'Male'" json:"sex"`
	// Relationships
	Loan *Loan `json:"loan" gorm:"foreignKey:LoanId;references:Id"`
}

func (m LoanDemographic) TableName() string {
	return "loan_demographics"
}

func (m *LoanDemographic) RemoveRelationships() {
	m.Loan = nil
}
